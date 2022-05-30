package models

// table utilisateur
type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	KmMax  int64  `json:"kmmax"`
	Niveau string `json:"niveau"`
}
