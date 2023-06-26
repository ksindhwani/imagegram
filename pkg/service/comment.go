package service

import (
	"fmt"
	"time"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/filesystem"
	"github.com/ksindhwani/imagegram/pkg/internal/tables"
)

type CommentService struct {
	Config     config.Config
	Database   database.Database
	FileSystem filesystem.FileSystem
}

func NewCommentService(
	Config *config.Config,
	database database.Database,
	fileSystem filesystem.FileSystem,
) *CommentService {
	return &CommentService{
		Config:     *Config,
		Database:   database,
		FileSystem: fileSystem,
	}
}

type Comment struct {
	PostId    int64     `json:"postId"`
	UserId    int64     `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentResponse struct {
	CommentId int64 `json:"commentId"`
	Success   bool  `json:"success"`
}

func (cs *CommentService) AddNewCommentOnPost(comment Comment) (CommentResponse, error) {
	commentTableRow := tables.CommentTable{
		PostId:  comment.PostId,
		UserId:  comment.UserId,
		Comment: comment.Content,
	}
	commentId, err := cs.Database.SaveComment(commentTableRow)
	if err != nil {
		return CommentResponse{}, fmt.Errorf("error in saving commment - %w", err)
	}
	return CommentResponse{
		CommentId: commentId,
		Success:   true,
	}, nil
}

func (cs *CommentService) DeleteComment(commentId int64) (CommentResponse, error) {
	err := cs.Database.DeleteComment(commentId)
	if err != nil {
		return CommentResponse{}, fmt.Errorf("error in deleting commment - %w", err)
	}
	return CommentResponse{
		CommentId: commentId,
		Success:   true,
	}, nil
}
