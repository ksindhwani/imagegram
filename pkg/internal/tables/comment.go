package tables

import "time"

type CommentTable struct {
	CommentId int64
	PostId    int64
	UserId    int64
	Comment   string
	CreatedAt time.Time
}
