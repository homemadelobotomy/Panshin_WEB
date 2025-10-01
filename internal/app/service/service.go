package service

import (
	"lab/internal/app/repository"
	"time"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {

	return &Service{
		repository: repository,
	}
}

func formateDate(date time.Time, layout string) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(layout)
}
