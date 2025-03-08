package main

import (
	// this package needs to always to be run first b4 all other custom packages
	"fmt"

	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
)

func main() {
	var text = "etry is a type of literature typically written in verse that uses figurative language, or language that can have different meanings from what is literally said, to give multiple shades of meaning to a word or a phrase. Examples o"
	//	ctx := pkg.ContextWithRequestID()
	//msg := services.DefineMessage("richard", "king", "textMessage")
	email := services.DefineMail("Test number 402", text, "richiebbaah@gmail.com") //services.Mailbody{
	m := services.NewMailer()
	fmt.Println("errors: ", m.Validate(email))
	//fmt.Println(email)
	println("No errors occured")
}
