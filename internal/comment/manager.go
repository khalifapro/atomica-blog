package comment

import (
	"database/sql"
	"errors"
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

// ListPostComments ...
func (m *Manager) ListPostComments(postID int64) ([]PostComment, error) {
	blogPosts := make([]PostComment, 0)
	query := `SELECT id, title, content, creator, updated_at, created_at, deleted_at FROM comments where post_id=$1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`
	rows, err := m.DB.Query(query, postID, 5, 0)
	if err != nil {
		return blogPosts, err
	}
	for rows.Next() {
		b := PostComment{}
		if err = rows.Scan(&b.ID, &b.Title, &b.Content, &b.UpdatedAt, &b.CreatedAt, &b.DeletedAt); err != nil {
			return blogPosts, err
		}
		blogPosts = append(blogPosts, b)
	}
	if rows.Err() != nil {
		return blogPosts, rows.Err()
	}

	return blogPosts, nil
}

// CreateComment ...
func (m *Manager) CreateComment(comment *PostComment) (int64, error) {
	id := int64(0)
	row := m.DB.QueryRow("INSERT INTO comments (title, content, creator, created_at, updated_at) Values ($1, $2, $3, now(), now()) RETURNING id", comment.Title, comment.Content, comment.Creator)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

// GetCommentByID ...
func (m *Manager) GetCommentByID(id int64) (*PostComment, error) {
	b := PostComment{}
	row := m.DB.QueryRow("select id, title, content, creator, updated_at, created_at, deleted_at, from comments WHERE id=$1 AND deleted_at is NULL", id)
	if row.Err() != nil {
		return &b, row.Err()
	}

	if err := row.Scan(&b.ID, &b.Title, &b.Content, &b.Creator, &b.UpdatedAt, &b.CreatedAt, &b.DeletedAt); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return &b, err
	}
	return &b, nil
}

// DeleteCommentByID ...
func (m *Manager) DeleteCommentByID(id int64) error {
	res := m.DB.MustExec("UPDATE comments SET deleted_at=now() WHERE id=$1 AND deleted_at is NULL", id)
	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("error deleting post, post does not exist")
	}
	return nil
}
