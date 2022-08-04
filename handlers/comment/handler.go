package comments

import (
	"fmt"
	"github.com/asaberwd/atomica-blog/swagger/models"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations/comment"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations/comments"
	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.AtomicaBlogServiceAPI, service CommentHandler) {
	api.CommentsGetPostCommentsByIDHandler = comments.GetPostCommentsByIDHandlerFunc(func(params comments.GetPostCommentsByIDParams) middleware.Responder {
		postComments, err := service.ListPostComments(params.PostID)
		if err != nil {
			fmt.Println(err)
			return comments.NewGetPostCommentsByIDBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		return comments.NewGetPostCommentsByIDOK().WithPayload(postComments)
	})

	api.CommentsAddPostCommentHandler = comments.AddPostCommentHandlerFunc(func(params comments.AddPostCommentParams) middleware.Responder {
		blogPost, err := service.CreateComment(params.Comment)
		if err != nil {
			fmt.Println(err)
			return comments.NewAddPostCommentBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		return comments.NewAddPostCommentOK().WithPayload(blogPost)
	})
	api.CommentUpdatePostCommentHandler = comment.UpdatePostCommentHandlerFunc(func(params comment.UpdatePostCommentParams) middleware.Responder {
		blogPost, err := service.CreateComment(params.Comment)
		if err != nil {
			fmt.Println(err)
			return comment.NewUpdatePostCommentBadRequest().WithPayload(&models.Error{Message: err.Error()})
		}
		return comment.NewUpdatePostCommentOK().WithPayload(blogPost)
	})

}
