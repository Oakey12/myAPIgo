package structs

import (
	"sync"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type NoteStore struct {
	mu     sync.RWMutex
	data   map[int]Note
	lastId int
}

func NewNoteSore() *NoteStore {
	return &NoteStore{
		data: make(map[int]Note),
	}
}

// Метод для добавления новой заметки
func (ns *NoteStore) SaveNote(title, content string) Note {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.lastId++

	note := Note{
		ID:      ns.lastId,
		Title:   title,
		Content: content,
	}
	ns.data[note.ID] = note
	return note
}

// метод для получения одной заметки по ID
func (ns *NoteStore) GetOneNote(id int) (Note, bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	note, exists := ns.data[id]
	return note, exists
}

// Метод для получения всех заметок
func (ns *NoteStore) GetAllNotes() []Note {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	notes := make([]Note, 0, len(ns.data))

	for _, v := range ns.data {
		notes = append(notes, v)
	}
	return notes
}

// Удаление заметки
func (ns *NoteStore) Delete(id int) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	_, exists := ns.data[id]
	if exists {
		delete(ns.data, id)
	}
	return exists
}
