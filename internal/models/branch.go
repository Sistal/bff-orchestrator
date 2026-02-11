package models

type Branch struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type BranchChangeRequestHistory struct {
	ID               int    `json:"id"`
	FechaSolicitud   string `json:"fechaSolicitud"`
	FechaEfectiva    string `json:"fechaEfectiva"`
	SucursalAnterior string `json:"sucursalAnterior"`
	SucursalNueva    string `json:"sucursalNueva"`
	Motivo           string `json:"motivo"`
	Estado           string `json:"estado"`
}

type CreateBranchChangeRequest struct {
	BranchID      int    `json:"branchId"`
	EffectiveDate string `json:"effectiveDate"`
	Reason        string `json:"reason"`
}
