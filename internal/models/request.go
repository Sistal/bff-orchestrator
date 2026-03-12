package models

// RequestSummary — resumen de solicitud.
// Mapea tabla "Petición Uniforme" con JOINs a Tipo Petición, Estado y Tallaje+Prenda.
type RequestSummary struct {
	ID         string   `json:"id"`                   // "SOL-" + id_peticion
	Tipo       string   `json:"tipo"`                 // Tipo Petición.nombre_tipo_peticion
	Fecha      string   `json:"fecha"`                // fecha_registro date
	Estado     string   `json:"estado"`               // Estado.nombre_estado
	Items      []string `json:"items"`                // []Tallaje JOIN Prenda.nombre_prenda
	Motivo     string   `json:"motivo"`               // observación text
	NuevoTalle string   `json:"nuevoTalle,omitempty"` // Tallaje.valor_talla (solo cambio-prenda)
}

// CreateReplenishmentRequest — body de POST /solicitudes/reposicion
type CreateReplenishmentRequest struct {
	Items  []string `json:"items"`  // nombres de prendas
	Reason string   `json:"reason"` // → observación
}

// CreateGarmentChangeRequest — body de POST /solicitudes/cambio-prenda
type CreateGarmentChangeRequest struct {
	Prenda  string `json:"prenda"`  // nombre_prenda
	Reason  string `json:"reason"`  // → observación
	NewSize string `json:"newSize"` // → Tallaje.valor_talla
}

// FileUploadResponse — respuesta de POST /archivos/upload
type FileUploadResponse struct {
	FileID string `json:"fileId"`
	URL    string `json:"url"`
	Name   string `json:"name"`
}
