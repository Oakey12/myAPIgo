package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Oakey12/myAPIGo/internal/handlers"
	"github.com/Oakey12/myAPIGo/internal/structs"
	"github.com/Oakey12/myAPIGo/utils"
)

func main() {

	store := structs.NewNoteSore()

	noteHandler := &handlers.NoteHandler{Store: store}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handlers.PingHandler)
	mux.HandleFunc("GET /hello", handlers.HelloHandler)

	mux.HandleFunc("POST /notes", noteHandler.CreateNote)
	mux.HandleFunc("GET /notes/{id}", noteHandler.GetNoteID)
	mux.HandleFunc("GET /notes", noteHandler.GetAllNote)
	mux.HandleFunc("DELETE /notes{id}", noteHandler.DeleteNote)

	fmt.Println("Server running!")

	handler := utils.Logger(mux)
	if err := http.ListenAndServe(":8012", handler); err != nil {
		log.Fatal(err)
	}
}
