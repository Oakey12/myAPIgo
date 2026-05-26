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

	jsonBody := []byte(`{"title": "test", "content": "test"}`)

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

func TestDeleteNote(t *testing.T) {
	store := structs.NewNoteStore()
	handler := &NoteHandler{Store: store}

	store.SaveNote("Заметка под удаление", "Текст заметки")
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /notes/{id}", handler.DeleteNote)

	test := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Позитивный: Успешное удаление существующей заметки",
			url:            "/notes/1",
			expectedStatus: http.StatusOK, // Ожидаем 200
		},
		{
			name:           "Негативный: Попытка удалить несуществующий ID",
			url:            "/notes/999",
			expectedStatus: http.StatusNotFound, // Ожидаем 404
		},
		{
			name:           "Негативный: Передача текста вместо числового ID",
			url:            "/notes/abc",
			expectedStatus: http.StatusBadRequest, // Ожидаем 400
		},
		{
			name:           "Негативный: Передача отрицательного ID",
			url:            "/notes/-5",
			expectedStatus: http.StatusBadRequest, // Ожидаем 400
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tc.url, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Сценарий '%s' провален: ожидался статус %d, получен %d", tc.name, tc.expectedStatus, rr.Code)
			}
		})
	}
	remainingNotes := store.GetAllNotes()
	if len(remainingNotes) != 0 {
		t.Errorf("После успешного удаления в хранилище остались записи. Ожидалось 0, найдено %d", len(remainingNotes))
	}
}

func TestUpdateNote(t *testing.T) {
	store := structs.NewNoteStore()
	handler := &NoteHandler{Store: store}

	store.SaveNote("Старый заголовок", "Старый текст")

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /notes/{id}", handler.UpdateNote)

	tests := []struct {
		name           string
		url            string
		body           []byte
		expectedStatus int
	}{
		{
			name:           "Позитивный: успешное обновление",
			url:            "/notes/1",
			body:           []byte(`{"title": "Новый заголовок", "content": "Новый текст"}`),
			expectedStatus: http.StatusOK, // ожидаемый статус 200
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, ts.url, bytes.NewBuffer(ts.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != ts.expectedStatus {
				t.Errorf("Сценарий '%s' провален: ожидался статус %d, получен %d", ts.name, ts.expectedStatus, rr.Code)
			}
		})

	}
	updatedNote, ok := store.GetOneNote(1)
	if !ok {
		t.Fatalf("Заметка была удалена или не найдена")
	}

	if updatedNote.Title != "Новый заголовок" {
		t.Errorf("Поле Title не обновилось. Текущее значение: %s", updatedNote.Title)
	}

}

func TestPartialUpdateNote(t *testing.T) {
	store := structs.NewNoteStore()
	handler := &NoteHandler{Store: store}

	store.SaveNote("Тестовый заголовок", "Тестовый текст")

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /notes/{id}", handler.PatchNote)

	tests := []struct {
		name           string
		url            string
		body           []byte
		expectedStatus int
	}{
		{
			name:           "Позитивный: обновление только заголовка",
			url:            "/notes/1",
			body:           []byte(`{"title": "Новый заголовок"}`),
			expectedStatus: http.StatusOK, // ожидаемый статус 200
		},
		{
			name:           "Позитивный: обновление только текста",
			url:            "/notes/1",
			body:           []byte(`{"content": "Новый текст"}`),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Позитивный: Обновление обоих полей",
			url:            "/notes/1",
			body:           []byte(`{"title": "Новое имя", "content": "Новое описание"}`),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Негативный: Пустой запрос (оба поля nil)",
			url:            "/notes/1",
			body:           []byte(`{}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Негативный: Обновление несуществующей записи",
			url:            "/notes/999",
			body:           []byte(`{"title": "Тест"}`),
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, ts.url, bytes.NewBuffer(ts.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != ts.expectedStatus {
				t.Errorf("Сценарий '%s' провален: ожидался статус %d, получен %d", ts.name, ts.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetAllNote(t *testing.T) {
	store := structs.NewNoteStore()

	store.SaveNote("Первая заметка", "Текстовый 1")
	store.SaveNote("Вторая заметка", "Текстовый 2")

	handler := &NoteHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes", handler.GetAllNote)

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получен %v", status)
	}
}
