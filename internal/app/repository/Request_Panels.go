package repository

func (r *Repository) DeleteSolarPanelFromRequest(userId uint, requestId uint, solarPanelId uint) {
	//TODO Выполнить DELETE из таблицы request_panels
}

func (r *Repository) ChangeSolarPanelArea(userId uint, requestId uint, solarPanelId uint) {
	// Выполнить изменение площади у услуги в заявке
}
