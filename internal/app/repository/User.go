package repository

import "lab/internal/app/ds"

func (r *Repository) GetUser(userID uint) (ds.User, error) {
	var user ds.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}
