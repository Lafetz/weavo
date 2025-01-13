package domain

type Location struct {
	Id          string
	UserID      string
	Notes       string
	Nickname    string
	City        string
	Coordinates Coordinates
	CreatedAt   string
}
type Coordinates struct {
	Lat float64
	Lon float64
}
