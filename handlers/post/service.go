package posts

import (
	"github.com/asaberwd/atomica-blog/internal/post"
	"github.com/asaberwd/atomica-blog/swagger/models"
	"time"
)

type PostHandler struct {
	PostManager Service
}

func New(manager Service) *PostHandler {
	return &PostHandler{
		PostManager: manager,
	}
}

type Service interface {
	CreatePost(blogPost *post.BlogPost) (int64, error)
	ListPosts() ([]post.BlogPost, error)
	GetPostByID(id int64) (*post.BlogPost, error)
	DeletePostByID(id int64) error
}

func (p *PostHandler) CreatePost(post *models.Post) (*models.Post, error) {
	bPost := mapModelToPost(post)
	id, err := p.PostManager.CreatePost(bPost)
	if err != nil {
		return post, err
	}

	blogPost, err := p.GetPostByID(id)
	if err != nil {
		return blogPost, err
	}
	return blogPost, nil
}

func (p *PostHandler) ListPosts() (models.Posts, error) {
	postsModel := models.Posts{Posts: []*models.Post{}}
	blogPosts, err := p.PostManager.ListPosts()
	if err != nil {
		return models.Posts{}, err
	}
	for _, b := range blogPosts {
		postsModel.Posts = append(postsModel.Posts, mapPostToModel(&b))
	}

	return postsModel, nil
}

func (p *PostHandler) GetPostByID(id int64) (*models.Post, error) {
	blogPost, err := p.PostManager.GetPostByID(id)
	if err != nil {
		return &models.Post{}, err
	}
	if blogPost == nil {
		return nil, nil
	}
	postModel := mapPostToModel(blogPost)
	return postModel, nil
}

func (p *PostHandler) DeletePostByID(id int64) error {
	err := p.PostManager.DeletePostByID(id)
	if err != nil {
		return err
	}

	return nil
}

func mapModelToPost(postModel *models.Post) *post.BlogPost {
	now := time.Now().UTC()
	p := post.BlogPost{
		Title:     *postModel.Title,
		Content:   *postModel.Content,
		Photos:    postModel.PhotoUrls,
		Tags:      postModel.Tags,
		CreatedAt: now.String(),
		UpdatedAt: now.String(),
	}
	return &p
}

func mapPostToModel(post *post.BlogPost) *models.Post {
	postModel := models.Post{
		Content:   &post.Content,
		ID:        post.ID,
		PhotoUrls: post.Photos,
		Tags:      post.Tags,
		Title:     &post.Title,
	}
	for range post.Tags {

	}
	return &postModel
}
