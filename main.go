package main

import (
	"context"
	"fmt"

	// this package needs to always to be run first b4 all other custom packages

	_ "github.com/Rich-T-kid/Notiffy/internal/Services" // this package needs to always to be run first b4 all other custom packages
	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
)

var (
	port = "9999"
	env  = "dev"
)
var text = `On the Grand Line, where the oceans roar,  
Luffy sails, always craving more.  
With a crew of dreams and a heart of fire,  
The One Piece treasure fuels his desire.  
Through storms and foes, their bond won't cease,  
For at the end lies freedom and peace.`

func main() {

	ctx := context.Background()

	// Instantiate the SMS notification service
	smsService := services.NewSMSNotification()

	// Register some users
	user1 := &services.RegisterINFO{
		Name:    "JohnDoe",
		Contact: 1234567890,
	}
	user2 := &services.RegisterINFO{
		Name:    "JaneSmith",
		Contact: 2345678901,
	}
	user3 := &services.RegisterINFO{
		Name:    "FrankWhite",
		Contact: 3456789012,
	}

	// Register them
	err := smsService.Register(ctx, user1, []services.Tag{"promo"})
	if err != nil {
		fmt.Printf("Could not register user1: %v\n", err)
	}
	err = smsService.Register(ctx, user2, []services.Tag{"alerts"})
	if err != nil {
		fmt.Printf("Could not register user2: %v\n", err)
	}
	err = smsService.Register(ctx, user3, []services.Tag{})
	if err != nil {
		fmt.Printf("Could not register user3: %v\n", err)
	}

	// Define multiple messages from different authors
	msg1 := services.DefineMessage("Alice", "Welcome!", "Hello everyone, thanks for joining!")
	msg2 := services.DefineMessage("Bob", "Weekly Update", "Here is the weekly update...")
	//msg3 := services.DefineMessage("Charlie", "Urgent Notice", "Please read this important notice ASAP.")

	// Start the SMS service (checks quota, etc.)
	startErr := smsService.Start(ctx)
	if startErr != nil {
		fmt.Printf("Start error: %v\n", startErr)
	}

	// Example filter: Only notify those whose name starts with "J"
	filterByJ := func(ctx context.Context, tags ...interface{}) bool {
		return true
	}

	// Notify (broadcast) with msg1, filtered by "name starts with J"
	notifiedCount, notifyErrors := smsService.Notify(ctx, msg1, filterByJ)
	fmt.Printf("Notify broadcast result: %d notified\n", notifiedCount)
	for _, ne := range notifyErrors {
		fmt.Printf("Notify error: %v\n", ne)
	}

	// Direct message user2 and user3 with msg2
	dmErrs := smsService.SendDirectMessage(ctx, msg2, "Bob", []string{"JaneSmith", "FrankWhite"})
	if len(dmErrs) > 0 {
		fmt.Println("Direct message errors:")
		for _, e := range dmErrs {
			fmt.Printf("  %v\n", e)
		}
	}

}
