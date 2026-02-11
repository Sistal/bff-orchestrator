package models

type EmployeeProfile struct {
	ID      int    `json:"id"` // Added based on context, though not explicitly in the snippet showing structure
	Nombre  string `json:"nombre"`
	Email   string `json:"email"`
	Cargo   string `json:"cargo"`
	Celular string `json:"celular"`
	// Add other profile fields as needed
}

type UpdateContactRequest struct {
	Celular string `json:"celular"`
	Email   string `json:"email"`
}

type NotificationsPreferences struct {
	Email bool `json:"email"`
}

type UpdatePreferencesRequest struct {
	Notifications NotificationsPreferences `json:"notifications"`
}

type HomeStats struct {
	SolicitudesPendientes int `json:"solicitudes_pendientes"`
	EntregasProximas      int `json:"entregas_proximas"`
}

type ActivityLog struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type BodyMeasurements struct {
	EstaturaM float64 `json:"estatura_m"`
	PechoCm   float64 `json:"pecho_cm"`
	// Add other measurements
}
