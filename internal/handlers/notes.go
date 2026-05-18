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

// POST для создания записи
// CreateNote для создания записи
// @Summary      Создать новую заметку
// @Description  Принимает заголовок и текст, генерирует ID и сохраняет в память
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        input body CreateNoteRequest true "Данные заметки"
// @Success      201  {object}  structs.Note
// @Failure      400  {string}  string "Неверный формат JSON"
// @Router       /notes [post]
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

	createNote := nh.Store.SaveNote(request.Title, request.Content)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createNote)
}

// GET получение заметки по ID
// GetNoteID получение заметки по ID
// @Summary      Получить заметку по ID
// @Description  Возвращает одну заметку по её числовому идентификатору
// @Tags         notes
// @Produce      json
// @Param        id   path      int  true  "Идентификатор заметки"
// @Success      200  {object}  structs.Note
// @Failure      400  {string}  string "Неверный ID"
// @Failure      404  {string}  string "Заметка не найдена"
// @Router       /notes/{id} [get]
func (nh *NoteHandler) GetNoteID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
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
// GetAllNote получение всех заметок
// @Summary      Получить все заметки
// @Description  Возвращает массив всех сохраненных заметок
// @Tags         notes
// @Produce      json
// @Success      200  {array}   structs.Note
// @Router       /notes [get]
func (nh *NoteHandler) GetAllNote(w http.ResponseWriter, r *http.Request) {
	notes := nh.Store.GetAllNotes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// DELETE заметку по id
// DeleteNote заметку по id
// @Summary      Удалить заметку
// @Description  Удаляет заметку из памяти по её ID
// @Tags         notes
// @Produce      json
// @Param        id   path      int  true  "Идентификатор заметки"
// @Success      200  {object}  map[string]string "status: ok"
// @Failure      400  {string}  string "Неверный ID"
// @Failure      404  {string}  string "Заметка не найдена"
// @Router       /notes/{id} [delete]
func (nh *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
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
// UpdateNote Полное обновление заметки
// @Summary      Полностью обновить заметку (PUT)
// @Description  Заменяет все поля существующей заметки новыми данными. Требует передачи и title, и content.
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "Идентификатор заметки"
// @Param        input body      UpdateNoteRequest  true  "Новые данные заметки"
// @Success      200   {object}  structs.Note
// @Failure      400   {string}  string "Неверный ID или формат JSON"
// @Failure      404   {string}  string "Заметка не найдена"
// @Router       /notes/{id} [put]
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

// PATCH
// PatchNote Частичное обновление заметки
// @Summary      Частично обновить заметку (PATCH)
// @Description  Обновляет только переданные поля заметки (title и/или content). Пропущенные поля остаются без изменений.
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "Идентификатор заметки"
// @Param        input body      PatchNoteRequest  true  "Поля для изменения"
// @Success      200   {object}  structs.Note
// @Failure      400   {string}  string "Неверный ID, формат JSON или пустой запрос"
// @Router       /notes/{id} [patch]
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
