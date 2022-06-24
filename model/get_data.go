package model

type AccountGet struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Type   *int   `json:"type"`
	Status *int   `json:"status"`
}
