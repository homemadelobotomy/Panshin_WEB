package service

import (
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
	"math"
)

func (s *Service) GetSolarPanelsInRequest(userId uint) (uint, int64, error) {
	return s.repository.GetSolarPanelsInRequest(userId)
}

func (s *Service) GetFilteredSolarPanelRequests(
	userId uint,
	filter dto.SolarPanleRequestFilter,
) ([]dto.SolarPanelsRequestsResponse, error) {
	solarPanelRequests, err := s.repository.GetFilteredSolarPanelRequests(userId, filter)
	var solarPanelsRequestsResponse []dto.SolarPanelsRequestsResponse
	layout := "02-01-2006 15:04:05"
	if err != nil {
		return nil, err
	}
	if len(solarPanelRequests) == 0 {
		return nil, ErrNoRecords
	}

	for _, panelRequest := range solarPanelRequests {
		solarPanelsRequestsResponse = append(solarPanelsRequestsResponse, dto.SolarPanelsRequestsResponse{
			ID:          panelRequest.ID,
			Status:      panelRequest.Status,
			Creator:     panelRequest.Creator.Login,
			CreatedAt:   formateDate(panelRequest.CreatedAt, layout),
			FormatedAt:  formateDate(panelRequest.FormatedAt, layout),
			ModeratedAt: formateDate(panelRequest.ModeratedAt, layout),
			Moderator:   panelRequest.Moderator.Login,
			TotalPower:  panelRequest.TotalPower,
			Insolation:  panelRequest.Insolation,
		})
	}
	return solarPanelsRequestsResponse, nil
}

func (s *Service) GetOneSolarPanelRequest(requestId uint, userId uint) (dto.OneSolarPanelRequestResponse, error) {
	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}

	if userId != solarPanelRequest.CreatorId {
		return dto.OneSolarPanelRequestResponse{}, ErrForbidden
	}
	return s.ValidateSolarPanelRequestResponse(&solarPanelRequest), nil
}

func (s *Service) ValidateSolarPanelRequestResponse(solarPanelRequest *ds.SolarPanelRequest) dto.OneSolarPanelRequestResponse {
	var solarPanelsDTO []dto.SolarPanelFromRequestResponse
	for _, requestPanel := range solarPanelRequest.Panels {
		if !requestPanel.SolarPanel.IsDelete {
			solarPanelsDTO = append(solarPanelsDTO, dto.SolarPanelFromRequestResponse{
				ID:       requestPanel.SolarPanel.ID,
				Title:    requestPanel.SolarPanel.Title,
				Type:     requestPanel.SolarPanel.Type,
				Power:    requestPanel.SolarPanel.Power,
				Image:    requestPanel.SolarPanel.Image,
				IsDelete: requestPanel.SolarPanel.IsDelete,
				Area:     requestPanel.Area,
			})
		}
	}

	response := dto.OneSolarPanelRequestResponse{
		ID:         solarPanelRequest.ID,
		TotalPower: solarPanelRequest.TotalPower,
		Insolation: solarPanelRequest.Insolation,
		Panels:     solarPanelsDTO,
	}
	return response
}

func (s *Service) ChangeSolarPanelRequest(userId uint, requestId uint, insolation float64) (dto.OneSolarPanelRequestResponse, error) {
	if insolation < 0.0 || insolation > 10.0 {
		return dto.OneSolarPanelRequestResponse{}, ErrBadRequest
	}
	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	if solarPanelRequest.CreatorId != userId {
		return dto.OneSolarPanelRequestResponse{}, ErrForbidden
	}
	err = s.repository.ChangeSolarPanelRequest(requestId, insolation)
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	solarPanelRequest, err = s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	return s.ValidateSolarPanelRequestResponse(&solarPanelRequest), nil

}

func (s *Service) FormateSolarPanelRequest(requestId uint, userId uint) (dto.OneSolarPanelRequestResponse, error) {
	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	if solarPanelRequest.CreatorId != userId {
		return dto.OneSolarPanelRequestResponse{}, ErrForbidden
	}

	solarpanels := solarPanelRequest.Panels

	for _, panel := range solarpanels {
		if panel.Area <= 0 || math.IsNaN(panel.Area) {
			return dto.OneSolarPanelRequestResponse{}, ErrBadRequest
		}
	}
	if solarPanelRequest.Insolation <= 0 ||
		solarPanelRequest.Insolation > 10 ||
		math.IsNaN(solarPanelRequest.Insolation) {
		return dto.OneSolarPanelRequestResponse{}, ErrBadRequest
	}

	err = s.repository.FormateSolarPanelRequest(requestId)
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	solarPanelRequest, err = s.repository.GetOneSolarPanelRequest(requestId, "сформирован")
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	return s.ValidateSolarPanelRequestResponse(&solarPanelRequest), nil
}

func (s *Service) ModeratorAction(requestId uint, action string, moderatorId uint) (dto.OneSolarPanelRequestResponse, error) {
	user, err := s.repository.GetUser(moderatorId)
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	if !user.IsModerator {
		return dto.OneSolarPanelRequestResponse{}, ErrForbidden
	}
	if action != "завершен" && action != "отклонен" {
		return dto.OneSolarPanelRequestResponse{}, ErrBadRequest
	}
	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "сформирован")
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	panels := solarPanelRequest.Panels
	insolation := solarPanelRequest.Insolation
	totalPower := 0.0
	if action == "завершен" {
		totalPower = CalculateTotalPower(panels, insolation)
	}
	err = s.repository.ModeratorAction(requestId, action, totalPower, moderatorId)
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	solarPanelRequest, err = s.repository.GetOneSolarPanelRequest(requestId, action)
	if err != nil {
		return dto.OneSolarPanelRequestResponse{}, err
	}
	return s.ValidateSolarPanelRequestResponse(&solarPanelRequest), nil
}

func (s *Service) DeleteSolarPanelRequest(requestId uint, userId uint) error {
	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return err
	}
	if solarPanelRequest.CreatorId != userId {
		return ErrForbidden
	}
	err = s.repository.DeleteSolarPanelRequest(requestId)
	if err != nil {
		return err
	}
	return nil
}
