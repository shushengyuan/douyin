package controller

import (
	"douyin/dao"
	"douyin/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "video_id is invalid"},
		})
		return
	}

	var video service.Video
	var num int32
	tx := dao.GetDB().Begin()
	if actionType == "1" {
		num = 1
		likes := service.Like{uint(user.Id), uint(video_id)}
		if err := tx.Create(likes).Error; err != nil {
			fmt.Println(err)
			tx.Rollback()
			return
		}
	} else if actionType == "2" {
		num = -1
		if err := tx.Where("user_id = ?", uint(user.Id)).Delete(&service.Like{}).Error; err != nil {
			fmt.Println(err)
			tx.Rollback()
			return
		}
	} else {
		tx.Rollback()
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "actionType is invalid"},
		})
	}

	if err := tx.Where("id = ?", video_id).
		Find(&video).
		UpdateColumn("videos.favorite_count", gorm.Expr("videos.favorite_count + ?", num)).
		Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return
	}
	if err := tx.Model(&service.User{}).
		Where("id = ?", video.AuthorID).
		UpdateColumn("total_favorited", gorm.Expr("total_favorited + ?", num)).
		Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	if err := tx.Model(&service.User{}).
		Where("id = ?", user.Id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", num)).
		Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	if err := tx.Commit().Error; err != nil {
		fmt.Println(err)
		// 如果提交失败，则回滚事务
		tx.Rollback()
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "点赞已上传成功！",
	})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	userID := c.Query("user_id")

	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	var videos []service.Video
	VideoForFavorite(userID, &videos)
	fmt.Println(videos)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
