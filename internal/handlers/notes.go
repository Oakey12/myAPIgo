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

// POST для создания записи
func (nh *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	createNote := nh.Store.SaveNote(request.Title, request.Content)
	w.Header().Set("Content-Type", "applicatin/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createNote)
}

// GET получение заметки по ID
func (nh *NoteHandler) GetNoteID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	note, ok := nh.Store.GetOneNote(id)
	if !ok {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// GET получение всех заметок
func (nh *NoteHandler) GetAllNote(w http.ResponseWriter, r *http.Request) {
	notes := nh.Store.GetAllNotes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// DELETE заметку по id
func (nh *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	deleted := nh.Store.Delete(id)
	if !deleted {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// PUT Полное обновление заметки
func (nh *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
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
