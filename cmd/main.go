package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/Oakey12/myAPIGo/docs"
	"github.com/Oakey12/myAPIGo/internal/handlers"
	"github.com/Oakey12/myAPIGo/internal/structs"
	"github.com/Oakey12/myAPIGo/utils"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "modernc.org/sqlite"
)

func main() {

	db := initBD("./notes.db")
	defer db.Close()

	store := structs.NewNoteStore(db)

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

func initBD(filePath string) *sql.DB {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы %v", err)
	}

	log.Println("База данных успешно инициализирована и таблица готова")
	return db
}
