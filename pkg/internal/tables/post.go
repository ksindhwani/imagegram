package tables

import "time"

type PostTable struct {
	PostId    int64
	UserId    int64
	Caption   string
	CreatedAt time.Time
}
