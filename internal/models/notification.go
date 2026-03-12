package models

// Notification — mapea entidad de notificación del ms-operations.
// Campos en inglés para alinearse con el contrato frontend.
type Notification struct {
	ID          string `json:"id"`
	Type        string `json:"type"`        // "approved" | "delivery" | "alert" | "update"
	Title       string `json:"title"`       // antes: titulo
	Message     string `json:"message"`     // antes: mensaje
	Timestamp   string `json:"timestamp"`   // antes: fecha — ISO 8601 o formato relativo
	IsRead      bool   `json:"isRead"`      // antes: leida
	ActionLabel string `json:"actionLabel"` // etiqueta del CTA, opcional
}

// MarkAllReadResponse — respuesta de PATCH /notificaciones/leer-todas
type MarkAllReadResponse struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
}
