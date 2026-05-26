package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", HelloHandler)

	req := httptest.NewRequest(http.MethodGet, "/hello/Name", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получен %v", status)
	}
}
