package ds

type RequestPanels struct {
	SolarPanelRequestID uint    `gorm:"primaryKey;auto_increment:false" json:"solar_panel_id"`
	SolarPanelID        uint    `gorm:"primaryKey;auto_increment:false" json:"solar_panel_request_id"`
	Area                float64 `json:"area"`

	SolarPanel SolarPanel `gorm:"foreignKey:SolarPanelID"`
}
