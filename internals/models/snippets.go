package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// insert a new snippet into the db
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// create a sql query with placeholders (?) for user input data
	query := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`
	// call query with params using db exec
	result, err := m.DB.Exec(query, title, content, expires)
	if err!=nil {
		return 0, err
	}
	// get id of last inserted record
	id, err := result.LastInsertId()
	if err!=nil {
		return 0, err
	}
	return int(id), nil
}

// return a specific snippet based on id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// create sql with placeholders
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE
			expires > UTC_TIMESTAMP() AND
			id = ?
	`
	// QueryRow: returns the first row and ignores rest
	row := m.DB.QueryRow(query, id)
	// map sql data to go struct
	s := &Snippet{}
	err := row.Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.Created,
		&s.Expires,
	)
	if err!=nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

// return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// create sql query
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		ORDER BY id DESC LIMIT 10
	`
	// get rows from db using query
	rows, err := m.DB.Query(query)
	if err!=nil {
		return nil, err
	}
	defer rows.Close()
	// map sql rows to array of struct
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Content,
			&s.Created,
			&s.Expires,
		)
		if err!=nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err !=nil {
		return nil, err
	}
	return snippets, nil
}
