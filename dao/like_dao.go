package dao

type Like struct {
	UserID  uint `gorm:"userid"`
	VideoID uint `gorm:"videoid"`
}
