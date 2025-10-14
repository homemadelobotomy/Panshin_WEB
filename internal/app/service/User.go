package service

import (
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
)

func (s *Service) GetUserData(userId uint) (dto.UserDataResposne, error) {

	user, err := s.repository.GetUser(userId)
	if err != nil {
		return dto.UserDataResposne{}, err
	}

	return dto.UserDataResposne{Login: user.Login,
		ID: user.ID}, nil
}

func (s *Service) AddNewUser(user dto.UserRegistration) (dto.UserDataResposne, error) {
	userId, err := s.repository.AddNewUser(&ds.User{Login: user.Login,
		Password: user.Password})
	if err != nil {
		return dto.UserDataResposne{}, err
	}
	return s.GetUserData(userId)
}

func (s *Service) ChangeUserData(user dto.ChangeUserData) (dto.UserDataResposne, error) {
	err := s.repository.ChangeUserData(ds.GetUser().GetId(), user)
	if err != nil {
		return dto.UserDataResposne{}, err
	}
	response, err := s.repository.GetUser(ds.GetUser().GetId())
	if err != nil {
		return dto.UserDataResposne{}, err
	}
	return dto.UserDataResposne{ID: response.ID, Login: response.Login}, nil
}
