package service

import "douyin/dao"

type Comment struct {
	Id         int64 `json:"id,omitempty" gorm:"primary_key;"`
	User       User  `json:"user" gorm:"foreignKey:Id;references:UserID;"`
	UserID     int64
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoId    int64
}

type CommentService interface {
	SendComment(comment dao.Comment) (Comment, error)
	DelComment(commentId int64) (Comment, error)
	GetList(videoId int64, userId int64) ([]Comment, error)
}
