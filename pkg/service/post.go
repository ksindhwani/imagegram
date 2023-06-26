package service

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/filesystem"
	"github.com/ksindhwani/imagegram/pkg/internal/tables"
)

type PostService struct {
	Config     config.Config
	Database   database.Database
	FileSystem filesystem.FileSystem
}

func NewPostService(
	Config *config.Config,
	database database.Database,
	fileSystem filesystem.FileSystem,
) *PostService {
	return &PostService{
		Config:     *Config,
		Database:   database,
		FileSystem: fileSystem,
	}
}

type Post struct {
	UserId    int64     `json:"userId"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostCommentResponse struct {
	PostId        int64     `json:"postId"`
	UserId        int64     `json:"userId"`
	Caption       string    `json:"caption"`
	ImageName     string    `json:"imageName"`
	ImageLocation string    `json:"imageLocation"`
	CreatedAt     time.Time `json:"createdAt"`
	Comments      []Comment `json:"comments"`
}
type PostResponse struct {
	PostId  int64 `json:"postId"`
	Success bool  `json:"success"`
}

func (ps *PostService) CreateNewPost(post Post, fileName string, file multipart.File) (PostResponse, error) {
	destinatinoUrl, err := ps.FileSystem.SaveFile(fileName, file)
	if err != nil {
		return PostResponse{}, fmt.Errorf("error in saving file - %w", err)
	}
	postId, err := ps.savePost(post, fileName, destinatinoUrl)
	if err != nil {
		return PostResponse{}, fmt.Errorf("error in saving post in database - %w", err)
	}
	return PostResponse{
		PostId:  postId,
		Success: true,
	}, nil
}

func (ps *PostService) GetAllPosts(cursor int, pageSize int) (map[int64]PostCommentResponse, error) {
	posts, err := ps.Database.GetAllPostWithLast2Comments(cursor, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error - %w", err)
	}

	response, err := parseDataIntoResponseFormat(posts)
	if err != nil {
		return nil, fmt.Errorf("error - %w", err)
	}
	return response, nil
}

func parseDataIntoResponseFormat(posts []database.AllPostsJoinQueryResult) (map[int64]PostCommentResponse, error) {
	postsMap := make(map[int64]PostCommentResponse)
	for _, post := range posts {
		key := post.PostId
		if _, ok := postsMap[key]; ok {
			postCommentValue := postsMap[key]
			comment := Comment{
				PostId:    key,
				UserId:    post.CommentUserId,
				Content:   post.Comment,
				CreatedAt: post.CommentCreatedAt,
			}
			postCommentValue.Comments = append(postCommentValue.Comments, comment)
			postsMap[key] = postCommentValue
		} else {
			postsMap[key] = PostCommentResponse{
				PostId:        key,
				UserId:        post.UserId,
				Caption:       post.Caption,
				CreatedAt:     post.CreatedAt,
				ImageName:     post.PostImageName,
				ImageLocation: post.PostImageLocation,
				Comments: []Comment{
					{
						PostId:    key,
						UserId:    post.CommentUserId,
						Content:   post.Comment,
						CreatedAt: post.CommentCreatedAt,
					},
				},
			}
		}
	}
	return postsMap, nil
}

func (ps *PostService) savePost(post Post, fileName string, destinatinoUrl string) (int64, error) {
	postTableRow := tables.PostTable{
		Caption: post.Caption,
		UserId:  post.UserId,
	}
	imageTableRow := tables.ImageTable{
		ImageFileName: fileName,
		Location:      destinatinoUrl,
	}
	return ps.Database.InsertNewPost(postTableRow, imageTableRow)
}
