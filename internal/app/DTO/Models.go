package dto

import "time"

type SolarPanleRequestFilter struct {
	Status     string
	Start_date time.Time
	End_date   time.Time
}
