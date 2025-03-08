package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	db "github.com/Rich-T-kid/Notiffy/internal/database" // this package needs to always to be run first b4 all other custom packages
	logger "github.com/Rich-T-kid/Notiffy/internal/log"
	"github.com/Rich-T-kid/Notiffy/pkg"
)

// TODO: Implement the rest of the notificationService (Start,Notify,SendDirectMessage)

var (
	ErrInvalidUserType = func(inputType string, expectedTypes []string) error {
		return fmt.Errorf("incorrect type passed into notification service. can accept %v  but recieved %s", expectedTypes, inputType)
	}
	ErrInvalidUserObject = errors.New("invalid user object passed in")
	ErrUsernameExists    = errors.New("username already exists, pick a new one")
	ErrUserInfoError     = errors.New("error in your userInfo")
	ErrNotImplemented    = errors.New("TODO: Implement")
	ErrUserMustExist     = func(userName string) error {
		return fmt.Errorf("user: %s must exist before you deregister them", userName)
	}
	ErrInvalidMessengerPassed = func(extra error) error {
		return fmt.Errorf("invalid messenger struct passed to notify function.additional context from validate %w", extra)
	}
)

type status string

const (
	DELIVERED status = "DELIVERED"
	SENT      status = "SENT"
	SENDING   status = "SENDING"
	FAILED    status = "FAILED"
	UNKNOWN   status = "UNKNOWN"
)

func init() {
	//print("richard make sure that the third party API that you going to intergrate with is up and running in this init func\n")
}

type RegisterINFO struct {
	Name    string `bson:"Name"`
	Contact int64  `bson:"Contact"` // phone number
	Tags    []Tag  `bson:"Tags"`
}

// in the future i could have this struct be composed of sub-structs but for now thats not needed
type SMSBody struct {
	from       string
	title      string
	categories tags
	mainText   string
}

// Messenger interface
func (s *SMSBody) Message() MessageBody  { return s }
func (s *SMSBody) Metadata() MessageMeta { return s }
func (s *SMSBody) Type() string          { return "sms" }

// Message body interface
func (s *SMSBody) Content() interface{} { return s.mainText } // actuall text to be sent to reciever

// Message Metadata interface
func (s *SMSBody) Tags() tags       { return s.categories } // target of broadcast message
func (s *SMSBody) Priority() int    { return 1 }            // for now this is hardcoded as 1
func (s *SMSBody) Timestamp() int64 { return time.Now().Unix() }
func (s *SMSBody) Title() string    { return s.title }
func (s *SMSBody) From() string     { return s.from }

type SMSNotification struct {
	apiKey string
	saver  db.Storage
	db     *mongo.Database // store as mongo collection for now
}
type textBeltMSGresponse struct {
	Success        bool   `json:"success"`
	TextId         string `json:"textId"`
	QuotaRemaining int    `json:"quotaRemaining"`
}

func DefineMessage(from, title, text string) *SMSBody {
	return &SMSBody{
		from:       from,
		title:      title,
		categories: tags{"sms"},
		mainText:   text,
	}
}
func NewSMSUserService() UserService {
	return NewSMSNotification()
}

func NewSMSService() NotificationService {
	return NewSMSNotification()
}

func NewSMSNotification() *SMSNotification {
	apikey := os.Getenv("TEXTBELT_API_KEY")
	// defensive programing :)
	if apikey == "" {
		log.Fatal("textbelt api key missing")
	}
	return &SMSNotification{
		apiKey: apikey,
		db:     db.EstablishMongoConnection(),
	}
}

// there should be a target to notify based on the filter
// notify and ensure that everyone that should have been notified was. if not return an error( this error will contain the users that were notified)
func (s *SMSNotification) Notify(ctx context.Context, body Messenger, filter Filter) (int, []error) {
	requestid, ok := ctx.Value(pkg.RequestIDKey{}).(string)
	if !ok {
		logger.Critical("context request ID is not present when it should be!")
		return 0, []error{pkg.ErrMissingRequestID}
	}
	startTime, ok := ctx.Value(pkg.StartTime{}).(time.Time)
	if !ok {
		logger.Critical("context start time is not present when it should be!")
		return 0, []error{pkg.ErrMissingStartTime}
	}
	formatedTime := startTime.Format("Mon, 02 Jan 2006 15:04:05 MST")
	if err := s.Validate(body); err != nil {
		return 0, []error{ErrInvalidMessengerPassed(err)}
	}
	var errorArray []error
	// again for now we are just going to send messages to every registerd user in the database. we can have users provide their clients later
	// this will be pretty basic. go through the mongo collection. collect all numbers and send a notification to each users number. if for what ever reason it cant be complete
	// add that users name to a string array and return it in the error
	collection := s.db.Collection("SMS")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return 0, []error{fmt.Errorf("requestId: %s failed to find documents: %w", requestid, err)}
	}
	defer cursor.Close(ctx)
	var res []RegisterINFO
	if err = cursor.All(ctx, &res); err != nil {
		return 0, []error{fmt.Errorf("requestId: %s failed to decode cursor results: %w", requestid, err)}
	}
	logger.Info(fmt.Sprintf("requestId: %s all users in database %v", requestid, res))
	var toNotify []RegisterINFO
	for _, i := range res {
		if filter(ctx, i) {
			toNotify = append(toNotify, i)
		}
	}
	var notified int
	// can wrap this in a concurent function to improve speed later if needed
	for _, user := range toNotify {
		phoneStr := strconv.FormatInt(user.Contact, 10)
		response, err := s.sendTXT(phoneStr, body.Message().Content().(string))
		if err != nil {
			errorArray = append(errorArray, err)
		}
		// im thinking of adding a context ID to be able to track things
		status, _ := s.messageStatus(response.TextId)
		if status != FAILED && status != "" {
			notified += 1
		}
		logger.Info(fmt.Sprintf("requestId: %s message sent from %s to number (%s) has resulted in this response status from textbelts API %v", requestid, body.Metadata().From(), phoneStr, status))

	}
	if len(toNotify) != notified {
		// TODO: not critical but it would be nice for callers to have more information about the users that we not notified
		str := fmt.Sprintf("requestId: %s notify should have sent %d messages but only sent %d context startTime: %s", requestid, len(toNotify), notified, formatedTime)
		errorArray = append(errorArray, errors.New(str))
	}

	return notified, errorArray
}

// dont waste api request. Ive tested this and i know it works
func (s *SMSNotification) SendDirectMessage(ctx context.Context, body Messenger, from string, recipient_usernames []string) []error {
	requestid, ok := ctx.Value(pkg.RequestIDKey{}).(string)
	if !ok {
		logger.Critical("context request ID is not present when it should be!")
		return []error{pkg.ErrMissingRequestID}
	}
	if len(recipient_usernames) == 0 {
		return []error{errors.New("recipiant list cannot be empty")}
	}
	newMessage := fmt.Sprintf("from : %s, %s", from, body.Message().Content().(string))
	userInfo, err := s.userbyID(ctx, recipient_usernames)
	if err != nil {
		return []error{err}
	}
	if len(userInfo) == 0 {
		return []error{errors.New("no valid usersnames were passed in. users must already be registred to be notified")}
	}
	var errorArray []error
	for _, user := range userInfo {
		phoneStr := strconv.FormatInt(user.Contact, 10)
		txResponse, err := s.sendTXT(phoneStr, newMessage)
		if err != nil {
			logger.Warn(fmt.Sprintf("requestID: %s  SendText function failed with error:%v", requestid, err))
			errorArray = append(errorArray, err)
		}
		status, err := s.messageStatus(txResponse.TextId)
		if err != nil {
			logger.Warn(fmt.Sprintf("requestID: %s SendText function failed with error:%v", requestid, err))
			errorArray = append(errorArray, err)

		}
		logger.Info(fmt.Sprintf("textId:%s has a current status of %s", txResponse.TextId, status))
	}
	return errorArray
}

// doesnt need to do alot right now. mabey ubon an extreranl api intergration butkeep minimal for as long as possible
func (s *SMSNotification) Start(ctx context.Context) error {
	requestid, ok := ctx.Value(pkg.RequestIDKey{}).(string)
	if !ok {
		logger.Critical("context request ID is not present when it should be!")
		return pkg.ErrMissingRequestID
	}
	_, err := s.quota()
	fmt.Printf("requestId: %s has began the SMSNotification service\n", requestid)
	return err
}

// add to db.SMS collection. for now these all go under the same Mangager(person users register under for notification) but when we build an external api we should take that in
func (s *SMSNotification) Register(ctx context.Context, userInfo Validator, subcategory []Tag) error {
	user, ok := userInfo.(*RegisterINFO)
	if !ok {
		// This should return an error but since i havnt handled the Validate file yet this will just crash
		logger.Critical("Incorrect type passed into SMS Notification Service")
		givenType := fmt.Sprintf("%T", userInfo)
		expectedTypeName := fmt.Sprintf("%T", &RegisterINFO{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})

	}
	collection := s.db.Collection("SMS")
	if err := userInfo.Validate(); err != nil {
		return ErrInvalidUserObject
	}
	exist := s.exist(user.Name)
	if exist {
		return ErrUsernameExists
	}
	// By defualt every user registered under this has an SMS tag
	user.Tags = append(user.Tags, Tag("SMS"))
	logger.Info(fmt.Sprintf("Inserting user %s into SMS collection", user.Name))
	_, err := collection.InsertOne(ctx, userInfo)
	return err
}

// interface defines the pretty well
func (s *SMSNotification) Unregister(ctx context.Context, userInfo Validator, subcategories []Tag) error {
	user, ok := userInfo.(*RegisterINFO)
	if !ok {
		// This should return an error but since i havnt handled the Validate file yet this will just crash
		logger.Critical("Incorrect type passed into SMS Notification Service")
		givenType := fmt.Sprintf("%T", userInfo)
		expectedTypeName := fmt.Sprintf("%T", &RegisterINFO{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})

	}
	if err := userInfo.Validate(); err != nil {
		return ErrInvalidUserObject
	}
	collection := s.db.Collection("SMS")

	exist := s.exist(user.Name)
	if !exist {
		return ErrUserMustExist(user.Name)
	}
	filter := bson.D{{Key: "Name", Value: user.Name}}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// Pass in the final object youd like to be stored in the Database.
// this means to add a tag to a user pass in that tag along with their user collection representation
// to remove a tag pass in their user collection representation and pass in the tags youd like to remain
func (s *SMSNotification) UpdateRegistration(ctx context.Context, userInfo Validator, subcategories []Tag) error {
	user, ok := userInfo.(*RegisterINFO)
	if !ok {
		logger.Critical("Incorrect type passed into SMS Notification Service")
		givenType := fmt.Sprintf("%T", userInfo)
		expectedTypeName := fmt.Sprintf("%T", &RegisterINFO{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})
	}
	if err := userInfo.Validate(); err != nil {
		return ErrInvalidUserObject
	}

	collection := s.db.Collection("SMS")
	filter := bson.D{{Key: "Name", Value: user.Name}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "Tags", Value: bson.M{"$in": subcategories}},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update registration: %w", err)
	}

	return nil
}
func (s *SMSNotification) ListUsers(filter Filter) ([]string, error) {
	collection := s.db.Collection("SMS")

	// Retrieve all registered users from the SMS collection
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		logger.Critical(fmt.Sprintf("Failed to retrieve users: %v", err))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []RegisterINFO
	if err = cursor.All(context.TODO(), &results); err != nil {
		logger.Critical(fmt.Sprintf("Failed to parse users: %v", err))
		return nil, err
	}

	var filteredUsers []string
	for _, user := range results {
		if filter(context.TODO(), user.Name, user.Tags) {
			filteredUsers = append(filteredUsers, user.Name)
		}
	}

	return filteredUsers, nil
}

// be carful not to confuse the function signatures below. they accept and return different types
// :aeaws
func (r *RegisterINFO) Validate() error {
	// Name cannot already be taken
	// Phone number must be valid
	if r.Name == "" {
		return errors.New("name cannot be empty")
	}
	if r.Contact == 0.0 {
		return errors.New("contact number must be a valid phone number")
	}
	return nil

}
func (s *SMSNotification) Validate(input Messenger) error {
	body, ok := input.(*SMSBody)
	if !ok {
		logger.Info("Incorrect type passed into SMS Notification Service")
		givenType := fmt.Sprintf("%T", input)
		expectedTypeName := fmt.Sprintf("%T", &SMSBody{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})
	}

	if body.from == "" {
		return errors.New("sender 'From' field cannot be empty")
	}
	if body.title == "" {
		return errors.New("title field cannot be empty")
	}
	if body.mainText == "" {
		return errors.New("main text field cannot be empty")
	}
	if len(body.categories) == 0 {
		return errors.New("categories field cannot be empty")
	}

	return nil
}

// Internal functions
func (s *SMSNotification) sendTXT(number, message string) (*textBeltMSGresponse, error) {
	// returns 200 even if it doesnt work need to check the success value
	// invalid numbers to cost us api calls but number validation should still be done on input
	var endpoint = "https://textbelt.com/text"
	// none of these fields can be empty. numbers must exist. message bodys cannot be nothing
	values := url.Values{
		"phone":   {number},
		"message": {message},
		"key":     {s.apiKey},
	}

	resp, err := http.PostForm(endpoint, values)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorString := fmt.Sprintf("TextBelt Responded with non 200 status code for textMessage endpoint, status code %d", resp.StatusCode)
		logger.Critical(errorString)
		return nil, errors.New(errorString)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading textBelt response body %w", err)
	}
	var response textBeltMSGresponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshelling textbelt json response %w", err)
	}
	return &response, nil

}
func (s *SMSNotification) userbyID(ctx context.Context, toFind []string) ([]RegisterINFO, error) {
	collection := s.db.Collection("SMS")

	// Correct filter for MongoDB query
	filter := bson.M{"Name": bson.M{"$in": toFind}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)
	var res []RegisterINFO
	if err = cursor.All(ctx, &res); err != nil {
		return nil, fmt.Errorf("failed to decode cursor results: %w", err)
	}
	return res, nil
}
func (s *SMSNotification) quota() (int, error) {
	endpoint := fmt.Sprintf("https://textbelt.com/quota/%s", s.apiKey)
	tempResponse := struct {
		Success        bool `json:"success"`
		QuotaRemaining int  `json:"quotaRemaining"`
	}{}
	c := http.Client{Timeout: time.Second * 3}
	response, err := c.Get(endpoint)
	if err != nil {
		return 0, fmt.Errorf("error requesting endpoint %s: recieved error -> %w", endpoint, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid status code: %d returned from endpoint: %s", response.StatusCode, endpoint)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error parsing response body from endpoint: %s", endpoint)
	}
	if err := json.Unmarshal(body, &tempResponse); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %w", err)
	}
	logger.Info(fmt.Sprintf("remaining Api requst %d\n", tempResponse.QuotaRemaining))
	return tempResponse.QuotaRemaining, nil
}
func (s *SMSNotification) messageStatus(messageid string) (status, error) {
	var endpoint = fmt.Sprintf("https://textbelt.com/status/%s", messageid)
	response := struct {
		Success bool   `json:"success"`
		Status  string `json:"status"`
	}{}
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	switch status(response.Status) {
	case DELIVERED, SENT, SENDING, FAILED:
		return status(response.Status), nil
	default:
		if response.Status != "UNKNOWN" {
			logger.Warn(fmt.Sprintf("textbelt responded with non standard response. Full response %+v\n", response))
		}
		return UNKNOWN, nil
	}

}

func (s *SMSNotification) exist(name string) bool {
	collection := s.db.Collection("SMS")
	filter := bson.D{{Key: "Name", Value: name}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Printf("Error checking if user exists: %v", err)
		return false
	}
	fmt.Printf("Name:%s exist %d times in databse\n", name, count)
	return count > 0
}
