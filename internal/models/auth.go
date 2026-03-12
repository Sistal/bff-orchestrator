package models

// ─── DTOs compartidos ────────────────────────────────────────────────────────

// RolDTO representa un rol del sistema.
type RolDTO struct {
	IDRol       int    `json:"id_rol"`
	NombreRol   string `json:"nombre_rol"`
	Descripcion string `json:"descripcion"`
}

// EstadoDTO representa el estado de una entidad.
type EstadoDTO struct {
	IDEstado     int    `json:"id_estado"`
	NombreEstado string `json:"nombre_estado"`
	TablaEstado  string `json:"tabla_estado"`
}

// UsuarioResponseDTO datos completos de usuario expuestos en las respuestas.
type UsuarioResponseDTO struct {
	IDUsuario         int       `json:"id_usuario"`
	NombreUsuario     string    `json:"nombre_usuario"`
	NombreCompleto    string    `json:"nombre_completo"`
	RUT               string    `json:"rut"`
	Rol               RolDTO    `json:"rol"`
	Estado            EstadoDTO `json:"estado"`
	FechaCreacion     string    `json:"fecha_creacion"`
	FechaModificacion *string   `json:"fecha_modificacion"`
}

// PaginationMeta metadatos de paginación.
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ─── Request bodies ───────────────────────────────────────────────────────────

// LoginRequest body de POST /api/v1/auth/login
type LoginRequest struct {
	NombreUsuario string `json:"nombre_usuario" binding:"required"`
	Password      string `json:"password"       binding:"required"`
}

// RegisterRequest body de POST /api/v1/auth/register
type RegisterRequest struct {
	NombreUsuario   string `json:"nombre_usuario"     binding:"required"`
	NombreCompleto  string `json:"nombre_completo"    binding:"required"`
	RUT             string `json:"rut"                binding:"required"`
	Password        string `json:"password"           binding:"required"`
	IDRol           int    `json:"id_rol"             binding:"required"`
	IDEstadoUsuario int    `json:"id_estado_usuario"`
}

// RefreshTokenRequest body de POST /api/v1/auth/refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest body de PUT /api/v1/auth/change-password
type ChangePasswordRequest struct {
	PasswordActual       string `json:"password_actual"       binding:"required"`
	PasswordNueva        string `json:"password_nueva"        binding:"required"`
	PasswordConfirmacion string `json:"password_confirmacion" binding:"required"`
}

// UpdateUsuarioRequest body de PUT /api/v1/auth/users/:id
type UpdateUsuarioRequest struct {
	NombreCompleto  string `json:"nombre_completo"`
	IDRol           int    `json:"id_rol"`
	IDEstadoUsuario int    `json:"id_estado_usuario"`
}

// ─── Response bodies ──────────────────────────────────────────────────────────

// LoginResponseData payload de data en POST /api/v1/auth/login
type LoginResponseData struct {
	Usuario   UsuarioResponseDTO `json:"usuario"`
	Token     string             `json:"token"`
	ExpiresIn int                `json:"expires_in"`
	TokenType string             `json:"token_type"`
}

// TokenResponseData payload de data en POST /api/v1/auth/refresh
type TokenResponseData struct {
	Token        string  `json:"token"`
	ExpiresIn    int     `json:"expires_in"`
	TokenType    string  `json:"token_type"`
	RefreshToken *string `json:"refresh_token"`
}

// ValidateTokenData payload de data en GET /api/v1/auth/validate
type ValidateTokenData struct {
	IDUsuario       int    `json:"id_usuario"`
	NombreUsuario   string `json:"nombre_usuario"`
	NombreCompleto  string `json:"nombre_completo"`
	RUT             string `json:"rut"`
	IDRol           int    `json:"id_rol"`
	NombreRol       string `json:"nombre_rol"`
	IDEstadoUsuario int    `json:"id_estado_usuario"`
	NombreEstado    string `json:"nombre_estado"`
	Exp             int64  `json:"exp"`
	Iat             int64  `json:"iat"`
}

// ─── Tipos legacy (usados internamente por el middleware) ─────────────────────

// AuthValidateResponse respuesta interna del ms-authentication para validación de token.
type AuthValidateResponse struct {
	Valid     bool   `json:"valid"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	RoleName  string `json:"role_name"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiresAt int64  `json:"expires_at"`
}

// AuthMeResponse respuesta legacy de /auth/me (reemplazado por UsuarioResponseDTO).
type AuthMeResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Rol   string `json:"rol"`
}
