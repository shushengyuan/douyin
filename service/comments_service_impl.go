package service

import (
	"douyin/dao"
)

type CommentsServiceImpl struct {
}

func (c CommentsServiceImpl) SendComment(comment dao.Comment) (dao.Comment, error) {
	commentDao, err := dao.InsertComment(comment)

	return commentDao, err

}

func (c CommentsServiceImpl) DelComment(comment dao.Comment) error {
	err := dao.DeleteComment(comment)

	return err

}
