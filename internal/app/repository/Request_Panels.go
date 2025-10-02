package repository

import "lab/internal/app/ds"

func (r *Repository) DeleteSolarPanelFromRequest(requestId uint, solarPanelId uint) error {
	//TODO Выполнить DELETE из таблицы request_panels
	err := r.db.
		Where("solar_panel_request_id = ? AND solar_panel_id = ?", requestId, solarPanelId).
		Delete(&ds.RequestPanels{}).Error
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) ChangeSolarPanelArea(requestId uint, solarPanelId uint, area float64) error {
	// Выполнить изменение площади у услуги в заявке
	err := r.db.Model(&ds.RequestPanels{}).
		Where("solar_panel_request_id = ? AND solar_panel_id = ?", requestId, solarPanelId).
		Update("area", area).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetSolarPanelFromRequest(requestId uint, solarPanelId uint) (ds.RequestPanels, error) {
	var requestPanel ds.RequestPanels
	err := r.db.Where("solar_panel_request_id = ? AND solar_panel_id = ?", requestId, solarPanelId).Preload("SolarPanel").
		First(&requestPanel).Error
	if err != nil {
		return ds.RequestPanels{}, err
	}
	return requestPanel, nil
}
