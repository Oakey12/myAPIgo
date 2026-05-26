package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()
	PingHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получен %v", status)
	}
}
