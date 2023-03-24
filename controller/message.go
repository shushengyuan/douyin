package controller

import (
	// "fmt"
	"douyin/dao"
	"douyin/service"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"

	// "sync/atomic"
	// "time"

	"github.com/gin-gonic/gin"
)

var tempChat = map[string][]service.Message{}
var chatConnMap = sync.Map{}

type ChatResponse struct {
	Response
	MessageList []service.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}
	// "to_user_id"转换为int64类型
	toUserIdInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id非法"})
		return
	}
	// 查找对方用户
	var toUser service.User
	toUserExitErr := db.Where("id = ?", toUserIdInt).Take(&toUser).Error
	if toUserExitErr == nil {
		CreateMessageErr := db.Create(&service.Message{ToUserId: toUserIdInt, FromUserId: user.Id, Content: content}).Error
		if CreateMessageErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "发送消息失败！"})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 0})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	var user dao.User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	// "to_user_id"转换为int64类型
	toUserIdInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id非法"})
		return
	}
	// 查找对方用户
	var toUser service.User
	toUserExitErr := db.Where("id = ?", toUserIdInt).Take(&toUser).Error
	if toUserExitErr == nil {
		// 查找relation数据库中用户与对方对应的MessageId
		var relation Relation
		relationFindErr := db.Where("follower = ? AND follow = ?", user.Id, toUserIdInt).Take(&relation).Error
		if relationFindErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查找relations数据库失败"})
			return
		}
		// -1表明对话刚刚开始，查找与双方有关的所有消息
		messages := []service.Message{}
		if relation.MessageId == -1 {
			messagesFindErr := db.Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)",
				toUserIdInt, user.Id, user.Id, toUserIdInt).Order("Id").Find(&messages).Error // 暂时按照主键顺序查找
			if messagesFindErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查找messages数据库失败"})
				return
			}

		} else { // 对话正在进行，只返回对方发来的最新消息
			messagesFindErr := db.Where("to_user_id = ? AND from_user_id = ? AND Id > ?",
				user.Id, toUserIdInt, relation.Id).Order("Id").Find(&messages).Error // 暂时按照主键降序查找
			if messagesFindErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查找messages数据库失败"})
				return
			}
		}

		// 若有新消息，则更新MessageId
		if len(messages) > 0 {
			relation.MessageId = messages[len(messages)-1].Id
			relationSaveErr := db.Save(&relation).Error
			if relationSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库更新失败"})
				return
			}
		}
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messages})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func RunMessageServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("Run message sever failed: %v\n", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept conn failed: %v\n", err)
			continue
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read message failed: %v\n", err)
			continue
		}

		var event = service.MessageSendEvent{}
		_ = json.Unmarshal(buf[:n], &event)
		fmt.Printf("Receive Message：%+v\n", event)

		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
		if len(event.MsgContent) == 0 {
			chatConnMap.Store(fromChatKey, conn)
			continue
		}

		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
		writeConn, exist := chatConnMap.Load(toChatKey)
		if !exist {
			fmt.Printf("User %d offline\n", event.ToUserId)
			continue
		}

		pushEvent := service.MessagePushEvent{
			FromUserId: event.UserId,
			MsgContent: event.MsgContent,
		}
		pushData, _ := json.Marshal(pushEvent)
		_, err = writeConn.(net.Conn).Write(pushData)
		if err != nil {
			fmt.Printf("Push message failed: %v\n", err)
		}
	}
}
