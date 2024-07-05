package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define SnippetModel which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This function will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	res, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, nil
	}

	// Get newly inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(query, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}

	// Convert the raw output from SQL to GO types.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Encapsulate the model by return `ErrNoRecord` instead return `sql.ErrNoRows`.
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Important
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		// Create a pointer to a new zeored Snippet struct
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// `rows.Err()` to retrieve any error that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
