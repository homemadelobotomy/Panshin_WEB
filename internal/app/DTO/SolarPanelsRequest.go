package dto

type NumberOfPanelsResponse struct {
	RequestId      uint  `json:"request_id"`
	NumberOfPanels int64 `json:"panels_in_request"`
}

type SolarPanelsRequestsResponse struct {
	ID          uint    `json:"id"`
	Status      string  `json:"status"`
	Creator     string  `json:"creator"`
	CreatedAt   string  `json:"created_at"`
	FormatedAt  string  `json:"formated_at"`
	ModeratedAt string  `json:"moderated_at"`
	Moderator   string  `json:"moderator"`
	TotalPower  float64 `json:"total_power"`
	Insolation  float64 `json:"insolation"`
}

type OneSolarPanelRequestResponse struct {
	ID         uint                            `json:"id"`
	TotalPower float64                         `json:"total_power"`
	Insolation float64                         `json:"insolation"`
	Panels     []SolarPanelFromRequestResponse `json:"solarpanels"`
}

type ChangeSolarPanelRequest struct {
	Insolation float64 `json:"insolation" binding:"required"`
}

type ModeratorAction struct {
	Action string `json:"action" binding:"required"`
}
