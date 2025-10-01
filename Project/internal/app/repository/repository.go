package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type SolarPanel struct {
	ID          int
	Title       string
	Type        string
	Description string
	Power       int
	Size        string
	Efficiency  string
	Image       string
	Area        float64
}
type Order struct {
	ID    int
	Title string
}

func (r *Repository) GetSolarPanels() ([]SolarPanel, error) {
	panels := []SolarPanel{
		{
			ID:         1,
			Title:      "",
			Type:       "Монокристаллическая",
			Power:      50,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://127.0.0.1:9000/images/mono.jpeg",
			Area:       41,
		},
		{
			ID:         2,
			Title:      "Самая лучшая панель",
			Type:       "Поликристаллическая",
			Power:      280,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://127.0.0.1:9000/images/poly.jpeg",
			Area:       122,
		},
		{
			ID:         3,
			Title:      "",
			Type:       "Тонкопленочная",
			Power:      40,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://127.0.0.1:9000/images/tonko.jpeg",
			Area:       5,
		},
		{
			ID:         4,
			Title:      "",
			Type:       "Монокристаллическая",
			Power:      50,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://127.0.0.1:9000/images/mono.jpeg",
			Area:       12,
		},
	}
	if len(panels) == 0 {
		return nil, fmt.Errorf("array is empty")
	}
	return panels, nil
}

func (r *Repository) GetSolarPanel(id int) (SolarPanel, error) {
	panels, err := r.GetSolarPanels()

	if err != nil {
		logrus.Error("error while getting panel")
		return SolarPanel{}, err
	}
	for _, panel := range panels {
		if panel.ID == id {
			return panel, nil
		}
	}
	return SolarPanel{}, fmt.Errorf("can not find panel with id:%d", id)
}

type SolarPanelRequest struct {
	ID          int
	Insolation  float64
	Status      string
	TotalPower  float64
	SolarPanels []SolarPanel
	Area        float64
}

func (r *Repository) GetSolarPanelsInRange(begin int, end int) ([]SolarPanel, error) {
	panels, err := r.GetSolarPanels()
	if err != nil {
		return []SolarPanel{}, err
	}
	result := []SolarPanel{}
	for _, panel := range panels {
		if panel.Power >= begin && panel.Power <= end {
			result = append(result, panel)
		}
	}
	return result, nil
}
func (r *Repository) GetSolarPanelsRequest(id int) (SolarPanelRequest, error) {
	panel1, err := r.GetSolarPanel(1)
	if err != nil {
		logrus.Error(err)
	}
	panel2, err := r.GetSolarPanel(2)
	if err != nil {
		logrus.Error(err)
	}

	var solar_panels []SolarPanel
	solar_panels = append(solar_panels, panel1, panel2)

	solar_panel_request := SolarPanelRequest{
		ID:          1,
		Insolation:  22,
		Status:      "черновик",
		TotalPower:  30,
		SolarPanels: solar_panels,
	}
	SolarPanelRequests := make(map[int]SolarPanelRequest)
	SolarPanelRequests[1] = solar_panel_request

	if len(SolarPanelRequests) == 0 {
		return SolarPanelRequest{}, fmt.Errorf("solarpanel request is empty")
	}
	return SolarPanelRequests[id], nil
}
