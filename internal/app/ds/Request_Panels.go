package ds

type RequestPanels struct {
	SolarPanelRequestID uint `gorm:"primaryKey;auto_increment:false"`
	SolarPanelID        uint `gorm:"primaryKey;auto_increment:false"`
	Area                float64

	SolarPanel SolarPanel `gorm:"foreignKey:SolarPanelID"`
}
