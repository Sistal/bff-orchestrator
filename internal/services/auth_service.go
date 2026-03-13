package services

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/models"
	"go.uber.org/zap"
)

// AuthService define las operaciones expuestas por el BFF de autenticación.
type AuthService interface {
	// Públicos
	Login(nombreUsuario, password string) (int, *models.APIResponse, error)
	Register(req *models.RegisterRequest) (int, *models.APIResponse, error)
	Validate(token string) (int, *models.APIResponse, error)
	Refresh(refreshToken string) (int, *models.APIResponse, error)

	// Protegidos (requieren JWT)
	GetMe(token string) (int, *models.APIResponse, error)
	Logout(token string) (int, *models.APIResponse, error)
	ChangePassword(token string, req *models.ChangePasswordRequest) (int, *models.APIResponse, error)
	GetRoles(token string, activosSolo *bool) (int, *models.APIResponse, error)

	// Administración
	CreateUser(token string, req *models.RegisterRequest) (int, *models.APIResponse, error)
	ListUsers(token string, params url.Values) (int, *models.APIResponse, error)
	GetUserByID(token string, id int) (int, *models.APIResponse, error)
	UpdateUser(token string, id int, req *models.UpdateUsuarioRequest) (int, *models.APIResponse, error)
}

// HTTPAuthService delega todas las operaciones al ms-authentication.
type HTTPAuthService struct {
	client          *clients.IdentityClient
	employeeService EmployeeService
}

func NewHTTPAuthService(client *clients.IdentityClient, empSvc EmployeeService) AuthService {
	return &HTTPAuthService{
		client:          client,
		employeeService: empSvc,
	}
}

func (s *HTTPAuthService) Login(nombreUsuario, password string) (int, *models.APIResponse, error) {
	return s.client.Login(map[string]string{
		"nombre_usuario": nombreUsuario,
		"password":       password,
	})
}

func (s *HTTPAuthService) Register(req *models.RegisterRequest) (int, *models.APIResponse, error) {
	log := logger.Get()
	log.Info("Register: Iniciando registro de usuario",
		zap.String("username", req.NombreUsuario),
		zap.String("rut", req.RUT),
		zap.String("rut_funcionario", req.RutFuncionario),
	)

	// 1. Crear usuario en ms-authentication
	code, resp, err := s.client.Register(req)
	if err != nil {
		log.Error("Register: Error al llamar a ms-authentication", zap.Error(err))
		return code, resp, err
	}

	// Si la creación del usuario falló, retornamos el error tal cual
	if code != http.StatusCreated || resp == nil || !resp.Success {
		log.Warn("Register: Creación de usuario fallida en ms-authentication",
			zap.Int("status_code", code),
			zap.Any("response", resp),
		)
		return code, resp, nil
	}

	log.Info("Register: Usuario creado exitosamente en ms-authentication",
		zap.Int("status_code", code),
	)

	// 2. Si hay datos de funcionario, intentar crearlo en ms-funcionario
	// Verificamos si vienen campos clave del funcionario
	if req.RutFuncionario != "" && req.Email != "" {
		log.Info("Register: Detectados datos de funcionario, procediendo a crear perfil asociado")

		// Extraer ID del usuario recién creado desde la respuesta
		// Se asume que resp.Data es un map o struct que tiene "id_usuario"
		// La estructura esperada es models.UsuarioResponseDTO
		var createdUserID int

		// Intentar hacer marshal/unmarshal para obtener el ID de forma segura
		dataBytes, _ := json.Marshal(resp.Data)
		var userData models.UsuarioResponseDTO
		if err := json.Unmarshal(dataBytes, &userData); err == nil && userData.IDUsuario != 0 {
			createdUserID = userData.IDUsuario
		} else {
			// Fallback: intentar leer como map[string]interface{} si viene genérico
			var userMap map[string]interface{}
			if err := json.Unmarshal(dataBytes, &userMap); err == nil {
				if idFloat, ok := userMap["id_usuario"].(float64); ok {
					createdUserID = int(idFloat)
				}
			}
		}

		if createdUserID != 0 {
			log.Info("Register: ID de usuario extraído exitosamente", zap.Int("created_user_id", createdUserID))

			// Preparar request para ms-funcionario
			empReq := models.CreateEmployeeRequestFromLogin{
				UserID:          createdUserID,
				RutFuncionario:  req.RutFuncionario,
				Nombres:         req.Nombres,
				ApellidoPaterno: req.ApellidoPaterno,
				ApellidoMaterno: req.ApellidoMaterno,
				Email:           req.Email,
				Genero:          req.Genero,
			}

			log.Info("Register: Llamando a employeeService.CreateEmployeeFromLogin",
				zap.Int("user_id", createdUserID),
				zap.String("rut_funcionario", empReq.RutFuncionario),
			)

			// Llamar al servicio de empleados
			_, errEmp := s.employeeService.CreateEmployeeFromLogin(empReq)
			if errEmp != nil {
				// LOGUEAR ERROR pero NO fallar el request completo, ya que el usuario ya fue creado.
				// Podríamos retornar un warning en el mensaje o meta.
				log.Error("Register: Usuario creado pero falló creación de funcionario",
					zap.Int("id_usuario", createdUserID),
					zap.String("error", errEmp.Error()),
				)
				// Opcional: Agregar warning al response
				resp.Message += " (Advertencia: Perfil de funcionario no pudo ser creado)"
			} else {
				log.Info("Register: Funcionario creado exitosamente vinculado al usuario",
					zap.Int("id_usuario", createdUserID),
				)
			}
		} else {
			log.Warn("Register: No se pudo extraer ID del usuario creado, omitiendo creación de funcionario")
		}
	} else {
		log.Info("Register: No se detectaron datos de funcionario, finalizando registro solo como usuario")
	}

	return code, resp, nil
}

func (s *HTTPAuthService) Validate(token string) (int, *models.APIResponse, error) {
	return s.client.ValidateTokenRaw(token)
}

func (s *HTTPAuthService) Refresh(refreshToken string) (int, *models.APIResponse, error) {
	return s.client.Refresh(map[string]string{
		"refresh_token": refreshToken,
	})
}

func (s *HTTPAuthService) GetMe(token string) (int, *models.APIResponse, error) {
	return s.client.GetMe(token)
}

func (s *HTTPAuthService) Logout(token string) (int, *models.APIResponse, error) {
	return s.client.Logout(token)
}

func (s *HTTPAuthService) ChangePassword(token string, req *models.ChangePasswordRequest) (int, *models.APIResponse, error) {
	return s.client.ChangePassword(token, req)
}

func (s *HTTPAuthService) GetRoles(token string, activosSolo *bool) (int, *models.APIResponse, error) {
	return s.client.GetRoles(token, activosSolo)
}

func (s *HTTPAuthService) CreateUser(token string, req *models.RegisterRequest) (int, *models.APIResponse, error) {
	return s.client.CreateUser(token, req)
}

func (s *HTTPAuthService) ListUsers(token string, params url.Values) (int, *models.APIResponse, error) {
	return s.client.ListUsers(token, params)
}

func (s *HTTPAuthService) GetUserByID(token string, id int) (int, *models.APIResponse, error) {
	return s.client.GetUserByID(token, id)
}

func (s *HTTPAuthService) UpdateUser(token string, id int, req *models.UpdateUsuarioRequest) (int, *models.APIResponse, error) {
	return s.client.UpdateUser(token, id, req)
}
