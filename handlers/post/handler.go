package posts

import (
	"fmt"
	"github.com/asaberwd/atomica-blog/swagger/models"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations/post"
	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.AtomicaBlogServiceAPI, service PostHandler) {
	api.PostGetPostsHandler = post.GetPostsHandlerFunc(func(params post.GetPostsParams) middleware.Responder {
		blogPosts, err := service.ListPosts()
		if err != nil {
			fmt.Println(err)
			return post.NewGetPostsBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		return post.NewGetPostsOK().WithPayload(&blogPosts)
	})

	api.PostAddPostHandler = post.AddPostHandlerFunc(func(params post.AddPostParams) middleware.Responder {
		blogPost, err := service.CreatePost(params.Post)
		if err != nil {
			fmt.Println(err)
			return post.NewAddPostBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		return post.NewAddPostOK().WithPayload(blogPost)
	})

	api.PostGetPostByIDHandler = post.GetPostByIDHandlerFunc(func(params post.GetPostByIDParams) middleware.Responder {
		blogPost, err := service.GetPostByID(params.PostID)
		if err != nil {
			fmt.Println(err)
			return post.NewGetPostByIDBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		if blogPost == nil {
			return post.NewGetPostByIDNotFound()
		}
		return post.NewGetPostByIDOK().WithPayload(blogPost)
	})

	api.PostDeletePostHandler = post.DeletePostHandlerFunc(func(params post.DeletePostParams) middleware.Responder {
		err := service.DeletePostByID(params.PostID)
		if err != nil {
			fmt.Println(err)
			if err.Error() == "error deleting post, post does not exist" {
				return post.NewDeletePostNotFound()
			}
			return post.NewDeletePostBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		return post.NewDeletePostOK()
	})

	api.PostUpdatePostHandler = post.UpdatePostHandlerFunc(func(params post.UpdatePostParams) middleware.Responder {
		return middleware.NotImplemented("operation post.UpdatePost has not yet been implemented")
	})

}
