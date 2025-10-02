package repository

import (
	"errors"
	"fmt"
	"time"

	"lab/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetSolarPanels() ([]ds.SolarPanel, error) {
	var panels []ds.SolarPanel

	err := r.db.Find(&panels).Error

	if err != nil {
		return nil, err
	}

	if len(panels) == 0 {
		return nil, fmt.Errorf("array is empty")
	}

	return panels, nil

}

func (r *Repository) GetSolarPanel(id int) (*ds.SolarPanel, error) {
	panel := &ds.SolarPanel{}

	err := r.db.Where("id = ?", id).First(&panel).Error
	if err != nil {
		return &ds.SolarPanel{}, err
	}
	return panel, nil
}

func (r *Repository) GetSolarPanelsInRange(begin int, end int) ([]ds.SolarPanel, error) {
	panels := []ds.SolarPanel{}

	err := r.db.Where("power BETWEEN ? AND ?", begin, end).Find(&panels).Error

	if err != nil {
		return nil, err
	}
	return panels, nil
}

func (r *Repository) GetSolarPanelRequestID() (int, error) {
	user_id := 1
	solarpanel_request := &ds.SolarPanelRequest{}
	err := r.db.Where("creator_id = ? AND status = 'черновик'", user_id).Preload("Panels.SolarPanel").First(&solarpanel_request).Error
	if err != nil {
		return -1, err
	}
	return solarpanel_request.ID, nil
}

func (r *Repository) GetSolarPanelRequest(id int) (*ds.SolarPanelRequest, error) {
	solarpanel_request := &ds.SolarPanelRequest{}

	err := r.db.Where("id = ? AND status = 'черновик'", id).Preload("Panels.SolarPanel").First(&solarpanel_request).Error
	if err != nil {
		return nil, err
	}
	return solarpanel_request, nil
}

func (r *Repository) GetNumberOfPanelsInRequest() (int64, error) {
	userId := 1
	var requestId int
	var count int64
	err := r.db.Model(&ds.SolarPanelRequest{}).
		Where("creator_id = ? AND status = 'черновик'", userId).
		Select("id").
		First(&requestId).Error

	if err != nil {
		return 0, fmt.Errorf("request not found for user %d : %w", userId, err)
	}

	err = r.db.Model(&ds.RequestPanels{}).
		Where("solar_panel_request_id = ?", requestId).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("can't count panels: %w", err)
	}

	return count, nil
}

func (r *Repository) AddSolarPanelToRequest(solarpanel_id int, user_id int) error {
	var request ds.SolarPanelRequest
	err := r.db.Where("creator_id = ? AND status = 'черновик'", user_id).
		First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			request = ds.SolarPanelRequest{
				CreatorId: user_id,
				Status:    "черновик",
				CreatedAt: time.Now(),
			}
			if err := r.db.Create(&request).Error; err != nil {
				return fmt.Errorf("error while creating new request: %w", err)
			}
		} else {
			return fmt.Errorf("error getting solar request: %w", err)
		}
	}
	request_panels := ds.RequestPanels{
		SolarPanelRequestID: uint(request.ID),
		SolarPanelID:        uint(solarpanel_id),
	}
	if err := r.db.Create(&request_panels).Error; err != nil {
		return fmt.Errorf("error adding new solar panel to request: %w", err)
	}
	return nil
}

func (r *Repository) DeleteSolarPanelRequest(request_id int) error {
	err := r.db.Exec("UPDATE solar_panel_requests SET status='удален', formated_at = $1 WHERE id=$2 AND status='черновик'", time.Now(), request_id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CalculateTotalPower(request_id int) {
	solar_panel_request, err := r.GetSolarPanelRequest(request_id)
	if err != nil {
		logrus.Error(err)
	}
	panels := solar_panel_request.Panels
	insolation := solar_panel_request.Insolation
	power := 0.0
	for _, panel := range panels {
		power += float64(panel.SolarPanel.Power) * panel.Area / (float64(panel.SolarPanel.Width * panel.SolarPanel.Height))
	}

	total_power := power * insolation / 1000

	r.db.Model(&ds.SolarPanelRequest{}).Where("id = ?", request_id).Update("total_power", total_power)
}
