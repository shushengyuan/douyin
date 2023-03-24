package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"douyin/dao"
	"douyin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []service.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment dao.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
	}

	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	commentService := new(service.CommentsServiceImpl)

	if actionType == "1" {
		text := c.Query("comment_text")
		currentTime := fmt.Sprintf("%d-%d", time.Now().Month(), time.Now().Day())

		comment := dao.Comment{User: user, Content: text, CreateDate: currentTime, VideoId: videoId}
		commentInfo, createCommentErr := commentService.SendComment(comment)

		if createCommentErr != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  createCommentErr.Error(),
			})
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: commentInfo})
		return
	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		deleteCommentErr := dao.GetDB().Where("id = ?", commentId).Delete(&service.Comment{}).Error
		dao.GetDB().Model(&service.Video{}).
			Where("id = ?", videoId).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", -1))
		if deleteCommentErr != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  deleteCommentErr.Error(),
			})
		}
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	token := c.Query("token")
	videoId := c.Query("video_id")

	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: verifyErr.Error()},
		})
		return
	}

	var comments []service.Comment
	err := dao.GetDB().Where("comments.video_id = ?", videoId).Order("comments.id desc").
		Find(&comments).Error

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "cannot get comments"},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
		},
		CommentList: comments,
	})

}
