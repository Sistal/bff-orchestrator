package models

type CatalogItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Campaign struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	FechaInicio string `json:"fechaInicio"`
	FechaFin    string `json:"fechaFin"`
	Activa      bool   `json:"activa"`
}

type MasterDataResponse struct {
	Sizes         []CatalogItem `json:"sizes"`
	GarmentTypes  []CatalogItem `json:"garmentTypes"`
	ChangeReasons []CatalogItem `json:"changeReasons"`
	Campaign      Campaign      `json:"campaign"`
}
