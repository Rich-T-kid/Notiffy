package services

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	db "github.com/Rich-T-kid/Notiffy/internal/database"
	logger "github.com/Rich-T-kid/Notiffy/internal/log"
	"github.com/Rich-T-kid/Notiffy/pkg"
)

var (
	ErrInvalidEmailResponse = func(reqId string, err error) error {
		return fmt.Errorf("request id %s resulted in invalid email response from custom email client %w", reqId, err)
	}
)

const (
	Email_Title = "Notiffy"
)

func NewMailNotiffyer() NotificationService {
	return NewMailer()
}
func NewMailRegister() UserService {
	return NewMailer()
}

type EmailReigisterInfo struct {
	Name  string `bson:"Name"`
	Email string `bson:"Email"`
	Tags  Tags   `bson:"Tags"`
}

func (e *EmailReigisterInfo) Validate() error {
	if e.Name == "" {
		return errors.New("user's name cannot be empty")
	}
	if e.Email == "" {
		return errors.New("user's Email cannot be empty")
	}
	if err := ValidateEmail(e.Email); err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}
	return nil
}

type register struct {
	db *mongo.Database
}
type Mailer struct {
	register
	senderEmail    string
	senderPassword string
	db             *mongo.Database
}
type Mailbody struct {
	Subject  string // I.E Subject
	MailList Tags
	Body     string
	To       string // Recipiant
}

// MessageMeta
func (m *Mailbody) Tags() Tags       { return m.MailList }
func (m *Mailbody) Priority() int    { return 1 }
func (m *Mailbody) Timestamp() int64 { return time.Now().Unix() }
func (m *Mailbody) Title() string    { return m.Subject }
func (m *Mailbody) From() string     { return m.To } // not applicable here. all emails are sent from the same account
// call it bad design lol but this is a convient implentation
func (m *Mailbody) Type() string { return "email" }

// MessageBody
func (m *Mailbody) Content() interface{}  { return m.Body }
func (m *Mailbody) Message() MessageBody  { return m }
func (m *Mailbody) Metadata() MessageMeta { return m }

func NewMailer() *Mailer {
	mongodb := db.EstablishMongoConnection()
	return &Mailer{
		register:       register{db: mongodb},
		senderEmail:    os.Getenv("SENDER_EMAIL"),
		senderPassword: os.Getenv("GOOGLE_GMAIL_PASSWORD"),
		db:             mongodb,
	}
}
func DefineMail(subject, body, to string, topics Tags) *Mailbody {
	return &Mailbody{
		Subject:  subject,
		MailList: topics,
		Body:     body,
		To:       to,
	}
}
func (r *register) Register(ctx context.Context, userInfo Validator, subcategory []Tag) error {
	user, ok := userInfo.(*EmailReigisterInfo)
	if !ok {
		// This should return an error but since i havnt handled the Validate file yet this will just crash
		logger.Critical("Incorrect type passed into Email Notification Service")
		givenType := fmt.Sprintf("%T", userInfo)
		expectedTypeName := fmt.Sprintf("%T", &EmailReigisterInfo{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})

	}
	if err := userInfo.Validate(); err != nil {
		return ErrInvalidUserObject
	}
	exist := r.exist(user.Name)
	if exist {
		return ErrUsernameExists
	}
	// By defualt every user registered under this has an SMS tag
	collection := r.db.Collection("EMAIL")
	user.Tags = append(user.Tags, Tag("EMAIL"))
	logger.Info(fmt.Sprintf("Inserting user %s into EMAIL collection", user.Name))
	_, err := collection.InsertOne(ctx, userInfo)
	return err
}
func (r *register) Unregister(ctx context.Context, userInfo Validator, subcategory []Tag) error {
	user, ok := userInfo.(*EmailReigisterInfo)
	if !ok {
		// This should return an error but since i havnt handled the Validate file yet this will just crash
		logger.Critical("Incorrect type passed into Email Notification Service")
		givenType := fmt.Sprintf("%T", userInfo)
		expectedTypeName := fmt.Sprintf("%T", &EmailReigisterInfo{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})

	}
	if err := userInfo.Validate(); err != nil {
		return ErrInvalidUserObject
	}
	exist := r.exist(user.Name)
	if !exist {
		return ErrUserMustExist(user.Name)
	}
	collection := r.db.Collection("EMAIL")
	filter := bson.D{{Key: "Name", Value: user.Name}}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}
func (r *register) UpdateRegistration(ctx context.Context, userInfo Validator, subcategories []Tag) error {
	user, ok := userInfo.(*EmailReigisterInfo)
	if !ok {
		// This should return an error but since i havnt handled the Validate file yet this will just crash
		logger.Critical("Incorrect type passed into Email Notification Service")
		givenType := fmt.Sprintf("%T", userInfo)
		expectedTypeName := fmt.Sprintf("%T", &EmailReigisterInfo{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})

	}
	if err := userInfo.Validate(); err != nil {
		logger.Debug(fmt.Sprintf("validate EmailRegisterInfo struct %v\n", err))
		return ErrInvalidUserObject
	}
	exist := r.exist(user.Name)
	if !exist {
		return fmt.Errorf("user %s must already exist before atempty to update their registration ", user.Name) //ErrUserMustExist(user.Name)
		//return ErrUserMustExist(user.Name)
	}
	user.Tags = addifNotexist("EMAIL", user.Tags)
	fmt.Printf("Passed object thats going to be sent to mongodb %+v\n", userInfo)
	collection := r.db.Collection("EMAIL")
	filter := bson.D{{Key: "Name", Value: user.Name}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "Email", Value: user.Email},
			{Key: "Tags", Value: user.Tags},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("mongo failed to update registration: %w", err)
	}
	return nil

}
func (r *register) ListUsers(filter Filter) ([]string, error) {
	collection := r.db.Collection("EMAIL")

	// Retrieve all registered users from the EMAIL collection
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		logger.Critical(fmt.Sprintf("Failed to retrieve users: %v", err))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []EmailReigisterInfo
	if err = cursor.All(context.TODO(), &results); err != nil {
		logger.Critical(fmt.Sprintf("Failed to parse users: %v", err))
		return nil, err
	}

	var filteredUsers []string
	for _, user := range results {
		// filter based on tag
		if filter(context.TODO(), TagToString(user.Tags)) {
			filteredUsers = append(filteredUsers, user.Name)
		}
	}

	return filteredUsers, nil
}

func (m *Mailer) Start(ctx context.Context) error {
	const DevEmail = "richiebbaah@gmail.com"
	err := m.sendEmail(ctx, DevEmail, "Email Service is starting", "Hi Rich, Just a heads-up that the server is starting up now.")
	if err != nil {
		requestid := ctx.Value(pkg.RequestIDKey{}).(string)
		logger.Critical("Email service is not working. Basic request to owners email isnt working")
		return ErrInvalidEmailResponse(requestid, err)
	}
	return nil
}

func (m *Mailer) Notify(ctx context.Context, body Messenger, filter Filter) (int, []error) {
	requestID, ok := ctx.Value(pkg.RequestIDKey{}).(string)
	if !ok {
		logger.Critical("context request ID is not present when it should be!")
		return 0, []error{pkg.ErrMissingRequestID}
	}
	startTime, ok := ctx.Value(pkg.StartTime{}).(time.Time)
	if !ok {
		logger.Critical("context start time is not present when it should be!")
		return 0, []error{pkg.ErrMissingStartTime}
	}
	formattedTime := startTime.Format("Mon, 02 Jan 2006 15:04:05 MST")

	if err := m.Validate(body); err != nil {
		return 0, []error{ErrInvalidMessengerPassed(err)}
	}
	allUsers, err := m.ListUsers(filter)
	if err != nil {
		return 0, []error{err}
	}

	users, err := m.userbyID(ctx, allUsers) // Fetch all users since the filter will determine who to notify
	if err != nil {
		return 0, []error{fmt.Errorf("requestId: %s failed to retrieve users: %w", requestID, err)}
	}
	if len(users) == 0 {
		return 0, []error{errors.New("no users found to notify")}
	}

	var toNotify []EmailReigisterInfo
	for _, user := range users {
		if filter(ctx, TagToString(user.Tags)) {
			toNotify = append(toNotify, user)
		}
	}

	var errorArray []error
	var notified int

	finalMessage := fmt.Sprintf("email sender: %s\n    %s", body.Metadata().From(), body.Message().Content().(string))

	for _, user := range toNotify {
		err := m.sendEmail(ctx, user.Email, body.Metadata().Title(), finalMessage)
		if err != nil {
			errorArray = append(errorArray, err)
			logger.Debug(fmt.Sprintf("request ID: %s failed to send email to %s with error: %v", requestID, user.Email, err))
			continue
		}
		notified++
		logger.Info(fmt.Sprintf("requestId: %s email sent to %s with subject %s", requestID, user.Email, body.Metadata().Title()))
	}

	if len(toNotify) != notified {
		str := fmt.Sprintf("requestId: %s notify should have sent %d emails but only sent %d context startTime: %s",
			requestID, len(toNotify), notified, formattedTime)
		errorArray = append(errorArray, errors.New(str))
	}

	return notified, errorArray
}

func (m *Mailer) SendDirectMessage(ctx context.Context, body Messenger, from string, recipients []string) []error {
	// dont forget to not use
	requestID := ctx.Value(pkg.RequestIDKey{}).(string)
	if len(recipients) == 0 {
		return []error{errors.New("recipiant list cannot be empty")}
	}
	users, err := m.userbyID(ctx, recipients)
	if err != nil {
		return []error{err}
	}
	if len(users) == 0 {
		return []error{errors.New("no valid usersnames were passed in. users must already be registred to be notified")}
	}
	// cleaner formating
	finalMessage := fmt.Sprintf("email sender:%s\n    %s", from, body.Message().Content().(string))
	var errorArray []error
	for _, reciever := range users {
		err := m.sendEmail(ctx, reciever.Email, body.Metadata().Title(), finalMessage)
		if err != nil {
			logger.Debug(fmt.Sprintf("request ID: %s failed with error %v", requestID, err))
			errorArray = append(errorArray, err)
		}
	}
	return errorArray
}

func (m *Mailer) Validate(body Messenger) error {
	mail, ok := body.(*Mailbody)
	if !ok {
		givenType := fmt.Sprintf("%T", body)
		expectedTypeName := fmt.Sprintf("%T", &Mailbody{})
		return ErrInvalidUserType(givenType, []string{expectedTypeName})
	}
	if mail.Body == "" {
		return errors.New("email body cannot be empty")
	}
	if mail.Subject == "" {
		return errors.New("email subject cannot be empty")

	}
	// TODO: we need to have a way to validate that the recipiant of these emails exist
	if mail.To == "" {
		return errors.New("email must have a valid recipiant")

	}
	return nil
}

// works but still need to implement the rest of the interface, mabey use composition to have a seperate concerns for registration and actauly notifying
func (m *Mailer) SendMail(ctx context.Context, recipient, subject, body string) error {
	return m.sendEmail(ctx, recipient, subject, body)
}

// Function to send an email
func (m *Mailer) sendEmail(ctx context.Context, recipientEmail, subject, body string) error {
	// Set up authentication information.
	// intergrate with context later
	requestid, ok := ctx.Value(pkg.RequestIDKey{}).(string)
	if !ok {
		logger.Critical("context request ID is not present when it should be!")
		return pkg.ErrMissingRequestID
	}
	auth := smtp.PlainAuth("", m.senderEmail, m.senderPassword, "smtp.gmail.com")

	// Format the message
	msg := "From: " + m.senderEmail + "\n" +
		"To: " + recipientEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// Send the email
	err := smtp.SendMail(
		"smtp.gmail.com:587",     // SMTP server and port
		auth,                     // Authentication
		m.senderEmail,            // Sender email
		[]string{recipientEmail}, // Recipient email
		[]byte(msg),              // Message body as byte array
	)
	if err != nil {
		logger.Warn(fmt.Sprintf("request ID: %s  Error sending email: %v\n", requestid, err))
		return err
	}

	return nil
}

func (r *register) exist(name string) bool {
	collection := r.db.Collection("EMAIL")
	filter := bson.D{{Key: "Name", Value: name}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Printf("Error checking if user exists: %v", err)
		return false
	}
	logger.Info(fmt.Sprintf("Name:%s exist %d times in databse\n", name, count))
	return count > 0
}

func (m *Mailer) userbyID(ctx context.Context, toFind []string) ([]EmailReigisterInfo, error) {
	collection := m.db.Collection("EMAIL")

	// Correct filter for MongoDB query
	filter := bson.M{"Name": bson.M{"$in": toFind}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)
	var res []EmailReigisterInfo
	if err = cursor.All(ctx, &res); err != nil {
		return nil, fmt.Errorf("failed to decode cursor results: %w", err)
	}
	return res, nil
}
