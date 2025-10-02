package dto

type SolarPanelFromRequestResponse struct {
	ID       uint    `json:"id"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	Power    float64 `json:"power"`
	Image    string  `json:"image"`
	IsDelete bool    `json:"is_deleted"`
	Area     float64 `json:"area"`
}

type AddSolarPanel struct {
	Title       string  `json:"title" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	Height      int     `json:"height" binding:"required"`
	Width       int     `json:"width" binding:"required"`
	Depth       int     `json:"depth" binding:"required"`
	Efficiency  string  `json:"efficiency" binding:"required"`
	Power       float64 `json:"power" binding:"required"`
}

type ChangeSolarPanel struct {
	Title       string  `json:"title,omitempty"`
	Type        string  `json:"type,omitempty"`
	Description string  `json:"description,omitempty"`
	Height      int     `json:"height,omitempty"`
	Width       int     `json:"width,omitempty"`
	Depth       int     `json:"depth,omitempty"`
	Efficiency  string  `json:"efficiency,omitempty"`
	Power       float64 `json:"power,omitempty"`
}
