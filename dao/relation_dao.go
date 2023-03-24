package dao

type Relation struct {
	Id        int64 `json:"id,omitempty"`
	Follow    int64 `json:"follow,omitempty"`
	Follower  int64 `json:"follower,omitempty"`
	MessageId int64 `json:"message_id,omitempty"` // 为了确认聊天是否是刚开始已经对方发的消息是否最新
}
