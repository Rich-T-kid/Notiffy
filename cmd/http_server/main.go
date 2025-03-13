package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
)

var port = os.Getenv("HTTP_SERVER_PORT")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Status 200 OK")
	})
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}

}
