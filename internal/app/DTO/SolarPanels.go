package dto

type SolarPanelResponse struct {
	ID       uint    `json:"id"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	Power    float64 `json:"power"`
	Image    string  `json:"image"`
	IsDelete bool    `json:"is_deleted"`
	Area     float64 `json:"area"`
}
