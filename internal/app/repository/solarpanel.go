package repository

import (
	"fmt"

	"lab/internal/app/ds"
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
