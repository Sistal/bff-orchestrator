package models

// APIResponse es el envelope estándar de todas las respuestas del BFF,
// alineado con el contrato de ms-authentication.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorDetail error con código específico.
type ErrorDetail struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

// ValidationError error de validación de campo.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse respuesta de error simple (legacy).
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{Success: true, Message: message, Data: data}
}

func SuccessResponseWithMeta(message string, data interface{}, meta interface{}) APIResponse {
	return APIResponse{Success: true, Message: message, Data: data, Meta: meta}
}

func ErrorResponseWithCode(message string, code string, details string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Errors:  ErrorDetail{Code: code, Details: details},
	}
}

func ErrorResponseWithValidation(message string, errs []ValidationError) APIResponse {
	return APIResponse{Success: false, Message: message, Errors: errs}
}

func SimpleErrorResponse(message string) APIResponse {
	return APIResponse{Success: false, Message: message}
}
