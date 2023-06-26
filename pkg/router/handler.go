package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ksindhwani/imagegram/pkg/httputils"
	"github.com/ksindhwani/imagegram/pkg/service"
)

const (
	htmlImageTagName   = "image"
	htmlCaptionTagName = "caption"
	htmlUserIdTag      = "userId"
	defaultCursor      = "0"
	defaultPageSize    = "10"
)

type PostHandler struct {
	Service *service.PostService
}

type CommentHandler struct {
	Service *service.CommentService
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{
		Service: service,
	}
}

func NewCommentHandler(service *service.CommentService) *CommentHandler {
	return &CommentHandler{
		Service: service,
	}
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func (ph *PostHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	// Making sure image is not more than 100 Mb
	r.ParseMultipartForm(100 << 20)

	// Retrieve the file and caption from the form data
	file, handler, err := r.FormFile(htmlImageTagName)
	if err != nil {
		httputils.NewBadRequestError(err, "error in image reterival")
		return
	}
	defer file.Close()
	caption := r.FormValue(htmlCaptionTagName)
	userIdParam := r.FormValue(htmlUserIdTag)
	if userIdParam == "" {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(errors.New("error "), "no userId present in the request"))
		return
	}
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(errors.New("userId in url should be integer"), ""))
		return
	}

	post := service.Post{
		Caption: caption,
		UserId:  int64(userId),
	}

	response, err := ph.Service.CreateNewPost(post, handler.Filename, file)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewInternalServerError(err, "unable to create post"))
		return
	}

	httputils.WriteResponse(w, http.StatusCreated, response)

}

func (ch *CommentHandler) CommentOnPost(w http.ResponseWriter, r *http.Request) {
	var comment service.Comment
	postIdParam, err := httputils.GetUrlParam(r, "postId")
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(
			fmt.Errorf("bad request - %w", err), "unable to fetch postId from url"))
		return
	}
	postId, err := strconv.Atoi(postIdParam)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(errors.New("postId in url should be integer"), ""))
		return
	}
	body, err := httputils.GetRequestBody(w, r)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(err, "unable to parse request body"))
		return
	}

	if err := json.Unmarshal(body, &comment); err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(err, "unable to marshal request body"))
		return
	}
	comment.PostId = int64(postId)
	response, err := ch.Service.AddNewCommentOnPost(comment)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewInternalServerError(err, "unable to save comment"))
		return
	}

	httputils.WriteResponse(w, http.StatusCreated, response)
}

func (ch *CommentHandler) DeleteCommentOnPost(w http.ResponseWriter, r *http.Request) {
	commentIdPAram, err := httputils.GetUrlParam(r, "commentId")
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(
			fmt.Errorf("bad request - %w", err), "unable to fetch params from url"))
		return
	}
	commentId, err := strconv.Atoi(commentIdPAram)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(errors.New("commentId in url should be integer"), ""))
		return
	}
	response, err := ch.Service.DeleteComment(int64(commentId))
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewInternalServerError(err, "unable to delete comment"))
		return
	}
	httputils.WriteResponse(w, http.StatusOK, response)
}

func (ph *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	cursor, pageSize, err := getCursorAndPageSize(r)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewBadRequestError(err, "invalid cursor or pagesize"))
		return
	}
	response, err := ph.Service.GetAllPosts(cursor, pageSize)
	if err != nil {
		httputils.WriteErrorResponse(w, httputils.NewInternalServerError(err, "unable to get all posts"))
		return
	}
	httputils.WriteResponse(w, http.StatusOK, response)
}

func getCursorAndPageSize(r *http.Request) (int, int, error) {
	// Parse query parameters
	cursor := r.URL.Query().Get("cursor")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Set default values for cursor and pageSize
	if cursor == "" {
		cursor = defaultCursor
	}
	if pageSizeStr == "" {
		pageSizeStr = defaultPageSize
	}
	// Convert cursor and pageSize to integers
	cursorInt, err := strconv.Atoi(cursor)
	if err != nil {
		return -1, -1, errors.New("invalid cursor")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return -1, -1, errors.New("invalid pagesize")
	}
	return cursorInt, pageSize, nil
}
