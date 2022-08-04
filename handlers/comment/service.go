package comments

import (
	"github.com/asaberwd/atomica-blog/internal/comment"
	"github.com/asaberwd/atomica-blog/swagger/models"
	"time"
)

type CommentHandler struct {
	CommentManager Service
}

func New(manager Service) *CommentHandler {
	return &CommentHandler{
		CommentManager: manager,
	}
}

type Service interface {
	CreateComment(comment *comment.PostComment) (int64, error)
	ListPostComments(postID int64) ([]comment.PostComment, error)
	GetCommentByID(id int64) (*comment.PostComment, error)
	DeleteCommentByID(id int64) error
}

func (c *CommentHandler) CreateComment(comment *models.Comment) (*models.Comment, error) {
	pComment := mapModelToComment(comment)
	id, err := c.CommentManager.CreateComment(pComment)
	if err != nil {
		return comment, err
	}

	postComment, err := c.GetCommentByID(id)
	if err != nil {
		return comment, err
	}

	return postComment, nil
}

func (c *CommentHandler) ListPostComments(postID int64) (models.Comments, error) {
	commentsModel := models.Comments{}
	blogPosts, err := c.CommentManager.ListPostComments(postID)
	if err != nil {
		return models.Comments{}, err
	}
	for _, b := range blogPosts {
		commentsModel = append(commentsModel, mapCommentToModel(&b))
	}

	return commentsModel, nil
}

func (c *CommentHandler) UpdatePostComment(id int64) (*models.Comment, error) {
	postComment, err := c.CommentManager.GetCommentByID(id)
	if err != nil {
		return &models.Comment{}, err
	}
	if postComment == nil {
		return nil, nil
	}
	postModel := mapCommentToModel(postComment)
	return postModel, nil
}

func (c *CommentHandler) GetCommentByID(id int64) (*models.Comment, error) {
	postComment, err := c.CommentManager.GetCommentByID(id)
	if err != nil {
		return &models.Comment{}, err
	}
	if postComment == nil {
		return nil, nil
	}
	commentModel := mapCommentToModel(postComment)
	return commentModel, nil
}

func (c *CommentHandler) DeleteCommentByID(id int64) error {
	err := c.CommentManager.DeleteCommentByID(id)
	if err != nil {
		return err
	}

	return nil
}

func mapModelToComment(commentModel *models.Comment) *comment.PostComment {
	now := time.Now().UTC()
	c := comment.PostComment{
		Title:     *commentModel.Title,
		Content:   *commentModel.Content,
		Creator:   *commentModel.Creator,
		CreatedAt: now.String(),
		UpdatedAt: now.String(),
	}
	return &c
}

func mapCommentToModel(post *comment.PostComment) *models.Comment {
	postModel := models.Comment{
		ID:      post.ID,
		Title:   &post.Title,
		Creator: &post.Creator,
		Content: &post.Content,
	}
	return &postModel
}
