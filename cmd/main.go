package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Oakey12/myAPIGo/internal/handlers"
	"github.com/Oakey12/myAPIGo/utils"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handlers.PingHandler)

	mux.HandleFunc("GET /hello", handlers.HelloHandler)

	fmt.Println("Server running!")

	handler := utils.Logger(mux)
	if err := http.ListenAndServe(":8012", handler); err != nil {
		log.Fatal(err)
	}
}
