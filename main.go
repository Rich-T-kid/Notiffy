package main

import (
	"context"
	"fmt"

	// this package needs to always to be run first b4 all other custom packages

	_ "github.com/Rich-T-kid/Notiffy/internal/Services" // this package needs to always to be run first b4 all other custom packages
	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
	"github.com/Rich-T-kid/Notiffy/internal/log"
)

var (
	port = "9999"
	env  = "dev"
)

func main() {
	SMS := services.NewSMSNotification()
	userInfo := &services.RegisterINFO{
		Name:    "Test2",
		Contact: 9239592375,
		Tags:    []services.Tag{services.TagSMS, services.TagEmail},
	}

	err := SMS.Unregister(context.TODO(), userInfo, []services.Tag{"Email"})
	if err != nil {
		log.Critical(fmt.Sprintf("issue registering userid: Richard with SMS Notification service %v", err))
	}
	fmt.Println("errors: ", err)

}
