package domain

import "time"

type Location struct {
	Id          string
	UserID      string
	Notes       string
	Nickname    string
	City        string
	Coordinates Coordinates
	CreatedAt   time.Time
}
type Coordinates struct {
	Lat float64
	Lon float64
}
