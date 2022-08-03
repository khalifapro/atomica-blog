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

func (p *CommentHandler) CreateComment(comment *models.Comment) (*models.Comment, error) {
	pComment := mapModelToComment(comment)
	_, err := p.CommentManager.CreateComment(pComment)
	if err != nil {
		return comment, err
	}

	/*	postComment, err := p.GetPostByID(id)
		if err != nil {
			return postComment, err
		}*/
	return comment, nil
}

func (p *CommentHandler) ListPostComments(postID int64) (models.Comments, error) {
	commentsModel := models.Comments{}
	blogPosts, err := p.CommentManager.ListPostComments(postID)
	if err != nil {
		return models.Comments{}, err
	}
	for _, b := range blogPosts {
		commentsModel = append(commentsModel, mapCommentToModel(&b))
	}

	return commentsModel, nil
}

func (p *CommentHandler) UpdatePostComment(id int64) (*models.Comment, error) {
	postComment, err := p.CommentManager.GetCommentByID(id)
	if err != nil {
		return &models.Comment{}, err
	}
	if postComment == nil {
		return nil, nil
	}
	postModel := mapCommentToModel(postComment)
	return postModel, nil
}

func (p *CommentHandler) DeleteCommentByID(id int64) error {
	err := p.CommentManager.DeleteCommentByID(id)
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
