package main

import (
	"fmt"
	// this package needs to always to be run first b4 all other custom packages

	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
	Logger "github.com/Rich-T-kid/Notiffy/internal/log"   // this package needs to always to be run first b4 all other custom packages
)

var (
	port = "9999"
	env  = "dev"
)

func main() {
	Logger.Info("should output to terminal")

	Logger.Critical("Should output to file")
	fmt.Println("vim-go")
}
