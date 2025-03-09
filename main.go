package main

import (
	// this package needs to always to be run first b4 all other custom packages
	"fmt"
	"net/http"

	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Status 200 OK")
	})

	println("No errors occured")
	fmt.Println("Server is running on http://localhost:9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}

}
