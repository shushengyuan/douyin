package controller

import (
	// "fmt"
	"douyin/dao"
	"douyin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []service.User `json:"user_list"`
}

type FriendUser struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty" default:"true"`
	FriendMessage string `json:"message,omitempty"`
	MsgType       int64  `json:"msgType,omitempty"`
}

type FriendUserListResponse struct {
	Response
	FriendUserList []FriendUser `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	if verifyErr == nil {
		// 对方用户id
		to_user_id := c.Query("to_user_id")

		var to_user service.User
		to_userExistErr := dao.GetDB().Where("id = ?", to_user_id).Take(&to_user).Error
		if to_userExistErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
			return
		}
		action_type := c.Query("action_type")

		// 关注操作
		if action_type == "1" {
			// 注意不能关注自己
			if user.Id == to_user.Id {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "您不能关注自己"})
				return
			}
			// 注意不能重复关注
			var relation dao.Relation
			relationTakeErr := dao.GetDB().Where("follow = ? AND follower = ?", to_user.Id, user.Id).Take(&relation).Error
			if relationTakeErr == nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "您不能重复关注该用户"})
				return
			}

			CreateFollowErr := dao.GetDB().Create(&dao.Relation{Follow: to_user.Id, Follower: user.Id}).Error
			if CreateFollowErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库插入失败"})
				return
			}

			user.FollowCount += 1
			userSaveErr := dao.GetDB().Save(&user).Error
			if userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}

			to_user.FollowerCount += 1
			to_userSaveErr := dao.GetDB().Save(&to_user).Error
			if to_userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}

		} else if action_type == "2" {
			// 取关操作
			DeleteFollowErr := dao.GetDB().Where("follow = ? AND follower = ?", to_user.Id, user.Id).Delete(&dao.Relation{}).Error
			if DeleteFollowErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库删除失败"})
				return
			}

			if user.FollowCount > 0 {
				user.FollowCount -= 1
			}
			userSaveErr := dao.GetDB().Save(&user).Error
			if userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}

			if to_user.FollowerCount > 0 {
				to_user.FollowerCount -= 1
			}
			to_userSaveErr := dao.GetDB().Save(&to_user).Error
			if to_userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 查找用户
	var user service.User
	userExitErr := dao.GetDB().Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// 在Relation数据查找关注者的id
	relations := []dao.Relation{}
	relationsFinderr := dao.GetDB().Where("follower = ?", user_id).Find(&relations).Error
	if relationsFinderr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	follows := []int64{}
	for _, relation := range relations {
		follows = append(follows, relation.Follow)
	}
	// 根据id在User数据库查找
	users := []service.User{}
	followsFindErr := dao.GetDB().Where("id IN ?", follows).Find(&users).Error
	if followsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 把users的IsFollow设为true（已关注）
	for i, _ := range users {
		users[i].IsFollow = true
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 查找用户
	var user service.User
	userExitErr := dao.GetDB().Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// 在Relation数据查找粉丝的id
	relations := []dao.Relation{}
	relationsFinderr := dao.GetDB().Where("follow = ?", user_id).Find(&relations).Error
	if relationsFinderr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	followers := []int64{}
	for _, relation := range relations {
		followers = append(followers, relation.Follower)
	}
	// 根据id在User数据库查找
	users := []service.User{}
	followersFindErr := dao.GetDB().Where("id IN ?", followers).Find(&users).Error
	if followersFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 把users的IsFollow设为false（未关注）
	for i, _ := range users {
		users[i].IsFollow = false
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 查找用户
	var user service.User
	userExitErr := dao.GetDB().Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// 在Relation数据查找和用户相关的记录（只查找一次）
	relations := []dao.Relation{}
	relationsFinderr := dao.GetDB().Where("follow = ? OR follower = ?", user_id, user_id).Find(&relations).Error
	if relationsFinderr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	// 提取用户的关注者id和粉丝id
	follows, followers := []int64{}, []int64{}
	userIdInt, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id非法"})
		return
	}
	for _, relation := range relations {
		if relation.Follow == userIdInt {
			followers = append(followers, relation.Follower)
		} else if relation.Follower == userIdInt {
			follows = append(follows, relation.Follow)
		}
	}
	// 寻找关注者和粉丝的共同者id
	friends := []int64{}
	for i, _ := range follows {
		for j, _ := range followers {
			if follows[i] == followers[j] {
				friends = append(friends, follows[i])
				break
			}
		}
	}
	// 根据id在User数据库查找好友
	users := []service.User{}
	friendsFindErr := dao.GetDB().Where("id IN ?", friends).Find(&users).Error
	if friendsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 将Relation数据库中用户为followers的MessageId字段全置为-1
	setRelations := []dao.Relation{}
	setRelationsFindErr := dao.GetDB().Where("follower = ?", user_id).Find(&setRelations).Error
	if setRelationsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	if len(setRelations) > 0 {
		for i, _ := range setRelations {
			setRelations[i].MessageId = -1
		}
		setRelationsSaveErr := dao.GetDB().Save(&setRelations).Error
		if setRelationsSaveErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库更新失败"})
			return
		}
	}
	// 查找与好友双方有关的最新消息
	var message service.Message
	friendUsers := []FriendUser{}
	for i, _ := range users {
		messageFindErr := dao.GetDB().Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)",
			userIdInt, users[i].Id, users[i].Id, userIdInt).Last(&message).Error // 暂时按照主键降序查找
		if messageFindErr == nil {
			msgType := 0
			if message.FromUserId == userIdInt {
				msgType = 1
			}
			friendUsers = append(friendUsers, FriendUser{Id: users[i].Id, Name: users[i].Name, FollowCount: users[i].FollowCount,
				FollowerCount: users[i].FollowerCount, IsFollow: true, FriendMessage: message.Content, MsgType: int64(msgType)}) // 把好友的IsFollow设为true（已关注）
		} else {
			friendUsers = append(friendUsers, FriendUser{Id: users[i].Id, Name: users[i].Name, FollowCount: users[i].FollowCount,
				FollowerCount: users[i].FollowerCount, IsFollow: true}) // 把好友的IsFollow设为true（已关注）
		}
	}

	c.JSON(http.StatusOK, FriendUserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		FriendUserList: friendUsers,
	})
}
