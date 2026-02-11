package models

type Notification struct {
	ID      string `json:"id"`
	Titulo  string `json:"titulo"`
	Mensaje string `json:"mensaje"`
	Leida   bool   `json:"leida"`
	Fecha   string `json:"fecha"`
}
