package comment

type PostComment struct {
	ID        int64   `db:"id"`
	Title     string  `db:"title"`
	Content   string  `db:"content"`
	Creator   string  `db:"creator"`
	CreatedAt string  `db:"created_at"`
	UpdatedAt string  `db:"updated_at"`
	DeletedAt *string `db:"deleted_at"`
}
