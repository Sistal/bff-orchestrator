package models

// CargoRef — referencia a tabla cargo (JOIN)
type CargoRef struct {
	ID          int    `json:"id"`
	NombreCargo string `json:"nombre_cargo"`
}

// SucursalRef — referencia a tabla Sucursal (JOIN)
// NOTA: el DDL no tiene columna "region"; se usa direccion para contexto geográfico.
type SucursalRef struct {
	ID             int    `json:"id"`
	NombreSucursal string `json:"nombre_sucursal"`
	Direccion      string `json:"direccion"`
}

// EstadoRef — referencia a tabla Estado (JOIN)
type EstadoRef struct {
	ID           int    `json:"id"`
	NombreEstado string `json:"nombre_estado"`
}

// NotificationPreferences — preferencias de notificación del funcionario
// No existe en DDL; se gestiona como mock temporal sin persistencia.
type NotificationPreferences struct {
	Email bool `json:"email"`
	Push  bool `json:"push"`
	SMS   bool `json:"sms"`
}

// UserPreferences — preferencias de la aplicación
// No existe en DDL; se gestiona como mock temporal sin persistencia.
type UserPreferences struct {
	Notifications NotificationPreferences `json:"notifications"`
	Theme         string                  `json:"theme"`
}

// EmployeeProfile — perfil completo del funcionario.
// Mapea tabla Funcionario con JOINs a cargo, Sucursal y Estado.
type EmployeeProfile struct {
	ID              int             `json:"id"`               // id_funcionario
	RutFuncionario  string          `json:"rut_funcionario"`  // rut_funcionario
	Nombres         string          `json:"nombres"`          // nombres
	ApellidoPaterno string          `json:"apellido_paterno"` // apellido_paterno
	ApellidoMaterno string          `json:"apellido_materno"` // apellido_materno
	Email           string          `json:"email"`            // email (contacto, puede diferir del email de login)
	Celular         string          `json:"celular"`          // celular
	Telefono        string          `json:"telefono"`         // telefono
	Direccion       string          `json:"direccion"`        // direccion
	Cargo           CargoRef        `json:"cargo"`            // JOIN cargo
	Sucursal        SucursalRef     `json:"sucursal"`         // JOIN Sucursal
	Estado          EstadoRef       `json:"estado"`           // JOIN Estado
	Preferences     UserPreferences `json:"preferences"`      // mock temporal, sin persistencia
}

// UpdateContactRequest — campos editables de Funcionario (PUT /api/v1/funcionarios/me)
type UpdateContactRequest struct {
	Nombres         string `json:"nombres"`
	ApellidoPaterno string `json:"apellido_paterno"`
	ApellidoMaterno string `json:"apellido_materno"`
	Celular         string `json:"celular"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
	Direccion       string `json:"direccion"`
}

// UpdatePreferencesRequest — preferencias (mock temporal sin persistencia en DDL)
type UpdatePreferencesRequest struct {
	Notifications NotificationPreferences `json:"notifications"`
	Theme         string                  `json:"theme"`
}

// UpdateSecurityRequest — configuración de seguridad
// No existe en DDL; se reserva para futuras implementaciones.
type UpdateSecurityRequest struct {
	RecoveryEmail string `json:"recoveryEmail"`
}

// HomeStats — estadísticas del dashboard agregadas desde Petición Uniforme y Despacho
type HomeStats struct {
	UserID                int `json:"user_id"`
	TotalSolicitudes      int `json:"total_solicitudes"`
	SolicitudesPendientes int `json:"solicitudes_pendientes"`
	EntregasProximas      int `json:"entregas_proximas"`
}

// ActivityLog — log de actividad (no existe en DDL, stub)
type ActivityLog struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
	Date   string `json:"date"`
	IP     string `json:"ip"`
	Device string `json:"device"`
}

// BodyMeasurements — medidas corporales.
// Mapea tabla "Medidas Funcionario". FuncionarioID se infiere del funcionario consultado.
// Activa se DERIVA: activa = (FechaFin == nil).
type BodyMeasurements struct {
	ID            int     `json:"id"`             // id_medidas
	FuncionarioID int     `json:"funcionario_id"` // inferido del funcionario consultado
	EstaturaM     float64 `json:"estatura_m"`     // estatura_m numeric(5,2)
	PechoCm       float64 `json:"pecho_cm"`       // pecho_cm numeric(5,2)
	CinturaCm     float64 `json:"cintura_cm"`     // cintura_cm numeric(5,2)
	CaderaCm      float64 `json:"cadera_cm"`      // cadera_cm numeric(5,2)
	MangaCm       float64 `json:"manga_cm"`       // manga_cm numeric(5,2)
	FechaInicio   string  `json:"fecha_inicio"`   // fecha_inicio date
	FechaFin      *string `json:"fecha_fin"`      // fecha_fin date (nullable)
	Activa        bool    `json:"activa"`         // derivado: FechaFin == nil
}
