package ds

import (
	"time"
)

type SolarPanelRequest struct {
	ID          uint      `gorm:"primaryKey"`
	Status      string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	CreatorId   uint      `gorm:"not null"`
	FormatedAt  time.Time `gorm:"default:NULL"`
	DeletedAt   time.Time `gorm:"default:NULL"`
	ModeratedAt time.Time `gorm:"default:NULL"`
	ModeratorID uint      `gorm:"default:NULL"`
	TotalPower  float64
	Insolation  float64

	Panels    []RequestPanels `gorm:"foreignKey:SolarPanelRequestID"`
	Creator   User
	Moderator User
}
