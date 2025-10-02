package service

import (
	"context"
	"errors"
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
	"lab/internal/app/dsn"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func (s *Service) GetSolarPanels(startValue float64, endValue float64) ([]ds.SolarPanel, error) {
	if startValue < 0 || endValue < 0 {
		return nil, ErrBadRequest
	}

	panels, err := s.repository.GetFilteredSolarPanels(startValue, endValue)

	if err != nil {
		return nil, err
	}
	if len(panels) == 0 {
		return nil, ErrNoRecords
	}
	return panels, nil
}

func (s *Service) GetSolarPanel(panelId uint) (ds.SolarPanel, error) {
	solarPanel, err := s.repository.GetOneSolarPanel(panelId)
	if err != nil {
		return ds.SolarPanel{}, err
	}
	return solarPanel, nil
}

func (s *Service) AddNewSolarPanel(newSolarPanel dto.AddSolarPanel) (ds.SolarPanel, error) {

	solarPanelToBD := ds.SolarPanel{
		Title:       newSolarPanel.Title,
		Type:        newSolarPanel.Title,
		Description: newSolarPanel.Description,
		Power:       newSolarPanel.Power,
		Height:      newSolarPanel.Height,
		Width:       newSolarPanel.Width,
		Depth:       newSolarPanel.Depth,
		Efficiency:  newSolarPanel.Efficiency,
		Image:       "",
	}

	response, err := s.repository.AddNewSolarPanel(solarPanelToBD)
	if err != nil {
		return ds.SolarPanel{}, err
	}
	return response, nil

}

func (s *Service) ChangeSolarPanel(solarPanelId uint, solarPanel dto.ChangeSolarPanel) (ds.SolarPanel, error) {
	if solarPanel.Depth < 0 ||
		solarPanel.Height < 0 ||
		solarPanel.Power < 0 ||
		solarPanel.Width < 0 {
		return ds.SolarPanel{}, ErrBadRequest
	}
	err := s.repository.ChangeSolarPanel(solarPanelId, solarPanel)
	if err != nil {
		return ds.SolarPanel{}, err
	}
	return s.repository.GetOneSolarPanel(solarPanelId)
}

func (s *Service) DeleteSolarPanel(solarPanelId uint) error {
	return s.repository.DeleteSolarPanel(solarPanelId)
}

func (s *Service) AddSolarPanelToRequest(solarPanelId uint, userId uint) error {
	var solarPanelRequest ds.SolarPanelRequest
	solarPanel, err := s.repository.GetOneSolarPanel(solarPanelId)
	if err != nil {
		return ErrNoRecords
	}
	if solarPanel.IsDelete {
		return ErrSolarPanelDeleted
	}
	solarPanelRequestId, _, err := s.repository.GetSolarPanelsInRequest(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			solarPanelRequest = ds.SolarPanelRequest{
				CreatorId: userId,
				CreatedAt: time.Now(),
				Status:    "черновик",
			}
			createErr := s.repository.CreateNewSolarPanelRequest(&solarPanelRequest)
			if createErr != nil {
				return createErr
			}
		} else {
			return err
		}

	}
	solarPanelRequest, err = s.repository.GetOneSolarPanelRequest(solarPanelRequestId, "черновик")
	if err != nil {
		return err
	}
	if userId != solarPanelRequest.CreatorId {
		return ErrForbidden
	}
	return s.repository.AddSolarPanelToRequest(solarPanelId, solarPanelRequest.ID)

}

func (s *Service) UploadImageToMinio(file *multipart.FileHeader, filename string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = s.minioClient.PutObject(context.Background(),
		"images",
		filename,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})
	if err != nil {
		return "", err
	}
	imageURL := dsn.GetMinioURL(filename)

	return imageURL, nil
}

func (s *Service) AddImageToSolarPanel(solarPanelId uint, file *multipart.FileHeader, filename string) (ds.SolarPanel, error) {
	panel, err := s.repository.GetOneSolarPanel(solarPanelId)
	if err != nil {
		return ds.SolarPanel{}, err
	}

	if panel.Image != "" {
		oldFilename := extractFilenameFromURL(panel.Image)
		s.minioClient.RemoveObject(context.Background(), "images", oldFilename, minio.RemoveObjectOptions{})
	}
	newImageUrl, err := s.UploadImageToMinio(file, filename)
	if err != nil {
		return ds.SolarPanel{}, err
	}
	err = s.repository.AddImageToSolarPanel(solarPanelId, newImageUrl)
	if err != nil {
		return ds.SolarPanel{}, err
	}
	return s.repository.GetOneSolarPanel(solarPanelId)
}
