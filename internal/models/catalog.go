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

type TipoEmpresa struct {
	ID   int    `json:"id_tipo_empresa"`
	Name string `json:"nombre_tipo_empresa"`
}

type Empresa struct {
	ID                      int         `json:"id_empresa"`
	RazonSocial             string      `json:"razon_social"`
	IdentificadorTributario string      `json:"identificador_tributario"`
	Direccion               string      `json:"direccion"`
	Telefono                string      `json:"telefono"`
	Email                   string      `json:"email"`
	TipoEmpresa             TipoEmpresa `json:"tipo_empresa"`
}

type Segmento struct {
	ID          int    `json:"id_segmento"`
	Nombre      string `json:"nombre_segmento"`
	Descripcion string `json:"descripcion"`
}

type Sucursal struct {
	ID        int    `json:"id_sucursal"`
	Nombre    string `json:"nombre_sucursal"`
	Direccion string `json:"direccion"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages,omitempty"`
}

type CatalogResponse[T any] struct {
	Success bool  `json:"success"`
	Data    []T   `json:"data"`
	Meta    *Meta `json:"meta,omitempty"`
}
