package service

type Like struct {
	UserID  uint `gorm:"userid"`
	VideoID uint `gorm:"videoid"`
}
