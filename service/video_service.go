package service

type Video struct {
	Id            int64 `json:"id,omitempty" gorm:"primary_key"`
	Author        User  `json:"author" gorm:"foreignKey:Id;references:AuthorID;"`
	AuthorID      int64
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	PublishTime   int64
}
