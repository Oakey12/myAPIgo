package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
	if err != nil {
		log.Printf("Ошибка при отправке JSON в /ping: %v", err)
	}
}
