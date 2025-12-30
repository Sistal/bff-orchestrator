package model

// Dashboard represents aggregated data for the frontend
type Dashboard struct {
	User     *User      `json:"user"`
	Products []*Product `json:"products"`
}
