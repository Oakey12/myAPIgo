package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "имя не введено"
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello " + name,
	})
	if err != nil {
		fmt.Printf("Ошибка при кодировании JSON: %v", err)
	}
}
