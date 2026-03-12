package models

// Branch — mapea tabla Sucursal.
// NOTA: el DDL no tiene columna "region"; se expone "direccion" en su lugar.
type Branch struct {
	ID        int    `json:"id"`        // id_sucursal
	Name      string `json:"name"`      // nombre_sucursal
	Direccion string `json:"direccion"` // direccion varchar(255)
}

// BranchChangeRequestHistory — historial de solicitudes de cambio de sucursal.
// No existe tabla dedicada en el DDL; el ms-funcionario provee esta información.
type BranchChangeRequestHistory struct {
	ID               int    `json:"id"`
	FechaSolicitud   string `json:"fechaSolicitud"`
	FechaEfectiva    string `json:"fechaEfectiva"`
	SucursalAnterior string `json:"sucursalAnterior"`
	SucursalNueva    string `json:"sucursalNueva"`
	Motivo           string `json:"motivo"`
	Estado           string `json:"estado"`
}

// CreateBranchChangeRequest — body de POST /solicitudes/cambio-sucursal
type CreateBranchChangeRequest struct {
	BranchID      int    `json:"branchId"`
	EffectiveDate string `json:"effectiveDate"`
	Reason        string `json:"reason"`
}
