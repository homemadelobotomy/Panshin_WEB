package service

import "lab/internal/app/ds"

func(s *Service) GetSolarPanels()([]ds.SolarPanel, error) {
	return  s.repository.GetSolarPanels()
}