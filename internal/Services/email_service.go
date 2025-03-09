package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"
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

type EmailReigisterInfo struct {
	Name  string `bson:"Name"`
	Email string `bson:"Email"`
	Tags  tags   `bson:"Tags"`
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
	MailList tags
	Body     string
	To       string // Recipiant
}

// MessageMeta
func (m *Mailbody) Tags() tags       { return m.MailList }
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
func DefineMail(subject, body, to string, topics tags) *Mailbody {
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
		return ErrInvalidUserObject
	}
	exist := r.exist(user.Name)
	if !exist {
		return fmt.Errorf("user %s must already exist before atempty to update their registration ", user.Name) //ErrUserMustExist(user.Name)
		//return ErrUserMustExist(user.Name)
	}
	collection := r.db.Collection("EMAIL")
	filter := bson.D{{Key: "Name", Value: user.Name}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "Tags", Value: bson.M{"$in": subcategories}},
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
		if filter(context.TODO(), user.Name, user.Tags) {
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
	return 0, nil
}

func (m *Mailer) SendDirectMessage(ctx context.Context, body Messenger, from string, recipient []string) []error {
	return nil
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

// Validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

var validTLDs = map[string]bool{
	".com": true, ".net": true, ".org": true, ".edu": true,
	".gov": true, ".io": true, ".co": true, ".us": true,
}

// ValidateEmail performs a comprehensive email validation
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	if len(email) > 254 {
		return errors.New("email exceeds the maximum length of 254 characters")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("email does not match the required format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("email must contain a single @ character")
	}

	domain := parts[1]

	if len(domain) < 3 || len(domain) > 255 {
		return errors.New("domain part of the email is invalid")
	}

	// Validate the domain contains at least one dot
	if !strings.Contains(domain, ".") {
		return errors.New("domain must contain a dot (.)")
	}

	// Check if the domain has a valid TLD
	tld := strings.ToLower(domain[strings.LastIndex(domain, "."):])
	if !validTLDs[tld] {
		return fmt.Errorf("invalid top-level domain: %s", tld)
	}

	// Verify the domain has MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return fmt.Errorf("domain does not have valid MX records: %s", domain)
	}

	return nil
}
