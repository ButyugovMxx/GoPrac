package domain

type Movieload struct {
	Id       int32  `json:"id"`
	Title    string `json:"title"`
	Year     int32  `json:"year"`
	Director string `json:"director"`
}

type MovieError struct {
	Id  int
	Err error
}

type Loadres struct {
	Movies []Movieload
	Errors []MovieError
}
