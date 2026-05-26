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

func TestGetAllNote(t *testing.T) {
	store := NewNoteStore()

	store.SaveNote("Первая", "Текст 1")
	store.SaveNote("Вторая", "Текст 2")
	store.SaveNote("Третья", "Текст 3")

	notes := store.GetAllNotes()

	if len(notes) != 3 {
		t.Errorf("Ожидалось 3 заметки, а получили %d", len(notes))
	}
}

func TestDeleteNote(t *testing.T) {
	store := NewNoteStore()

	newNote := store.SaveNote("Тестовый текст", "Тестовый текст")

	store.Delete(newNote.ID)

	_, err := store.GetOneNote(newNote.ID)
	if err {
		t.Errorf("Заметка с ID %d всё ещё существует в базе после удаления", newNote.ID)
	}
	notes := store.GetAllNotes()
	if len(notes) != 0 {
		t.Errorf("Ожидалось 0 заметок в пустой базе, но найдено %d", len(notes))
	}
}

func TestUpdateNote(t *testing.T) {
	store := NewNoteStore()

	newNote := store.SaveNote("Тестовый текст", "Тестовый текст")

	store.Update(newNote.ID, "Новый заголовок", "Новый контент")

	updateNote, ok := store.GetOneNote(newNote.ID)
	if !ok {
		t.Fatalf("Заметка по ID: %d исчезла после обновления", newNote.ID)
	}
	if updateNote.Title != "Новый заголовок" {
		t.Errorf("Заголовок не обновился! Ожидалось 'Новый заголовок', получено '%s'", updateNote.Title)
	}
	if updateNote.Content != "Новый контент" {
		t.Errorf("Текст не обновился! Ожидалось 'Новый контент', получено '%s'", updateNote.Content)
	}
}
