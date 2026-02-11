package models

type DeliverySummary struct {
	ID           string `json:"id"`
	RequestID    string `json:"requestId"`
	Status       string `json:"status"`
	TrackingCode string `json:"trackingCode"`
	Type         string `json:"type"`
	Garments     string `json:"garments"`
}
