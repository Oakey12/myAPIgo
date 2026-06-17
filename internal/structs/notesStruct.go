package structs

import (
	"database/sql"
	"log"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type NoteStore struct {
	db *sql.DB
}

func NewNoteStore(db *sql.DB) *NoteStore {
	return &NoteStore{db: db}
}

// Метод для добавления новой заметки
func (ns *NoteStore) CreateNote(title, content string) (Note, error) {
	query := "INSERT INTO notes (title, content) VALUES (?, ?)"

	result, err := ns.db.Exec(query, title, content)
	if err != nil {
		log.Println("Ошибка записи в БД", err)
		return Note{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Ошибка получения LastInsertId:", err)
		return Note{}, err
	}

	return Note{
		ID:      int(id),
		Title:   title,
		Content: content,
	}, nil
}

// метод для получения одной заметки по ID
func (ns *NoteStore) GetNoteID(id int) (Note, bool) {
	query := "SELECT id, title, content FROM notes WHERE id = ?"
	row := ns.db.QueryRow(query, id)

	var note Note

	err := row.Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, false
		}
		log.Println("Ошибка чтения из БД", err)
		return Note{}, false
	}
	return note, true
}

// Метод для получения всех заметок
func (ns *NoteStore) GetAllNotes() []Note {
	query := "SELECT id, title, content FROM notes"

	rows, err := ns.db.Query(query)
	if err != nil {
		log.Println("Ошибка запроса всех заметок:", err)
		return []Note{}
	}
	defer rows.Close()
	var notes []Note

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			continue
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		log.Println("Критическая ошибка во время чтения строк:", err)
		return []Note{}
	}

	if notes == nil {
		notes = []Note{}
	}
	return notes

}

// Удаление заметки
func (ns *NoteStore) DeleteNoteID(id int) bool {
	query := "DELETE FROM notes WHERE id = ?"
	result, err := ns.db.Exec(query, id)
	if err != nil {
		log.Println("Ошибка удаления", err)
		return false
	}
	rowsDel, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return rowsDel > 0

}

// метод для полного изменения заметки PUT
func (ns *NoteStore) Update(id int, title, content string) (Note, bool) {
	query := "UPDATE notes SET title = ?, content = ? WHERE id = ?"
	result, err := ns.db.Exec(query, title, content, id)
	if err != nil {
		log.Println("Ошибка обновления", err)
		return Note{}, false
	}
	rowsUp, err := result.RowsAffected()
	if err != nil || rowsUp == 0 {
		return Note{}, false
	}
	return Note{ID: id, Title: title, Content: content}, true

}

// PAtch обновление записи по конкретным параметрам
func (ns *NoteStore) PatchNote(id int, title *string, content *string) (Note, bool) {
	currentNote, ok := ns.GetNoteID(id)
	if !ok {
		return Note{}, false
	}
	if title != nil {
		currentNote.Title = *title
	}
	if content != nil {
		currentNote.Content = *content
	}
	return ns.Update(id, currentNote.Title, currentNote.Content)
}
