package repository

import (
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
)

func (r *Repository) GetFilteredSolarPanels(startValue float64, endValue float64) ([]ds.SolarPanel, error) {
	//TODO Вернкть список панелей, если понадобится, то отсортированный. Обязательно проверить что выведет ctx.Param("start"), при пустых параметрах в запросе,
	// если ошибку то можно ее обработать, но лучше разделить логику этого метода на два, с фильтрацией и без
	var panels []ds.SolarPanel
	db := r.db.Model(&ds.SolarPanel{}).Where("is_delete = false")
	if startValue > 0 {
		db = db.Where("power > ?", startValue)
	}
	if endValue > 0 {
		db = db.Where("power < ?", endValue)
	}
	err := db.Find(&panels).Error
	if err != nil {
		return nil, err
	}
	return panels, nil

}

func (r *Repository) GetOneSolarPanel(panelId uint) (ds.SolarPanel, error) {
	//TODO Вернуть одну панель по ее id
	var solarPanel ds.SolarPanel
	err := r.db.Where("id = ?", panelId).
		First(&solarPanel).Error
	if err != nil {
		return ds.SolarPanel{}, err
	}
	return solarPanel, nil
}

func (r *Repository) AddNewSolarPanel(solarPanel ds.SolarPanel) (ds.SolarPanel, error) {
	//TODO Добавить новую панель без изображения
	err := r.db.Create(&solarPanel).Error
	if err != nil {
		return ds.SolarPanel{}, err
	}
	return solarPanel, nil
}

func (r *Repository) ChangeSolarPanel(solarPanelId uint, solarPanel dto.ChangeSolarPanel) error {
	err := r.db.Model(&ds.SolarPanel{}).
		Where("id = ?", solarPanelId).
		Updates(&solarPanel).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteSolarPanel(solarPanelId uint) error {
	err := r.db.Model(&ds.SolarPanel{}).
		Where("id = ?", solarPanelId).
		Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateNewSolarPanelRequest(solarPanelRequest *ds.SolarPanelRequest) error {
	err := r.db.Create(solarPanelRequest).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddSolarPanelToRequest(solarPanelId uint, requestId uint) error {
	err := r.db.Create(&ds.RequestPanels{
		SolarPanelRequestID: requestId,
		SolarPanelID:        solarPanelId,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddImageToSolarPanel(solarPanelId uint, imageURL string) error {
	return r.db.Model(&ds.SolarPanel{}).
		Where("id = ?", solarPanelId).
		Update("image", imageURL).Error

}

// func (r *Repository) GetSolarPanel(id int) (*ds.SolarPanel, error) {
// 	panel := &ds.SolarPanel{}

// 	err := r.db.Where("id = ?", id).First(&panel).Error
// 	if err != nil {
// 		return &ds.SolarPanel{}, err
// 	}
// 	return panel, nil
// }

// func (r *Repository) GetSolarPanelsInRange(begin int, end int) ([]ds.SolarPanel, error) {
// 	panels := []ds.SolarPanel{}

// 	err := r.db.Where("power BETWEEN ? AND ?", begin, end).Find(&panels).Error

// 	if err != nil {
// 		return nil, err
// 	}
// 	return panels, nil
// }

// func (r *Repository) GetSolarPanelRequestID() (uint, error) {
// 	user_id := 1
// 	solarpanel_request := &ds.SolarPanelRequest{}
// 	err := r.db.Where("creator_id = ? AND status = 'черновик'", user_id).Preload("Panels.SolarPanel").First(&solarpanel_request).Error
// 	if err != nil {
// 		return 0, err
// 	}
// 	return solarpanel_request.ID, nil
// }

// func (r *Repository) GetSolarPanelRequest(id int) (*ds.SolarPanelRequest, error) {
// 	solarpanel_request := &ds.SolarPanelRequest{}

// 	err := r.db.Where("id = ? AND status = 'черновик'", id).Preload("Panels.SolarPanel").First(&solarpanel_request).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return solarpanel_request, nil
// }

// func (r *Repository) GetNumberOfPanelsInRequest() (int64, error) {
// 	var userId = ds.GetUser().GetId()
// 	var requestId int
// 	var count int64
// 	err := r.db.Model(&ds.SolarPanelRequest{}).
// 		Where("creator_id = ? AND status = 'черновик'", userId).
// 		Select("id").
// 		First(&requestId).Error

// 	if err != nil {
// 		return 0, fmt.Errorf("request not found for user %d : %w", userId, err)
// 	}

// 	err = r.db.Model(&ds.RequestPanels{}).
// 		Where("solar_panel_request_id = ?", requestId).
// 		Count(&count).Error

// 	if err != nil {
// 		return 0, fmt.Errorf("can't count panels: %w", err)
// 	}

// 	return count, nil
// }

// func (r *Repository) AddSolarPanelToRequest(solarpanel_id int) error {
// 	userId := ds.GetUser().GetId()
// 	var request ds.SolarPanelRequest
// 	err := r.db.Where("creator_id = ? AND status = 'черновик'", userId).
// 		First(&request).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			request = ds.SolarPanelRequest{
// 				CreatorId: userId,
// 				Status:    "черновик",
// 				CreatedAt: time.Now(),
// 			}
// 			if err := r.db.Create(&request).Error; err != nil {
// 				return fmt.Errorf("error while creating new request: %w", err)
// 			}
// 		} else {
// 			return fmt.Errorf("error getting solar request: %w", err)
// 		}
// 	}
// 	request_panels := ds.RequestPanels{
// 		SolarPanelRequestID: uint(request.ID),
// 		SolarPanelID:        uint(solarpanel_id),
// 	}
// 	if err := r.db.Create(&request_panels).Error; err != nil {
// 		return fmt.Errorf("error adding new solar panel to request: %w", err)
// 	}
// 	return nil
// }

// // func (r *Repository) DeleteSolarPanelRequest(request_id int) error {
// // 	err := r.db.Exec("UPDATE solar_panel_requests SET status='удален', delete_date = $1 WHERE id=$2 AND status='черновик'", time.Now(), request_id).Error
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }

// func (r *Repository) CalculateTotalPower(request_id int) {
// 	solar_panel_request, err := r.GetSolarPanelRequest(request_id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	panels := solar_panel_request.Panels
// 	insolation := solar_panel_request.Insolation
// 	power := 0.0
// 	for _, panel := range panels {
// 		power += float64(panel.SolarPanel.Power) * panel.Area / (float64(panel.SolarPanel.Width * panel.SolarPanel.Height))
// 	}

// 	total_power := power * insolation / 1000

// 	r.db.Model(&ds.SolarPanelRequest{}).Where("id = ?", request_id).Update("total_power", total_power)
// }
