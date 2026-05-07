package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
	if err != nil {
		log.Printf("Ошибка при отправке JSON в /ping: %v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "noName"
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello " + name,
	})
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)

	mux.HandleFunc("GET /hello", helloHandler)

	

	handler := logger(mux)
	if err := http.ListenAndServe(":8012", handler); err != nil {
		log.Fatal(err)
	}
}
