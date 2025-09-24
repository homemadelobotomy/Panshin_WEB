package ds

type SolarPanel struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(100);not null"`
	Type        string `gorm:"type:varchar(30);not null"`
	Description string `gorm:"type:varchar(500)"`
	Power       int    `gorm:"not null"`
	Size        string `gorm:"type:varchar(100);not null"`
	Efficiency  string `gorm:"type:varchar(10);not null"`
	Image       string `gorm:"type:varchar(100)"`
	IsDelete    bool   `gorm:"type:boolean not null;default:true"`
}
