package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ksindhwani/imagegram/pkg/app"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/service"
)

func New(deps *app.Dependencies) (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/ping", PingHandler).Methods(http.MethodGet)

	database := database.New(deps.DB)
	postService := service.NewPostService(deps.Config, database, deps.LocalFileSystem)
	commmentService := service.NewCommentService(deps.Config, database, deps.LocalFileSystem)
	postHandler := NewPostHandler(postService)
	commentHandler := NewCommentHandler(commmentService)

	r.HandleFunc("/posts", postHandler.CreateNewPost).Methods(http.MethodPost)
	r.HandleFunc("/posts/{postId}/comments", commentHandler.CommentOnPost).Methods(http.MethodPost)
	r.HandleFunc("/posts/{postId}/comments/{commentId}", commentHandler.DeleteCommentOnPost).Methods(http.MethodDelete)
	r.HandleFunc("/posts", postHandler.GetAllPosts).Methods(http.MethodGet)
	return r, nil
}
