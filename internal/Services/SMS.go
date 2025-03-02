package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	db "github.com/Rich-T-kid/Notiffy/internal/database" // this package needs to always to be run first b4 all other custom packages
	logger "github.com/Rich-T-kid/Notiffy/internal/log"
)

var (
	ErrInvalidUserType   = errors.New("incorrect type passed into sms notification service")
	ErrInvalidUserObject = errors.New("invalid user object passed in")
	ErrUsernameExists    = errors.New("username already exists, pick a new one")
	ErrUserInfoError     = errors.New("error in your userInfo")
	ErrNotImplemented    = errors.New("TODO: Implement")
	ErrUserMustExist     = func(userName string) error {
		return fmt.Errorf("user: %s must exist before you deregister them", userName)
	}
)

func init() {
	//print("richard make sure that the third party API that you going to intergrate with is up and running in this init func\n")
}

type RegisterINFO struct {
	Name    string  `bson:"Name"`
	Contact float32 `bson:"Contact"` // phone number
	Tags    []Tag   `bson:"Tags"`
}

func (r *RegisterINFO) Validate(interface{}) (bool, error) {
	// Name cannot already be taken
	// Phone number must be valid
	return true, nil
}

type SMSNotification struct {
	apiKey string
	saver  db.Storage
	db     *mongo.Database // store as mongo collection for now
}

func NewSMSNotification() *SMSNotification {
	return &SMSNotification{
		apiKey: os.Getenv("TEXTBELT_API_KEY"),
		saver:  db.NewStorage(),
		db:     db.EstablishMongoConnection(),
	}
}

// Completed the UserService interface for SMS
func NewSMSUserService() UserService {
	return NewSMSNotification()
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

// add to db.SMS collection. for now these all go under the same Mangager(person users register under for notification) but when we build an external api we should take that in
func (s *SMSNotification) Register(ctx context.Context, userInfo Validator, subcategory []Tag) error {
	user, ok := userInfo.(*RegisterINFO)
	if !ok {
		// This should return an error but since i havnt handled the Validate file yet this will just crash
		logger.Critical("Incorrect type passed into SMS Notification Service")
		return ErrInvalidUserType

	}
	collection := s.db.Collection("SMS")
	if _, err := userInfo.Validate(userInfo); err != nil {
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
		return ErrInvalidUserType

	}
	if _, err := userInfo.Validate(userInfo); err != nil {
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
		return ErrInvalidUserType
	}
	if _, err := userInfo.Validate(userInfo); err != nil {
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
