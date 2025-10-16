package ds

import (
	"lab/internal/app/role"
)

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Login       string    `gorm:"type:varchar(255);unique"`
	Password    string    `gorm:"type:varchar(255)"`
	IsModerator role.Role `gorm:"default:false"`
}
