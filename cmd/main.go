package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/Oakey12/myAPIGo/docs"
	"github.com/Oakey12/myAPIGo/internal/handlers"
	"github.com/Oakey12/myAPIGo/internal/structs"
	"github.com/Oakey12/myAPIGo/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	store := structs.NewNoteStore()

	noteHandler := &handlers.NoteHandler{Store: store}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handlers.PingHandler)
	mux.HandleFunc("GET /hello", handlers.HelloHandler)

	mux.HandleFunc("POST /notes", noteHandler.CreateNote)
	mux.HandleFunc("GET /notes/{id}", noteHandler.GetNoteID)
	mux.HandleFunc("GET /notes", noteHandler.GetAllNote)
	mux.HandleFunc("DELETE /notes/{id}", noteHandler.DeleteNote)
	mux.HandleFunc("PUT /notes/{id}", noteHandler.UpdateNote)
	mux.HandleFunc("PATCH /notes/{id}", noteHandler.PatchNote)

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	log.Println("Server running on http://localhost:8012")
	log.Println("Swagger documentation available on http://localhost:8012/swagger/index.html")

	fmt.Println("Server running!")

	handler := utils.Logger(mux)
	if err := http.ListenAndServe(":8012", handler); err != nil {
		log.Fatal(err)
	}
}
