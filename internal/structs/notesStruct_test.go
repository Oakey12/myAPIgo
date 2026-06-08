package structs

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"testing"
)

func createTestBD(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Не удалось создать БД: %v", err)
	}
	createBDSQL := `CREATE TABLE notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	);`

	if _, err := db.Exec(createBDSQL); err != nil {
		t.Fatalf("Ошибка при создании тестовой таблицы: %v", err)
	}
	return db
}

func TestCreateNote(t *testing.T) {
	db := createTestBD(t)
	defer db.Close()

	store := NewNoteStore(db)

	title := "Тестовая заметка"
	content := "Тестовый контент"

	saveNote := store.CreateNote(title, content)

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
	db := createTestBD(t)
	defer db.Close()

	store := NewNoteStore(db)

	title := "Тестовый заголовок"
	content := "Тестовый контент"

	newNote := store.CreateNote(title, content) // новая заметка

	getNote, ok := store.GetNoteID(newNote.ID) // заметка, которую хотим получить
	if !ok {
		t.Errorf("Не удалось получить заметку с ID %d", newNote.ID)
	}
	if getNote.Title != newNote.Title {
		t.Errorf("Не совпадают заголовки")
	}
	if getNote.ID != newNote.ID {
		t.Errorf("ID заметки для получения совпадает с ID созданной заметки")
	}
	_, ok = store.GetNoteID(999)
	if ok {
		t.Errorf("Ожидалось, что заметка с ID 999 не будет найдена, но она нашлась")
	}

}

func TestGetAllNote(t *testing.T) {
	db := createTestBD(t)
	defer db.Close()

	store := NewNoteStore(db)

	store.CreateNote("Первая", "Текст 1")
	store.CreateNote("Вторая", "Текст 2")
	store.CreateNote("Третья", "Текст 3")

	notes := store.GetAllNotes()

	if len(notes) != 3 {
		t.Errorf("Ожидалось 3 заметки, а получили %d", len(notes))
	}
}

func TestDeleteNote(t *testing.T) {
	db := createTestBD(t)
	defer db.Close()

	store := NewNoteStore(db)

	newNote := store.CreateNote("Тестовый текст", "Тестовый текст")

	store.DeleteNoteID(newNote.ID)

	_, err := store.GetNoteID(newNote.ID)
	if err {
		t.Errorf("Заметка с ID %d всё ещё существует в базе после удаления", newNote.ID)
	}
	notes := store.GetAllNotes()
	if len(notes) != 0 {
		t.Errorf("Ожидалось 0 заметок в пустой базе, но найдено %d", len(notes))
	}
}

func TestUpdateNote(t *testing.T) {
	db := createTestBD(t)
	defer db.Close()

	store := NewNoteStore(db)

	newNote := store.CreateNote("Тестовый текст", "Тестовый текст")

	store.Update(newNote.ID, "Новый заголовок", "Новый контент")

	updateNote, ok := store.GetNoteID(newNote.ID)
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

func TestPatchNote(t *testing.T) {
	db := createTestBD(t)
	defer db.Close()

	store := NewNoteStore(db)

	firstNote := store.CreateNote("Старый заголовок", "Старый текст")

	titleUp := "Новый Заголовок"

	updateNote, ok := store.PatchNote(firstNote.ID, &titleUp, nil)
	if !ok {
		t.Fatalf("Заметка по ID: %d исчезла после обновления", firstNote.ID)
	}

	if updateNote.Title != titleUp {
		t.Errorf("Заголовок не обновился. Заголовок сейчас: %s", updateNote.Title)
	}
	if updateNote.Content != "Старый текст" {
		t.Errorf("Текст должен был остаться старым, но он изменился на '%s'", updateNote.Content)
	}

}
