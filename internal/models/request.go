package models

type RequestSummary struct {
	ID     string   `json:"id"`
	Tipo   string   `json:"tipo"`
	Fecha  string   `json:"fecha"`
	Estado string   `json:"estado"`
	Items  []string `json:"items"`
	Motivo string   `json:"motivo"`
}

type CreateReplenishmentRequest struct {
	Items  []string `json:"items"`
	Reason string   `json:"reason"`
}

type CreateGarmentChangeRequest struct {
	Prenda  string `json:"prenda"`
	Reason  string `json:"reason"`
	NewSize string `json:"newSize"`
}

type FileUploadResponse struct {
	FileID string `json:"fileId"`
	URL    string `json:"url"`
	Name   string `json:"name"`
}
