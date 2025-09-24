package repository

import (
	"fmt"

	"lab/internal/app/ds"

	"github.com/sirupsen/logrus"
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

func (r *Repository) GetSolarPanelRequest(id int) (*ds.SolarPanelRequest, error) {
	solarpanel_request := &ds.SolarPanelRequest{}
	userId := 1
	err := r.db.Where("creator_id = ? AND status = 'черновик'", userId).Preload("Panels.SolarPanel").First(&solarpanel_request).Error
	if err != nil {
		return nil, err
	}
	return solarpanel_request, nil
}

func (r *Repository) GetNumberOfPanelsInRequest(id int) int64 {
	userId := 1
	var requestId int
	var count int64
	err := r.db.Model(&ds.SolarPanelRequest{}).Where("creator_id = ? AND status = 'черновик'", userId).Select("id").First(&requestId).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.RequestPanels{}).Where("solar_panel_request_id = ?", requestId).Count(&count).Error

	if err != nil {
		logrus.Println("Error counting records in lists_chats:", err)
	}
	return count
}

func (r* Repository) AddSolarPanelToRequest(request_id int, solarpanel_id int) error {
	
}