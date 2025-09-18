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

type Solarpanel struct {
	ID          int
	Title       string
	Type        string
	Description string
	Power       int
	Size        string
	Efficiency  string
	Image       string
}
type Order struct {
	ID    int
	Title string
}

func (r *Repository) GetSolarpanels() ([]Solarpanel, error) {
	panels := []Solarpanel{
		{
			ID:         1,
			Title:      "",
			Type:       "Монокристаллическая",
			Power:      50,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://172.19.0.2:9000/lab1/mono.jpeg",
		},
		{
			ID:         2,
			Title:      "",
			Type:       "Поликристаллическая",
			Power:      280,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://172.19.0.2:9000/lab1/poly.jpeg",
		},
		{
			ID:         3,
			Title:      "",
			Type:       "Тонкопленочная",
			Power:      40,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://172.19.0.2:9000/lab1/tonko.jpeg",
		},
		{
			ID:         4,
			Title:      "",
			Type:       "Монокристаллическая",
			Power:      50,
			Size:       "1640×992×35",
			Efficiency: "20-22",
			Image:      "http://172.19.0.2:9000/lab1/mono.jpeg",
		},
	}
	if len(panels) == 0 {
		return nil, fmt.Errorf("array is empty")
	}
	return panels, nil
}

func (r *Repository) GetSolarpanel(id int) (Solarpanel, error) {
	panels, err := r.GetSolarpanels()

	if err != nil {
		logrus.Error("error while getting panel")
		return Solarpanel{}, err
	}
	for _, panel := range panels {
		if panel.ID == id {
			return panel, nil
		}
	}
	return Solarpanel{}, fmt.Errorf("can not find panel with id:%d", id)
}

type panelBid struct {
	Solarpanel
	Area float64
}

func (r *Repository) GetSolarPanelsInRange(begin int, end int) ([]Solarpanel, error) {
	panels, err := r.GetSolarpanels()
	if err != nil {
		return []Solarpanel{}, err
	}
	result := []Solarpanel{}
	for _, panel := range panels {
		if panel.Power >= begin && panel.Power <= end {
			result = append(result, panel)
		}
	}
	return result, nil
}
func (r *Repository) GetBid(id int) ([]panelBid, error) {
	bid := []panelBid{
		{
			Solarpanel: Solarpanel{
				ID:         1,
				Title:      "",
				Type:       "Монокристаллическая",
				Power:      50,
				Size:       "1640×992×35",
				Efficiency: "20-22",
				Image:      "http://172.19.0.2:9000/lab1/mono.jpeg",
			},
			Area: 0,
		},
		{
			Solarpanel: Solarpanel{
				ID:         3,
				Title:      "",
				Type:       "Тонкопленочная",
				Power:      40,
				Size:       "1640×992×35",
				Efficiency: "20-22",
				Image:      "http://172.19.0.2:9000/lab1/tonko.jpeg",
			},
			Area: 0,
		},
	}

	if len(bid) == 0 {
		return nil, fmt.Errorf("bid is empty")
	}
	return bid, nil
}
