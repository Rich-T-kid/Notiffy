package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"

	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment"
	"github.com/Rich-T-kid/Notiffy/pkg"
)

// Generate unique timestamps to avoid duplicates
func generateTimestamp() string {
	return time.Now().Format("20060102150405") // YYYYMMDDHHMMSS format
}

// Generate a random phone number// Generate unique email with timestamp
func generatePhoneNumber() int64 {
	min := int64(9000000000)            // Minimum 10-digit number
	max := int64(9999999999)            // Maximum 10-digit number
	return min + rand.Int63n(max-min+1) // Ensures 10-digit number
}
func generateEmail(name string) string {
	domains := []string{"gmail.com", "yahoo.com", "outlook.com", "protonmail.com"}
	timestamp := generateTimestamp()
	return fmt.Sprintf("%s_%s_%s@%s", name, uuid.New().String()[:4], timestamp, domains[rand.Intn(len(domains))])
}

// Generate random tags
var testTags = []services.Tag{"SMS", "Email", "Sports", "Dance", "News", "Tech", "Entertainment"}

// List of test names
var testNames = []string{
	"JohnDoe", "JaneSmith", "AliceBrown", "BobMiller", "CharlieDavis",
	"EmmaWilson", "LiamAnderson", "OliviaThomas", "NoahMartinez", "AvaHernandez",
	"WilliamLopez", "SophiaGonzalez",
}

// Register test users using the Notification Service interfaces
func registerTestUsers() {
	ctx := pkg.ContextWithRequestID()

	// Initialize Services
	smsService := services.NewSMSUserService()
	emailService := services.NewMailRegister()

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Loop through test users
	for _, name := range testNames {
		timestamp := generateTimestamp()
		phone := generatePhoneNumber()
		email := generateEmail(name)

		// Assign 2-3 random tags
		selectedTags := []services.Tag{testTags[rand.Intn(len(testTags))], testTags[rand.Intn(len(testTags))]}

		// **Register SMS User**
		smsUser := &services.RegisterINFO{
			Name:    fmt.Sprintf("%s_%s", name, timestamp), // Unique name
			Contact: phone,
			Tags:    selectedTags,
		}
		err := smsService.Register(ctx, smsUser, selectedTags)
		if err != nil {
			log.Printf("Failed to register SMS user %s: %v", smsUser.Name, err)
		} else {
			log.Printf("Successfully registered SMS user: %s", smsUser.Name)
		}

		// **Register Email User**
		emailUser := &services.EmailReigisterInfo{
			Name:  fmt.Sprintf("%s_%s", name, timestamp), // Unique name
			Email: email,
			Tags:  selectedTags,
		}
		err = emailService.Register(ctx, emailUser, selectedTags)
		if err != nil {
			log.Printf("Failed to register Email user %s: %v", emailUser.Name, err)
		} else {
			log.Printf("Successfully registered Email user: %s", emailUser.Name)
		}
	}
}

func main() {
	registerTestUsers()
}
