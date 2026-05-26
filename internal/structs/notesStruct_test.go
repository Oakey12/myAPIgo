package structs

import (
	"testing"
)

func TestSaveNote(t *testing.T) {
	store := NewNoteStore()

	title := "Тестовая заметка"
	content := "Тестовый контент"

	saveNote := store.SaveNote(title, content)

	if saveNote.ID != 1 {
		t.Errorf("Ожидался ID 1, получен %d", saveNote.ID)
	}

	if saveNote.Title != title {
		t.Errorf("Ожидался заголовок '%s', получен '%s'", title, saveNote.Title)
	}

	notes := store.GetAllNotes()
	if len(notes) != 1 {
		t.Errorf("Ожидалась 1 заметка в хранилище, найдено %d", len(notes))
	}

}
