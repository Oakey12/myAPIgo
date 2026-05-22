package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Oakey12/myAPIGo/internal/structs"
)

func TestCreateNote(t *testing.T) {
	store := structs.NewNoteStore()
	handler := &NoteHandler{Store: store}

	jsonBody := []byte(`{"title": "Учеба", "content": "Повторить Go"}`)

	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateNote(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус 201 Created, получен %v", status)
	}

	if len(store.GetAllNotes()) != 1 {
		t.Errorf("Заметка не сохранилась в базу данных через хэндлер")
	}
}

func TestGetNoteID(t *testing.T) {
	store := structs.NewNoteStore()
	handler := &NoteHandler{Store: store}

	store.SaveNote("Тестовый заголовок", "Тестовый контент")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes/{id}", handler.GetNoteID)

	test := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Позитивный: Успешное получение существующей заметки",
			url:            "/notes/1",
			expectedStatus: http.StatusOK, // Ожидаем 200
		},
		{
			name:           "Негативный: Запрос несуществующего ID",
			url:            "/notes/999",
			expectedStatus: http.StatusNotFound, // Ожидаем 404
		},
		{
			name:           "Негативный: Передача букв вместо ID",
			url:            "/notes/abc",
			expectedStatus: http.StatusBadRequest, // Ожидаем 400
		},
		{
			name:           "Негативный: Передача отрицательного ID",
			url:            "/notes/-5",
			expectedStatus: http.StatusBadRequest, // Ожидаем 400
		},
	}
	for _, ts := range test {
		t.Run(ts.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ts.url, nil)

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != ts.expectedStatus {
				t.Errorf("Сценарий '%s' провален: ожидался статус %d, получен %d", ts.name, ts.expectedStatus, rr.Code)
			}
		})
	}
}
