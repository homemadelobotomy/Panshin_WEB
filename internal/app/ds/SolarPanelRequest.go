package ds

import (
	"time"
)

type SolarPanelRequest struct {
	ID             int       `gorm:"primaryKey"`
	Status         string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	CreatorId      int       `gorm:"not null"`
	FormationDate  time.Time
	DeleteDate     time.Time
	CompletionDate time.Time
	RejectionDate  time.Time
	ModeratorID    int
	TotalPower     float64
	Insolation     float64

	Panels []RequestPanels `gorm:"foreignKey:SolarPanelRequestID"`
}
