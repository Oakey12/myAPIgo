package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Oakey12/myAPIGo/internal/structs"
)

type NoteHandler struct {
	Store *structs.NoteStore
}

type CreateNoteRequest struct {
	Title   string `json:"title" example:"Купить продукты"`
	Content string `json:"content" example:"Молоко, хлеб, сыр"`
}
type UpdateNoteRequest struct {
	Title   string `json:"title" example:"Обновленный заголовок"`
	Content string `json:"content" example:"Обновленный текст заметки"`
}
type PatchNoteRequest struct {
	Title   *string `json:"title" example:"Новое название"`
	Content *string `json:"content" example:"Новый контент заметки"`
}

// Создание заметки
func (nh *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var request CreateNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	if request.Title == "" {
		http.Error(w, "Поле должно быть заполнено", http.StatusBadRequest)
		return
	}

	createNote := nh.Store.CreateNote(request.Title, request.Content)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createNote)
}

// Получение заметки по ID
func (nh *NoteHandler) GetNoteID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	note, ok := nh.Store.GetNoteID(id)
	if !ok {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (nh *NoteHandler) GetAllNote(w http.ResponseWriter, r *http.Request) {
	notes := nh.Store.GetAllNotes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// DELETE note (id)
func (nh *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	deleted := nh.Store.DeleteNoteID(id)
	if !deleted {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// PUT note
func (nh *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	updateNote, ok := nh.Store.Update(id, req.Title, req.Content)
	if !ok {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateNote)
}

// PATCH note
func (nh *NoteHandler) PatchNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	var req PatchNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.Title == nil && req.Content == nil {
		http.Error(w, "Хотя бы одно поле должно быть заполнено", http.StatusBadRequest)
		return
	}
	updateNote, ok := nh.Store.PatchNote(id, req.Title, req.Content)
	if !ok {
		http.Error(w, "Не удалось найти заметку", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateNote)

}
