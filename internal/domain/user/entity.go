package user

type User struct {
	ID          string `gorm:"primaryKey"`
	Username    string `gorm:"uniqueIndex;not null"`
	PassHash    string `gorm:"not null"`
	HasUsername bool
	Bio         string
	IsActive    bool
}
