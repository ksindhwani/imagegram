package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ksindhwani/imagegram/pkg/internal/converter"
	"github.com/ksindhwani/imagegram/pkg/internal/tables"
)

type Database interface {
	InsertNewPost(postTableRow tables.PostTable, imageTableRow tables.ImageTable) (int64, error)
	SaveComment(comment tables.CommentTable) (int64, error)
	DeleteComment(commentId int64) error
	GetAllPostWithLast2Comments(cursor int, pageSize int) ([]AllPostsJoinQueryResult, error)
	GetAllImages() ([]tables.ImageTable, error)
	UpdateImageConvertedData(image converter.ImageConversionResponse) error
}

type database struct {
	Db *sql.DB
}

func New(db *sql.DB) Database {
	return &database{
		Db: db,
	}
}

type AllPostsJoinQueryResult struct {
	PostId            int64
	UserId            int64
	Caption           string
	CreatedAt         time.Time
	PostImageName     string
	PostImageLocation string
	CommentId         int64
	CommentUserId     int64
	Comment           string
	CommentCreatedAt  time.Time
}

// Insert New Post and image in database
func (d *database) InsertNewPost(postTableRow tables.PostTable, imageTableRow tables.ImageTable) (int64, error) {
	tx, err := d.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Rollback the transaction if there is an error

	insertPostQuery := "INSERT INTO `posts` (`user_id`, caption) VALUES (?, ?)"
	stmt, err := tx.Prepare(insertPostQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(postTableRow.UserId, postTableRow.Caption)
	if err != nil {
		return 0, err
	}

	// Get the inserted post ID
	postId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert Image
	imageTableRow.PostId = postId

	insertImageQuery := "INSERT INTO `images` (`post_id`, image_file_name, location) VALUES (?, ?, ?)"
	stmt, err = tx.Prepare(insertImageQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(imageTableRow.PostId, imageTableRow.ImageFileName, imageTableRow.Location)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return postId, nil
}

// Save new comment in database
func (d *database) SaveComment(comment tables.CommentTable) (int64, error) {
	insertQuery := "INSERT INTO comments (post_id, user_id, comment) VALUES (?, ?, ?)"
	result, err := d.Db.Exec(insertQuery, comment.PostId, comment.UserId, comment.Comment)
	if err != nil {
		return 0, err
	}
	// Get the inserted user's ID
	commentId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return commentId, nil
}

// Delete Comment from database
func (d *database) DeleteComment(commentId int64) error {
	deleteQuery := "DELETE FROM comments where comment_id = ?"
	result, err := d.Db.Exec(deleteQuery, commentId)
	if err != nil {
		return err
	}

	// Check the number of rows affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no row found with the given comment id")
	}
	return err

}

func (d *database) GetAllPostWithLast2Comments(cursor int, pageSize int) ([]AllPostsJoinQueryResult, error) {
	// Sql Query to get all posts with last 2 comments
	query := "SELECT " +
		"p.post_id, p.user_id, p.caption, p.created_at, " +
		"IFNULL(i.converted_image_name, ''), IFNULL(i.converted_image_location, ''), " +
		"c.comment_id, c.user_id,c.comment,c.created_at " +
		"FROM posts p " +
		"LEFT JOIN ( " +
		"SELECT comment_id, post_id,user_id, comment, created_at, " +
		"ROW_NUMBER() OVER (PARTITION BY post_id ORDER BY comment_id DESC) AS rn FROM comments" +
		") c ON p.post_id = c.post_id " +
		"INNER JOIN images i on p.post_id = i.post_id " +
		"WHERE (c.rn <= 2 OR c.comment_id IS NULL) AND p.post_id > ? " +
		"ORDER BY p.post_id, c.comment_id DESC " +
		"LIMIT ?"

	// Execute the query
	rows, err := d.Db.Query(query, cursor, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to store the results
	var results []AllPostsJoinQueryResult

	// Iterate over the rows and map the results to the struct
	for rows.Next() {
		var result AllPostsJoinQueryResult
		err := rows.Scan(
			&result.PostId,
			&result.UserId,
			&result.Caption,
			&result.CreatedAt,
			&result.PostImageName,
			&result.PostImageLocation,
			&result.CommentId,
			&result.CommentUserId,
			&result.Comment,
			&result.CommentCreatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *database) GetAllImages() ([]tables.ImageTable, error) {
	var images []tables.ImageTable
	selectQuery := "SELECT " +
		"`image_id`, " +
		"`post_id`, " +
		"`image_file_name`, " +
		"`location`, " +
		"IFNULL(`converted_image_name`, ''), " +
		"IFNULL(`converted_image_location`, ''), " +
		"`uploaded_at` " +
		"FROM `images` " +
		"WHERE `converted_image_name` is NULL"
	// Execute the query
	rows, err := d.Db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and map the results to the struct
	for rows.Next() {
		var image tables.ImageTable
		err := rows.Scan(
			&image.ImageId,
			&image.PostId,
			&image.ImageFileName,
			&image.Location,
			&image.ConvertedImageName,
			&image.ConvertedImageLocation,
			&image.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return images, nil
}

func (d *database) UpdateImageConvertedData(image converter.ImageConversionResponse) error {
	updateQuery := "UPDATE `images` SET `converted_image_name` = ?, converted_image_location = ? WHERE `image_id` = ?"
	result, err := d.Db.Exec(updateQuery, image.ConvertedImageName, image.ConvertedImageLocation, image.ImageId)
	if err != nil {
		return err
	}

	// Check the number of rows affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no row updated with the given image id")
	}
	return err

}
