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

type EmailReigisterInfo struct {
	Name  string
	Email string
	Tags  tags
}

func (e *EmailReigisterInfo) Validate() error {
	if e.Name == "" {
		return errors.New("user's name cannot be empty")
	}
	if e.Email == "" {
		return errors.New("user's Email cannot be empty")
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
func DefineMail(subject, body, to string) *Mailbody {
	return &Mailbody{
		Subject: subject,
		Body:    body,
		To:      to,
	}
}
func (r *register) Register(ctx context.Context, userinfo Validator, subcategory []Tag) error {
	return nil
}
func (r *register) Unregister(ctx context.Context, userinfo Validator, subcategory []Tag) error {
	return nil
}
func (r *register) UpdateRegistration(ctx context.Context, userinfo Validator, subcategory []Tag) error {
	return nil
}
func (r *register) ListUsers(filter Filter) ([]string, error) {
	return nil, nil
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

func (m *Mailer) exist(name string) bool {
	collection := m.db.Collection("EMAIL")
	filter := bson.D{{Key: "Name", Value: name}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Printf("Error checking if user exists: %v", err)
		return false
	}
	logger.Info(fmt.Sprintf("Name:%s exist %d times in databse\n", name, count))
	return count > 0
}
