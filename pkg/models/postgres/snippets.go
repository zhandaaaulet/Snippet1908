package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"se08.com/pkg/models"
	"strconv"
	"time"
)

const (
	insertSql                 = "INSERT INTO snippets (title,content,created,expires) VALUES ($1,$2,$3,$4) RETURNING id"
	getSnippetById            = "SELECT id, title, content, created, expires FROM snippets where id=$1 AND expires > now()"
	getLastTenCreatedSnippets = "SELECT id, title, content, created, expires FROM snippets WHERE expires > now() ORDER BY created DESC LIMIT 10"
)

type SnippetModel struct {
	Pool *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	var id uint64
	interval, err := strconv.Atoi(expires)
	row := m.Pool.QueryRow(context.Background(), insertSql,
		title, content, time.Now(),
		time.Now().AddDate(0, 0, interval))
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	err := m.Pool.QueryRow(context.Background(), getSnippetById, id).
		Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own models.ErrNoRecord error
		// instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	rows, err := m.Pool.Query(context.Background(), getLastTenCreatedSnippets)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// Initialize an empty slice to hold the models.Snippets objects.
	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err

	}

	return snippets, nil
}
