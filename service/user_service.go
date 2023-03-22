package service

type User struct {
	Id             int64  `json:"id,omitempty" gorm:"primary_key"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	IsFollow       bool   `json:"is_follow,omitempty"`
	TotalFavorited int64  `json:"total_favorited,omitempty" gorm:"default:0"`
	WorkCount      int64  `json:"work_count,omitempty" gorm:"default:0"`
	FavoriteCount  int64  `json:"favorite_count,omitempty" gorm:"default:0"`
}
