package main

import (
	"fmt"

	_ "github.com/Rich-T-kid/Notiffy/enviroment" // this package needs to always to be run first b4 all other custom packages
)

var (
	port = "9999"
	env  = "dev"
)

func main() {

	fmt.Println("vim-go")
}
