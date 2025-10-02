package dto

// type SolarPanelFromRequestResponse struct {
// 	SolarPanelId        uint    `json:"solar_panel_id"`
// 	SolarPanelRequestId uint    `json:"solar_panel_request_id"`
// 	Area                float64 `json:"area"`
// }

type ChangeSolarPanelAreaRequest struct {
	Area float64 `json:"area" binding:"required"`
}
