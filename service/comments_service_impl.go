package service

import (
	"douyin/dao"
)

type CommentsServiceImpl struct {
}

func (c CommentsServiceImpl) SendComment(comment dao.Comment) (Comment, error) {
	commentDao, err := dao.InsertComment(comment)
	var commentInfo Comment
	commentInfo.Id = commentDao.Id
	commentInfo.UserID = commentDao.UserID
	commentInfo.VideoId = commentDao.VideoId
	commentInfo.Content = commentDao.Content
	return commentInfo, err

}

func (c CommentsServiceImpl) DelComment(comment dao.Comment) error {
	err := dao.DeleteComment(comment)

	return err

}
