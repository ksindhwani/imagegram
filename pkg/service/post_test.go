package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type NewPostInput struct {
		post     Post
		fileName string
		file     multipart.File
	}

	tests := []struct {
		Name                          string
		Input                         NewPostInput
		ExpectedSaveFileResponse      string
		ExpectedSaveFileError         error
		ExpectedInsertNewPostResponse int64
		ExpectedInsertNewPostError    error
		ExpectedSaveFileCalls         int
		ExpectedInsertNewPostCalls    int
		ExpectedResponse              PostResponse
		ExpectedError                 error
	}{
		{
			Name: "Test All Valid ",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "/images/test.png",
			ExpectedSaveFileError:         nil,
			ExpectedInsertNewPostResponse: 1,
			ExpectedInsertNewPostError:    nil,
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    1,
			ExpectedResponse: PostResponse{
				PostId:  1,
				Success: true,
			},
			ExpectedError: nil,
		},
		{
			Name: "Test when lcoal directory not exists",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "",
			ExpectedSaveFileError:         errors.New("directory not exist"),
			ExpectedInsertNewPostResponse: 1,
			ExpectedInsertNewPostError:    nil,
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    0,
			ExpectedResponse:              PostResponse{},
			ExpectedError:                 fmt.Errorf("error in saving file - %w", errors.New("directory not exist")),
		},
		{
			Name: "Test when error in database occured",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "/images/test.png",
			ExpectedSaveFileError:         nil,
			ExpectedInsertNewPostResponse: 0,
			ExpectedInsertNewPostError:    errors.New("error in database"),
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    1,
			ExpectedResponse:              PostResponse{},
			ExpectedError:                 fmt.Errorf("error in saving post in database - %w", errors.New("error in database")),
		},
		{
			Name: "Test when error in rollback error occured",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "/images/test.png",
			ExpectedSaveFileError:         nil,
			ExpectedInsertNewPostResponse: 0,
			ExpectedInsertNewPostError:    errors.New("rollback"),
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    1,
			ExpectedResponse:              PostResponse{},
			ExpectedError:                 fmt.Errorf("error in saving post in database - %w", errors.New("rollback")),
		},
		{
			Name: "Test when error in commit error occured",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "/images/test.png",
			ExpectedSaveFileError:         nil,
			ExpectedInsertNewPostResponse: 0,
			ExpectedInsertNewPostError:    errors.New("error in commit"),
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    1,
			ExpectedResponse:              PostResponse{},
			ExpectedError:                 fmt.Errorf("error in saving post in database - %w", errors.New("error in commit")),
		},
		{
			Name: "Test when error in prepare statement occured",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "/images/test.png",
			ExpectedSaveFileError:         nil,
			ExpectedInsertNewPostResponse: 0,
			ExpectedInsertNewPostError:    errors.New("error in sql prepare statement"),
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    1,
			ExpectedResponse:              PostResponse{},
			ExpectedError:                 fmt.Errorf("error in saving post in database - %w", errors.New("error in sql prepare statement")),
		},
		{
			Name: "Test when error in parsing query response into model occured",
			Input: NewPostInput{
				post: Post{
					UserId:    1,
					Caption:   "Test Post Caption",
					CreatedAt: time.Now(),
				},
				fileName: "test.png",
				file:     nil,
			},
			ExpectedSaveFileResponse:      "/images/test.png",
			ExpectedSaveFileError:         nil,
			ExpectedInsertNewPostResponse: 0,
			ExpectedInsertNewPostError:    errors.New("can't find the column"),
			ExpectedSaveFileCalls:         1,
			ExpectedInsertNewPostCalls:    1,
			ExpectedResponse:              PostResponse{},
			ExpectedError:                 fmt.Errorf("error in saving post in database - %w", errors.New("can't find the column")),
		},
	}

	any := gomock.Any()
	config := config.Config{
		HostImageDirectory:  "test host directory",
		LocalImageDirectory: "test local directory",
	}
	localFileSystem := mocks.NewMockFileSystem(ctrl)
	database := mocks.NewMockDatabase(ctrl)
	postService := NewPostService(&config, database, localFileSystem)
	for _, test := range tests {
		localFileSystem.EXPECT().SaveFile(any, any).Return(test.ExpectedSaveFileResponse, test.ExpectedSaveFileError).Times(test.ExpectedSaveFileCalls)
		database.EXPECT().InsertNewPost(any, any).Return(test.ExpectedInsertNewPostResponse, test.ExpectedInsertNewPostError).Times(test.ExpectedInsertNewPostCalls)
		result, err := postService.CreateNewPost(test.Input.post, test.Input.fileName, test.Input.file)
		assert.Equal(t, test.ExpectedResponse, result, test.Name)
		assert.Equal(t, test.ExpectedError, err, test.Name)
	}
}

func TestGetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type GetAllPostsInput struct {
		cursor   int
		pageSize int
	}

	tests := []struct {
		Name                                         string
		Input                                        GetAllPostsInput
		ExpectedIGetAllPostWithLast2CommentsResponse []database.AllPostsJoinQueryResult
		ExpectedIGetAllPostWithLast2CommentsError    error
		ExpectedGetAllPostWithLast2CommentsCalls     int
		ExpectedResponse                             map[int64]PostCommentResponse
		ExpectedError                                error
	}{
		{
			Name: "Test All Valid ",
			Input: GetAllPostsInput{
				cursor:   0,
				pageSize: 10,
			},
			ExpectedIGetAllPostWithLast2CommentsResponse: []database.AllPostsJoinQueryResult{
				{
					PostId:            1,
					UserId:            1,
					Caption:           "test Caption post user 1",
					CreatedAt:         time.Date(2023, 6, 26, 0, 0, 0, 1, &time.Location{}),
					PostImageName:     "test.png",
					PostImageLocation: "/images/test.png",
					CommentId:         1,
					CommentUserId:     2,
					Comment:           "comment by user 2",
					CommentCreatedAt:  time.Date(2023, 6, 26, 0, 0, 0, 3, &time.Location{}),
				},
				{
					PostId:            1,
					UserId:            1,
					Caption:           "test Caption post user 1",
					CreatedAt:         time.Date(2023, 6, 26, 0, 0, 0, 1, &time.Location{}),
					PostImageName:     "test.png",
					PostImageLocation: "/images/test.png",
					CommentId:         1,
					CommentUserId:     3,
					Comment:           "comment by user 3",
					CommentCreatedAt:  time.Date(2023, 6, 26, 0, 0, 0, 2, &time.Location{}),
				},
			},
			ExpectedIGetAllPostWithLast2CommentsError: nil,
			ExpectedGetAllPostWithLast2CommentsCalls:  1,
			ExpectedResponse: map[int64]PostCommentResponse{
				1: {
					PostId:        1,
					UserId:        1,
					Caption:       "test Caption post user 1",
					CreatedAt:     time.Date(2023, 6, 26, 0, 0, 0, 1, &time.Location{}),
					ImageName:     "test.png",
					ImageLocation: "/images/test.png",
					Comments: []Comment{
						{
							PostId:    1,
							UserId:    2,
							Content:   "comment by user 2",
							CreatedAt: time.Date(2023, 6, 26, 0, 0, 0, 3, &time.Location{}),
						},
						{
							PostId:    1,
							UserId:    3,
							Content:   "comment by user 3",
							CreatedAt: time.Date(2023, 6, 26, 0, 0, 0, 2, &time.Location{}),
						},
					},
				},
			},
			ExpectedError: nil,
		},
		{
			Name: "Test all pages are traversed ",
			Input: GetAllPostsInput{
				cursor:   11,
				pageSize: 10,
			},
			ExpectedIGetAllPostWithLast2CommentsResponse: []database.AllPostsJoinQueryResult{},
			ExpectedIGetAllPostWithLast2CommentsError:    nil,
			ExpectedGetAllPostWithLast2CommentsCalls:     1,
			ExpectedResponse:                             map[int64]PostCommentResponse{},
			ExpectedError:                                nil,
		},
		{
			Name: "Test when error in row scaning after join query",
			Input: GetAllPostsInput{
				cursor:   0,
				pageSize: 10,
			},
			ExpectedIGetAllPostWithLast2CommentsResponse: nil,
			ExpectedIGetAllPostWithLast2CommentsError:    errors.New("error in row scanning"),
			ExpectedGetAllPostWithLast2CommentsCalls:     1,
			ExpectedResponse:                             nil,
			ExpectedError:                                fmt.Errorf("error - %w", errors.New("error in row scanning")),
		},
		{
			Name: "Test when error in db query execution",
			Input: GetAllPostsInput{
				cursor:   0,
				pageSize: 10,
			},
			ExpectedIGetAllPostWithLast2CommentsResponse: nil,
			ExpectedIGetAllPostWithLast2CommentsError:    errors.New("error in db query execution"),
			ExpectedGetAllPostWithLast2CommentsCalls:     1,
			ExpectedResponse:                             nil,
			ExpectedError:                                fmt.Errorf("error - %w", errors.New("error in db query execution")),
		},
		{
			Name: "Test when error in rows after join",
			Input: GetAllPostsInput{
				cursor:   0,
				pageSize: 10,
			},
			ExpectedIGetAllPostWithLast2CommentsResponse: nil,
			ExpectedIGetAllPostWithLast2CommentsError:    errors.New("error in rows after join"),
			ExpectedGetAllPostWithLast2CommentsCalls:     1,
			ExpectedResponse:                             nil,
			ExpectedError:                                fmt.Errorf("error - %w", errors.New("error in rows after join")),
		},
	}

	any := gomock.Any()
	config := config.Config{
		HostImageDirectory:  "test host directory",
		LocalImageDirectory: "test local directory",
	}
	database := mocks.NewMockDatabase(ctrl)
	postService := NewPostService(&config, database, nil)
	for _, test := range tests {
		database.EXPECT().GetAllPostWithLast2Comments(any, any).
			Return(test.ExpectedIGetAllPostWithLast2CommentsResponse, test.ExpectedIGetAllPostWithLast2CommentsError).
			Times(test.ExpectedGetAllPostWithLast2CommentsCalls)
		result, err := postService.GetAllPosts(test.Input.cursor, test.Input.pageSize)
		assert.Equal(t, test.ExpectedResponse, result, test.Name)
		assert.Equal(t, test.ExpectedError, err, test.Name)
	}
}
