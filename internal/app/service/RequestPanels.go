package service

import (
	"lab/internal/app/ds"
	"math"
)

func (s *Service) DeleteSolarPanelFromRequest(userId uint, requestId uint, solarPanelId uint) error {

	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return err
	}

	if solarPanelRequest.CreatorId != userId {
		return ErrForbidden
	}
	exists := false
	for _, panel := range solarPanelRequest.Panels {
		if panel.SolarPanel.ID == solarPanelId {
			exists = true
			continue
		}
	}
	if !exists {
		return ErrNoRecords
	}
	err = s.repository.DeleteSolarPanelFromRequest(requestId, solarPanelId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeSolarPanelAreaInRequest(userId uint, requestId uint, solarPanelId uint, area float64) (ds.RequestPanels, error) {
	if math.IsNaN(area) || area <= 0 {
		return ds.RequestPanels{}, ErrBadRequest
	}
	solarPanelRequest, err := s.repository.GetOneSolarPanelRequest(requestId, "черновик")
	if err != nil {
		return ds.RequestPanels{}, err
	}
	if solarPanelRequest.CreatorId != userId {
		return ds.RequestPanels{}, ErrForbidden
	}
	err = s.repository.ChangeSolarPanelArea(requestId, solarPanelId, area)
	if err != nil {
		return ds.RequestPanels{}, ErrNoRecords
	}
	response, err := s.repository.GetSolarPanelFromRequest(requestId, solarPanelId)
	if err != nil {
		return ds.RequestPanels{}, ErrNoRecords
	}
	return response, nil

}
