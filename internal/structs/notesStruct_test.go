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

func TestGetOneNote(t *testing.T) {
	store := NewNoteStore()

	title := "Тестовый заголовок"
	content := "Тестовый контент"

	newNote := store.SaveNote(title, content) // новая заметка

	getNote, ok := store.GetOneNote(newNote.ID) // заметка, которую хотим получить
	if !ok {
		t.Errorf("Не удалось получить заметку с ID %d", newNote.ID)
	}
	if getNote.Title != newNote.Title {
		t.Errorf("Не совпадают заголовки")
	}
	if getNote.ID != newNote.ID {
		t.Errorf("ID заметки для получения совпадает с ID созданной заметки")
	}
	_, ok = store.GetOneNote(999)
	if ok {
		t.Errorf("Ожидалось, что заметка с ID 999 не будет найдена, но она нашлась")
	}

}
