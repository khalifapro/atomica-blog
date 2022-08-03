package post

type BlogPost struct {
	ID        int64    `db:"id"`
	Title     string   `db:"title"`
	Content   string   `db:"content"`
	Photos    []string `db:"photos"`
	Tags      []string `db:"tags"`
	CreatedAt string   `db:"created_at"`
	UpdatedAt string   `db:"updated_at"`
	DeletedAt *string  `db:"deleted_at"`
}
