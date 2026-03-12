package models

// TimelineEvent — evento del timeline de una entrega.
// No existe en DDL; se diseña para visualización en detalle de entrega.
type TimelineEvent struct {
	Status    string `json:"status"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

// DeliverySummary — resumen de entrega.
// Mapea tabla Despacho con JOINs a Petición Uniforme, Estado y Tallaje+Prenda.
// NOTA: estimatedDate no existe en el DDL y fue excluida del contrato.
type DeliverySummary struct {
	ID           string          `json:"id"`
	RequestID    string          `json:"requestId"`          // "SOL-" + Petición Uniforme.id_peticion
	DispatchDate string          `json:"dispatchDate"`       // fecha_despacho date
	Garments     string          `json:"garments"`           // JOIN Tallaje + Prenda.nombre_prenda
	Address      string          `json:"address"`            // Despacho.sucursal varchar(100)
	Status       string          `json:"status"`             // Estado.nombre_estado ("in-transit"|"delivered")
	TrackingCode string          `json:"trackingCode"`       // guia_de_despacho varchar(100)
	Type         string          `json:"type"`               // Tipo Petición.nombre_tipo_peticion
	Timeline     []TimelineEvent `json:"timeline,omitempty"` // solo en GET /entregas/:id
}
