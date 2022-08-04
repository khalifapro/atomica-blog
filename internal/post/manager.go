package post

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

// Manager ...
type Manager struct {
	DB DBConnector
}

// NewManager ...
func NewManager(db DBConnector) *Manager {
	return &Manager{DB: db}
}

// DBConnector contains dataAccess functionalities
type DBConnector interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// ListPosts ...
func (m *Manager) ListPosts() ([]BlogPost, error) {
	blogPosts := make([]BlogPost, 0)
	query := `SELECT id, title, content, updated_at, created_at, deleted_at, tags, photos FROM posts ORDER BY created_at ASC LIMIT $1 OFFSET $2`
	rows, err := m.DB.Query(query, 5, 0)
	if err != nil {
		return blogPosts, err
	}
	for rows.Next() {
		b := BlogPost{}
		if err = rows.Scan(&b.ID, &b.Title, &b.Content, &b.UpdatedAt, &b.CreatedAt, &b.DeletedAt, pq.Array(&b.Tags), pq.Array(&b.Photos)); err != nil {
			return blogPosts, err
		}
		blogPosts = append(blogPosts, b)
	}
	if rows.Err() != nil {
		return blogPosts, rows.Err()
	}

	return blogPosts, nil
}

// CreatePost ...
func (m *Manager) CreatePost(post *BlogPost) (int64, error) {
	id := int64(0)
	row := m.DB.QueryRow("INSERT INTO posts (title, content, photos, tags, created_at, updated_at) Values ($1, $2, $3, $4, now(), now()) RETURNING id", post.Title, post.Content, pq.Array(post.Photos), pq.Array(post.Tags))
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

// UpdatePost ...
func (m *Manager) UpdatePost(post *BlogPost, id int64) error {
	row := m.DB.QueryRow("UPDATE posts SET title=$1, content=$2, photos=$3, tags=$4, created_at=now(), updated_at=now() WHERE id=$5", post.Title, post.Content, pq.Array(post.Photos), pq.Array(post.Tags), id)
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

// GetPostByID ...
func (m *Manager) GetPostByID(id int64) (*BlogPost, error) {
	b := BlogPost{}
	row := m.DB.QueryRow("select id, title, content, updated_at, created_at, deleted_at, tags, photos from posts WHERE id=$1 AND deleted_at is NULL", id)
	if row.Err() != nil {
		return &b, row.Err()
	}

	if err := row.Scan(&b.ID, &b.Title, &b.Content, &b.UpdatedAt, &b.CreatedAt, &b.DeletedAt, pq.Array(&b.Tags), pq.Array(&b.Photos)); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return &b, err
	}
	return &b, nil
}

// DeletePostByID ...
func (m *Manager) DeletePostByID(id int64) error {
	res := m.DB.MustExec("UPDATE posts SET deleted_at=now() WHERE id=$1 AND deleted_at is NULL", id)
	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("error deleting post, post does not exist")
	}
	return nil
}
