package main

import (
	// this package needs to always to be run first b4 all other custom packages
	"log"

	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
	"github.com/Rich-T-kid/Notiffy/pkg"
)

func main() {
	ctx := pkg.ContextWithRequestID()
	sms := services.NewSMSNotification()
	if err := sms.Start(ctx); err != nil {
		log.Fatal(err)
	}
	print("no errors")

}
