package dao

import (
	"gorm.io/gorm"
)

type Comment struct {
	Id         int64 `json:"id,omitempty" gorm:"primary_key;"`
	UserID     int64
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoId    int64
}

// InsertComment
// 2、发表评论
func InsertComment(comment Comment) (Comment, error) {
	createCommentErr := db.Create(&comment).Error
	db.Model(&Video{}).
		Where("id = ?", comment.VideoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	return comment, createCommentErr
}

func DeleteComment(comment Comment) error {
	deleteCommentErr := db.Where("id = ?", comment.Id).Delete(&Comment{}).Error
	db.Model(&Video{}).
		Where("id = ?", comment.VideoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", -1))

	return deleteCommentErr
}
