package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddNewCommentOnPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name                        string
		Input                       Comment
		ExpectedSaveCommentResponse int64
		ExpectedSaveCommentError    error
		ExpectedSaveCommentCalls    int
		ExpectedResponse            CommentResponse
		ExpectedError               error
	}{
		{
			Name: "Test All Valid",
			Input: Comment{
				PostId:    1,
				UserId:    1,
				Content:   "Test Comment",
				CreatedAt: time.Now(),
			},
			ExpectedSaveCommentResponse: 1,
			ExpectedSaveCommentError:    nil,
			ExpectedSaveCommentCalls:    1,
			ExpectedResponse: CommentResponse{
				CommentId: 1,
				Success:   true,
			},
			ExpectedError: nil,
		},
		{
			Name: "Test error in db query execution",
			Input: Comment{
				PostId:    1,
				UserId:    1,
				Content:   "Test Comment",
				CreatedAt: time.Now(),
			},
			ExpectedSaveCommentResponse: 0,
			ExpectedSaveCommentError:    errors.New("error in query execution"),
			ExpectedSaveCommentCalls:    1,
			ExpectedResponse:            CommentResponse{},
			ExpectedError:               fmt.Errorf("error in saving commment - %w", errors.New("error in query execution")),
		},
		{
			Name: "Test error in getting last comment id",
			Input: Comment{
				PostId:    1,
				UserId:    1,
				Content:   "Test Comment",
				CreatedAt: time.Now(),
			},
			ExpectedSaveCommentResponse: 0,
			ExpectedSaveCommentError:    errors.New("error in fetching last comment id"),
			ExpectedSaveCommentCalls:    1,
			ExpectedResponse:            CommentResponse{},
			ExpectedError:               fmt.Errorf("error in saving commment - %w", errors.New("error in fetching last comment id")),
		},
		{
			Name: "Test table not found",
			Input: Comment{
				PostId:    1,
				UserId:    1,
				Content:   "Test Comment",
				CreatedAt: time.Now(),
			},
			ExpectedSaveCommentResponse: 0,
			ExpectedSaveCommentError:    errors.New("table not found"),
			ExpectedSaveCommentCalls:    1,
			ExpectedResponse:            CommentResponse{},
			ExpectedError:               fmt.Errorf("error in saving commment - %w", errors.New("table not found")),
		},
	}

	any := gomock.Any()
	config := config.Config{
		HostImageDirectory:  "test host directory",
		LocalImageDirectory: "test local directory",
	}
	database := mocks.NewMockDatabase(ctrl)
	commentService := NewCommentService(&config, database, nil)
	for _, test := range tests {
		database.EXPECT().SaveComment(any).
			Return(test.ExpectedSaveCommentResponse, test.ExpectedSaveCommentError).
			Times(test.ExpectedSaveCommentCalls)
		result, err := commentService.AddNewCommentOnPost(test.Input)
		assert.Equal(t, test.ExpectedResponse, result, test.Name)
		assert.Equal(t, test.ExpectedError, err, test.Name)
	}
}

func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name                       string
		Input                      int64
		ExpectedDeleteCommentError error
		ExpectedDeleteCommentCalls int
		ExpectedResponse           CommentResponse
		ExpectedError              error
	}{
		{
			Name:                       "Test All Valid",
			Input:                      1,
			ExpectedDeleteCommentError: nil,
			ExpectedDeleteCommentCalls: 1,
			ExpectedResponse: CommentResponse{
				CommentId: 1,
				Success:   true,
			},
			ExpectedError: nil,
		},
		{
			Name:                       "Test error in db execution",
			Input:                      1,
			ExpectedDeleteCommentError: errors.New("error in db execution"),
			ExpectedDeleteCommentCalls: 1,
			ExpectedResponse:           CommentResponse{},
			ExpectedError:              fmt.Errorf("error in deleting commment - %w", errors.New("error in db execution")),
		},
		{
			Name:                       "Test when no comment id not found and no row is deleted",
			Input:                      2,
			ExpectedDeleteCommentError: errors.New("no row found with the given comment id"),
			ExpectedDeleteCommentCalls: 1,
			ExpectedResponse:           CommentResponse{},
			ExpectedError:              fmt.Errorf("error in deleting commment - %w", errors.New("no row found with the given comment id")),
		},
		{
			Name:                       "Test error in fetching rows affected",
			Input:                      1,
			ExpectedDeleteCommentError: errors.New("unable to fetch row effected"),
			ExpectedDeleteCommentCalls: 1,
			ExpectedResponse:           CommentResponse{},
			ExpectedError:              fmt.Errorf("error in deleting commment - %w", errors.New("unable to fetch row effected")),
		},
	}

	any := gomock.Any()
	config := config.Config{
		HostImageDirectory:  "test host directory",
		LocalImageDirectory: "test local directory",
	}
	database := mocks.NewMockDatabase(ctrl)
	commentService := NewCommentService(&config, database, nil)
	for _, test := range tests {
		database.EXPECT().DeleteComment(any).
			Return(test.ExpectedDeleteCommentError).
			Times(test.ExpectedDeleteCommentCalls)
		result, err := commentService.DeleteComment(test.Input)
		assert.Equal(t, test.ExpectedResponse, result, test.Name)
		assert.Equal(t, test.ExpectedError, err, test.Name)
	}
}
