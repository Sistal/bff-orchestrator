package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/models"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// ─── Públicos ─────────────────────────────────────────────────────────────────

// Login autentica un usuario y retorna JWT + datos de usuario.
// @Summary      Login
// @Description  Autentica un usuario y retorna un JWT de acceso
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.LoginRequest true "Credenciales"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Failure      403  {object}  models.APIResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	log := logger.Get()
	ip := c.ClientIP()

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Login: body inválido",
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.ErrorResponseWithValidation(
			"Datos de entrada inválidos",
			[]models.ValidationError{{Field: "request", Message: err.Error()}},
		))
		return
	}

	log.Info("Login: intento de autenticación",
		zap.String("nombre_usuario", req.NombreUsuario),
		zap.String("ip", ip),
	)

	code, resp, err := h.service.Login(req.NombreUsuario, req.Password)
	if err != nil {
		log.Error("Login: error al contactar ms-authentication",
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	if code == http.StatusOK {
		log.Info("Login: autenticación exitosa",
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Int("status", code),
		)

		// Extraer token y expires_in del data de la respuesta para almacenarlos en cookies.
		if resp != nil && resp.Data != nil {
			dataBytes, _ := json.Marshal(resp.Data)
			var loginData models.LoginResponseData
			if err := json.Unmarshal(dataBytes, &loginData); err == nil && loginData.Token != "" {
				maxAge := loginData.ExpiresIn
				if maxAge <= 0 {
					maxAge = cookieMaxAge()
				}
				domain := cookieDomain()
				secureCookie := cookieSecure()

				log.Info("Login: estableciendo cookie access_token",
					zap.String("nombre_usuario", req.NombreUsuario),
					zap.String("ip", ip),
					zap.String("token_prefix", safeTokenPrefix(loginData.Token)),
					zap.Int("max_age_seconds", maxAge),
					zap.String("cookie_domain", domain),
					zap.String("cookie_name", "access_token"),
					zap.Bool("cookie_secure", secureCookie),
				)

				// SameSiteNoneMode es necesario para requests cross-origin (ej. frontend s-dev vs api-s-dev subdomains) con credentials
				c.SetSameSite(http.SameSiteNoneMode)
				c.SetCookie("access_token", loginData.Token, maxAge, "/", domain, secureCookie, true)

				// ── DIAGNÓSTICO: confirmar Set-Cookie header enviado al browser ──
				log.Debug("Login: Set-Cookie header enviado al browser",
					zap.String("set_cookie_header", c.Writer.Header().Get("Set-Cookie")),
				)
				// ── FIN DIAGNÓSTICO ─────────────────────────────────────────────

				log.Debug("Login: cookie access_token seteada correctamente",
					zap.String("nombre_usuario", req.NombreUsuario),
					zap.String("ip", ip),
					zap.String("token_prefix", safeTokenPrefix(loginData.Token)),
				)
			} else if err != nil {
				log.Error("Login: error al deserializar LoginResponseData para setear cookie",
					zap.String("nombre_usuario", req.NombreUsuario),
					zap.String("ip", ip),
					zap.Error(err),
				)
			} else {
				log.Warn("Login: token vacío en la respuesta del ms-authentication, cookie no seteada",
					zap.String("nombre_usuario", req.NombreUsuario),
					zap.String("ip", ip),
				)
			}
		}
	} else {
		log.Warn("Login: autenticación fallida",
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// Register registra un nuevo usuario público.
// @Summary      Register
// @Description  Registra un nuevo usuario de forma pública
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.RegisterRequest true "Datos del usuario"
// @Success      201  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log := logger.Get()
	ip := c.ClientIP()

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Register: body inválido",
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.ErrorResponseWithValidation(
			"Datos de entrada inválidos",
			[]models.ValidationError{{Field: "request", Message: err.Error()}},
		))
		return
	}

	log.Info("Register: solicitud de registro de usuario",
		zap.String("nombre_usuario", req.NombreUsuario),
		zap.String("nombre_completo", req.NombreCompleto),
		zap.String("rut", req.RUT),
		zap.Int("id_rol", req.IDRol),
		zap.String("ip", ip),
		zap.String("nombres", req.Nombres),
		zap.String("apellido_paterno", req.ApellidoPaterno),
		zap.String("apellido_materno", req.ApellidoMaterno),
		zap.String("rut_funcionario", req.RutFuncionario),
	)

	code, resp, err := h.service.Register(&req)
	if err != nil {
		log.Error("Register: error al contactar ms-authentication",
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	if code == http.StatusCreated {
		log.Info("Register: usuario creado exitosamente",
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Int("status", code),
		)
	} else {
		log.Warn("Register: registro fallido",
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// Validate valida un JWT y retorna el payload de claims.
// @Summary      Validate token
// @Description  Valida un JWT y retorna el payload de claims
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Router       /api/v1/auth/validate [get]
func (h *AuthHandler) Validate(c *gin.Context) {
	log := logger.Get()
	ip := c.ClientIP()

	log.Info("Validate: solicitud de validación de token recibida",
		zap.String("ip", ip),
	)

	// ── DIAGNÓSTICO: loguear headers y cookies crudas recibidas ──────────────
	allCookies := c.Request.Cookies()
	cookieNames := make([]string, 0, len(allCookies))
	for _, ck := range allCookies {
		cookieNames = append(cookieNames, ck.Name)
	}
	log.Debug("Validate: headers y cookies recibidos en la request",
		zap.String("host", c.Request.Host),
		zap.String("origin", c.Request.Header.Get("Origin")),
		zap.String("authorization_header", func() string {
			if h := c.Request.Header.Get("Authorization"); h != "" {
				return safeTokenPrefix(strings.TrimPrefix(h, "Bearer "))
			}
			return "(ausente)"
		}()),
		zap.Strings("cookie_names", cookieNames),
		zap.String("cookie_header_raw", c.Request.Header.Get("Cookie")),
	)
	// ── FIN DIAGNÓSTICO ──────────────────────────────────────────────────────

	// Intentar obtener el token desde el header Authorization.
	var token string
	var tokenSource string

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
		tokenSource = "header"
		log.Debug("Validate: token obtenido desde header Authorization",
			zap.String("ip", ip),
			zap.String("token_prefix", safeTokenPrefix(token)),
		)
	} else if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
		// Fallback: intentar desde la cookie access_token.
		token = cookieToken
		tokenSource = "cookie"
		log.Debug("Validate: token obtenido desde cookie access_token (fallback)",
			zap.String("ip", ip),
			zap.String("token_prefix", safeTokenPrefix(token)),
		)
	} else {
		log.Warn("Validate: token ausente — ni header Authorization ni cookie access_token presentes",
			zap.String("ip", ip),
		)
		c.JSON(http.StatusUnauthorized, models.ErrorResponseWithCode(
			"Token inválido o expirado",
			"INVALID_TOKEN",
			"El token proporcionado no es válido",
		))
		return
	}

	log.Info("Validate: enviando token a ms-authentication para validación",
		zap.String("ip", ip),
		zap.String("token_source", tokenSource),
		zap.String("token_prefix", safeTokenPrefix(token)),
	)

	code, resp, err := h.service.Validate(token)
	if err != nil {
		log.Error("Validate: error al contactar ms-authentication",
			zap.String("ip", ip),
			zap.String("token_source", tokenSource),
			zap.String("token_prefix", safeTokenPrefix(token)),
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, models.ErrorResponseWithCode(
			"Token inválido o expirado",
			"INVALID_TOKEN",
			"El token proporcionado no es válido",
		))
		return
	}

	if code == http.StatusOK {
		log.Info("Validate: token válido — ms-authentication confirmó la validez",
			zap.String("ip", ip),
			zap.String("token_source", tokenSource),
			zap.String("token_prefix", safeTokenPrefix(token)),
			zap.Int("status", code),
		)
	} else {
		log.Warn("Validate: token rechazado por ms-authentication",
			zap.String("ip", ip),
			zap.String("token_source", tokenSource),
			zap.String("token_prefix", safeTokenPrefix(token)),
			zap.Int("status", code),
		)
	}

	log.Debug("Validate: respondiendo al cliente",
		zap.String("ip", ip),
		zap.Int("status", code),
	)

	c.JSON(code, resp)
}

// Refresh genera un nuevo JWT a partir de un refresh token.
// @Summary      Refresh token
// @Description  Genera un nuevo JWT de acceso a partir de un refresh token válido
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.RefreshTokenRequest true "Refresh token"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	log := logger.Get()
	ip := c.ClientIP()

	// Leer el refresh token desde la cookie HttpOnly.
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		log.Warn("Refresh: cookie refresh_token ausente",
			zap.String("ip", ip),
		)
		c.JSON(http.StatusUnauthorized, models.ErrorResponseWithCode(
			"Sesión expirada",
			"MISSING_REFRESH_TOKEN",
			"No se encontró el refresh token. Por favor inicia sesión nuevamente.",
		))
		return
	}

	log.Info("Refresh: solicitud de renovación de token",
		zap.String("ip", ip),
		zap.String("refresh_token_prefix", safeTokenPrefix(refreshToken)),
	)

	code, resp, err := h.service.Refresh(refreshToken)
	if err != nil {
		log.Error("Refresh: error al contactar ms-authentication",
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	if code == http.StatusOK {
		log.Info("Refresh: token renovado exitosamente",
			zap.String("ip", ip),
			zap.Int("status", code),
		)

		// Sobrescribir cookies con los nuevos tokens.
		if resp != nil && resp.Data != nil {
			dataBytes, _ := json.Marshal(resp.Data)
			var tokenData models.TokenResponseData
			if err := json.Unmarshal(dataBytes, &tokenData); err == nil {
				domain := cookieDomain()
				secureCookie := cookieSecure()
				if tokenData.Token != "" {
					maxAge := tokenData.ExpiresIn
					if maxAge <= 0 {
						maxAge = cookieMaxAge()
					}

					log.Info("Refresh: renovando cookie access_token",
						zap.String("ip", ip),
						zap.String("token_prefix", safeTokenPrefix(tokenData.Token)),
						zap.Int("max_age_seconds", maxAge),
						zap.String("cookie_domain", domain),
						zap.String("cookie_name", "access_token"),
						zap.Bool("cookie_secure", secureCookie),
					)

					c.SetSameSite(http.SameSiteNoneMode)
					c.SetCookie("access_token", tokenData.Token, maxAge, "/", domain, secureCookie, true)

					log.Debug("Refresh: cookie access_token renovada correctamente",
						zap.String("ip", ip),
						zap.String("token_prefix", safeTokenPrefix(tokenData.Token)),
					)
				}
				if tokenData.RefreshToken != nil && *tokenData.RefreshToken != "" {
					refreshMaxAge := cookieMaxAge() * 24

					log.Info("Refresh: renovando cookie refresh_token",
						zap.String("ip", ip),
						zap.String("refresh_token_prefix", safeTokenPrefix(*tokenData.RefreshToken)),
						zap.Int("max_age_seconds", refreshMaxAge),
						zap.String("cookie_domain", domain),
						zap.String("cookie_name", "refresh_token"),
						zap.Bool("cookie_secure", secureCookie),
					)

					c.SetSameSite(http.SameSiteNoneMode)
					c.SetCookie("refresh_token", *tokenData.RefreshToken, refreshMaxAge, "/", domain, secureCookie, true)

					log.Debug("Refresh: cookie refresh_token renovada correctamente",
						zap.String("ip", ip),
						zap.String("refresh_token_prefix", safeTokenPrefix(*tokenData.RefreshToken)),
					)
				}
			} else {
				log.Error("Refresh: error al deserializar TokenResponseData para renovar cookies",
					zap.String("ip", ip),
					zap.Error(err),
				)
			}
		}
	} else {
		log.Warn("Refresh: renovación fallida",
			zap.String("ip", ip),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// ─── Protegidos ───────────────────────────────────────────────────────────────

// GetMe retorna el perfil completo del usuario autenticado.
// @Summary      Get current user
// @Description  Retorna el perfil del usuario autenticado
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Status(c *gin.Context) {
	log := logger.Get()
	userIDStr, exists := c.Get("userID") // El middleware BearerAuthMiddleware lo guarda como "userID"
	if !exists {
		log.Warn("Status: Usuario no autenticado (falta userID en contexto)")
		c.JSON(http.StatusUnauthorized, models.SimpleErrorResponse("Usuario no autenticado"))
		return
	}

	userID, ok := userIDStr.(string)
	if !ok || userID == "" {
		log.Warn("Status: userID inválido o vacío en contexto", zap.Any("userID_raw", userIDStr))
		c.JSON(http.StatusUnauthorized, models.SimpleErrorResponse("Token inválido o sin información de usuario"))
		return
	}

	code, resp, err := h.service.Status(userID)
	if err != nil {
		log.Error("Status: error al verificar estado del perfil", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	c.JSON(code, resp)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	log := logger.Get()
	userID := c.GetString("userID")
	ip := c.ClientIP()

	log.Info("GetMe: solicitud de perfil",
		zap.String("user_id", userID),
		zap.String("ip", ip),
	)

	token := extractToken(c)
	code, resp, err := h.service.GetMe(token)
	if err != nil {
		log.Error("GetMe: error al obtener perfil del usuario",
			zap.String("user_id", userID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error al obtener información del usuario"))
		return
	}

	if code != http.StatusOK {
		log.Warn("GetMe: respuesta no exitosa del microservicio",
			zap.String("user_id", userID),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// Logout invalida la sesión del usuario.
// @Summary      Logout
// @Description  Invalida la sesión del usuario autenticado
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.APIResponse
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	log := logger.Get()
	userID := c.GetString("userID")
	ip := c.ClientIP()

	log.Info("Logout: cierre de sesión solicitado",
		zap.String("user_id", userID),
		zap.String("ip", ip),
	)

	token := extractToken(c)
	code, resp, err := h.service.Logout(token)
	if err != nil {
		log.Error("Logout: error al cerrar sesión",
			zap.String("user_id", userID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	// Expirar ambas cookies independientemente del resultado del microservicio.
	domain := cookieDomain()
	secureCookie := cookieSecure()

	log.Info("Logout: expirando cookies de sesión",
		zap.String("user_id", userID),
		zap.String("ip", ip),
		zap.Strings("cookies", []string{"access_token", "refresh_token"}),
		zap.String("cookie_domain", domain),
		zap.Bool("cookie_secure", secureCookie),
	)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("access_token", "", -1, "/", domain, secureCookie, true)
	c.SetCookie("refresh_token", "", -1, "/", domain, secureCookie, true)

	log.Debug("Logout: cookies access_token y refresh_token expiradas correctamente",
		zap.String("user_id", userID),
		zap.String("ip", ip),
	)

	log.Info("Logout: sesión cerrada",
		zap.String("user_id", userID),
		zap.String("ip", ip),
		zap.Int("status", code),
	)

	c.JSON(code, resp)
}

// ChangePassword cambia la contraseña del usuario autenticado.
// @Summary      Change password
// @Description  Cambia la contraseña del usuario autenticado
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body models.ChangePasswordRequest true "Contraseñas"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Router       /api/v1/auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	log := logger.Get()
	userID := c.GetString("userID")
	ip := c.ClientIP()

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("ChangePassword: body inválido",
			zap.String("user_id", userID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.ErrorResponseWithValidation(
			"Error en la validación",
			[]models.ValidationError{{Field: "request", Message: err.Error()}},
		))
		return
	}

	log.Info("ChangePassword: solicitud de cambio de contraseña",
		zap.String("user_id", userID),
		zap.String("ip", ip),
	)

	token := extractToken(c)
	code, resp, err := h.service.ChangePassword(token, &req)
	if err != nil {
		log.Error("ChangePassword: error al contactar ms-authentication",
			zap.String("user_id", userID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	if code == http.StatusOK {
		log.Info("ChangePassword: contraseña actualizada exitosamente",
			zap.String("user_id", userID),
			zap.String("ip", ip),
		)
	} else {
		log.Warn("ChangePassword: cambio de contraseña fallido",
			zap.String("user_id", userID),
			zap.String("ip", ip),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// GetRoles lista los roles disponibles en el sistema.
// @Summary      Get roles
// @Description  Lista los roles disponibles en el sistema
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Param        activos_solo query bool false "Filtrar solo roles activos (default: true)"
// @Success      200  {object}  models.APIResponse
// @Failure      500  {object}  models.APIResponse
// @Router       /api/v1/auth/roles [get]
func (h *AuthHandler) GetRoles(c *gin.Context) {
	log := logger.Get()
	userID := c.GetString("userID")
	ip := c.ClientIP()

	var activosSolo *bool
	if val, exists := c.GetQuery("activos_solo"); exists {
		b, err := strconv.ParseBool(val)
		if err == nil {
			activosSolo = &b
		} else {
			log.Warn("GetRoles: valor inválido para activos_solo",
				zap.String("user_id", userID),
				zap.String("valor", val),
			)
		}
	}

	log.Info("GetRoles: listando roles",
		zap.String("user_id", userID),
		zap.String("ip", ip),
		zap.Any("activos_solo", activosSolo),
	)

	token := extractToken(c)
	code, resp, err := h.service.GetRoles(token, activosSolo)
	if err != nil {
		log.Error("GetRoles: error al obtener roles",
			zap.String("user_id", userID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error al obtener roles"))
		return
	}

	c.JSON(code, resp)
}

// ─── Administración ───────────────────────────────────────────────────────────

// CreateUser crea un usuario (admin).
// @Summary      Create user (admin)
// @Description  Crea un nuevo usuario con control total de atributos
// @Tags         auth-admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body models.RegisterRequest true "Datos del usuario"
// @Success      201  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      403  {object}  models.APIResponse
// @Router       /api/v1/auth/users [post]
func (h *AuthHandler) CreateUser(c *gin.Context) {
	log := logger.Get()
	adminID := c.GetString("userID")
	ip := c.ClientIP()

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("CreateUser: body inválido",
			zap.String("admin_id", adminID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.ErrorResponseWithValidation(
			"Datos de entrada inválidos",
			[]models.ValidationError{{Field: "request", Message: err.Error()}},
		))
		return
	}

	log.Info("CreateUser: creación de usuario por admin",
		zap.String("admin_id", adminID),
		zap.String("nombre_usuario", req.NombreUsuario),
		zap.String("nombre_completo", req.NombreCompleto),
		zap.String("rut", req.RUT),
		zap.Int("id_rol", req.IDRol),
		zap.String("ip", ip),
	)

	token := extractToken(c)
	code, resp, err := h.service.CreateUser(token, &req)
	if err != nil {
		log.Error("CreateUser: error al contactar ms-authentication",
			zap.String("admin_id", adminID),
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	if code == http.StatusCreated {
		log.Info("CreateUser: usuario creado exitosamente",
			zap.String("admin_id", adminID),
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.Int("status", code),
		)
	} else {
		log.Warn("CreateUser: creación de usuario fallida",
			zap.String("admin_id", adminID),
			zap.String("nombre_usuario", req.NombreUsuario),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// ListUsers lista todos los usuarios con paginación y filtros.
// @Summary      List users (admin)
// @Description  Lista todos los usuarios con paginación y filtros opcionales
// @Tags         auth-admin
// @Produce      json
// @Security     BearerAuth
// @Param        page     query int    false "Número de página (default: 1)"
// @Param        limit    query int    false "Resultados por página (default: 20, max: 100)"
// @Param        id_rol   query int    false "Filtrar por ID de rol"
// @Param        id_estado query int   false "Filtrar por ID de estado"
// @Param        search   query string false "Búsqueda en nombre_usuario, nombre_completo, rut"
// @Param        sort_by  query string false "Campo para ordenar (default: fecha_creacion)"
// @Param        order    query string false "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  models.APIResponse
// @Failure      500  {object}  models.APIResponse
// @Router       /api/v1/auth/users [get]
func (h *AuthHandler) ListUsers(c *gin.Context) {
	log := logger.Get()
	adminID := c.GetString("userID")
	ip := c.ClientIP()

	params := url.Values{}
	for key, vals := range c.Request.URL.Query() {
		for _, v := range vals {
			params.Add(key, v)
		}
	}

	log.Info("ListUsers: listado de usuarios por admin",
		zap.String("admin_id", adminID),
		zap.String("ip", ip),
		zap.String("query", params.Encode()),
	)

	token := extractToken(c)
	code, resp, err := h.service.ListUsers(token, params)
	if err != nil {
		log.Error("ListUsers: error al listar usuarios",
			zap.String("admin_id", adminID),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error al listar usuarios"))
		return
	}

	if code != http.StatusOK {
		log.Warn("ListUsers: respuesta no exitosa del microservicio",
			zap.String("admin_id", adminID),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// GetUserByID obtiene un usuario por su ID.
// @Summary      Get user by ID (admin)
// @Description  Obtiene información detallada de un usuario por su ID
// @Tags         auth-admin
// @Produce      json
// @Security     BearerAuth
// @Param        id_usuario path int true "ID del usuario"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      404  {object}  models.APIResponse
// @Router       /api/v1/auth/users/{id_usuario} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	log := logger.Get()
	adminID := c.GetString("userID")
	ip := c.ClientIP()

	id, err := strconv.Atoi(c.Param("id_usuario"))
	if err != nil {
		log.Warn("GetUserByID: id_usuario inválido",
			zap.String("admin_id", adminID),
			zap.String("raw_id", c.Param("id_usuario")),
			zap.String("ip", ip),
		)
		c.JSON(http.StatusBadRequest, models.SimpleErrorResponse("ID de usuario inválido"))
		return
	}

	log.Info("GetUserByID: consulta de usuario por admin",
		zap.String("admin_id", adminID),
		zap.Int("target_user_id", id),
		zap.String("ip", ip),
	)

	token := extractToken(c)
	code, resp, errSvc := h.service.GetUserByID(token, id)
	if errSvc != nil {
		log.Error("GetUserByID: error al obtener usuario",
			zap.String("admin_id", adminID),
			zap.Int("target_user_id", id),
			zap.String("ip", ip),
			zap.Error(errSvc),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error interno del servidor"))
		return
	}

	if code == http.StatusNotFound {
		log.Warn("GetUserByID: usuario no encontrado",
			zap.String("admin_id", adminID),
			zap.Int("target_user_id", id),
		)
	}

	c.JSON(code, resp)
}

// UpdateUser actualiza parcialmente los campos de un usuario.
// @Summary      Update user (admin)
// @Description  Actualiza parcialmente los campos de un usuario
// @Tags         auth-admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id_usuario path int true "ID del usuario"
// @Param        body body models.UpdateUsuarioRequest true "Campos a actualizar"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      404  {object}  models.APIResponse
// @Router       /api/v1/auth/users/{id_usuario} [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	log := logger.Get()
	adminID := c.GetString("userID")
	ip := c.ClientIP()

	id, err := strconv.Atoi(c.Param("id_usuario"))
	if err != nil {
		log.Warn("UpdateUser: id_usuario inválido",
			zap.String("admin_id", adminID),
			zap.String("raw_id", c.Param("id_usuario")),
			zap.String("ip", ip),
		)
		c.JSON(http.StatusBadRequest, models.SimpleErrorResponse("ID de usuario inválido"))
		return
	}

	var req models.UpdateUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("UpdateUser: body inválido",
			zap.String("admin_id", adminID),
			zap.Int("target_user_id", id),
			zap.String("ip", ip),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, models.SimpleErrorResponse("Datos de actualización inválidos"))
		return
	}

	log.Info("UpdateUser: actualización de usuario por admin",
		zap.String("admin_id", adminID),
		zap.Int("target_user_id", id),
		zap.String("ip", ip),
	)

	token := extractToken(c)
	code, resp, errSvc := h.service.UpdateUser(token, id, &req)
	if errSvc != nil {
		log.Error("UpdateUser: error al actualizar usuario",
			zap.String("admin_id", adminID),
			zap.Int("target_user_id", id),
			zap.String("ip", ip),
			zap.Error(errSvc),
		)
		c.JSON(http.StatusInternalServerError, models.SimpleErrorResponse("Error al actualizar usuario"))
		return
	}

	if code == http.StatusOK {
		log.Info("UpdateUser: usuario actualizado exitosamente",
			zap.String("admin_id", adminID),
			zap.Int("target_user_id", id),
		)
	} else {
		log.Warn("UpdateUser: actualización fallida",
			zap.String("admin_id", adminID),
			zap.Int("target_user_id", id),
			zap.Int("status", code),
		)
	}

	c.JSON(code, resp)
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// extractToken obtiene el JWT desde el header Authorization o, como fallback,
// desde la cookie HttpOnly "access_token" establecida en el login.
func extractToken(c *gin.Context) string {
	if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		return cookie
	}
	return ""
}

// cookieMaxAge retorna el valor de COOKIE_MAX_AGE (en segundos) o 3600 por defecto.
func cookieMaxAge() int {
	if v := os.Getenv("COOKIE_MAX_AGE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 3600
}

// cookieDomain retorna el dominio a usar en SetCookie.
// En desarrollo local con COOKIE_DOMAIN vacío, se retorna "" para que el browser
// aplique la cookie al host exacto (localhost) sin fijar el atributo Domain.
// Fijar Domain=localhost explícitamente puede causar rechazo en Chrome/Firefox
// porque localhost es un TLD especial según RFC 6265.
// En producción usa el valor de COOKIE_DOMAIN tal cual.
func cookieDomain() string {
	if d := os.Getenv("COOKIE_DOMAIN"); d != "" {
		return d
	}
	return ""
}

// cookieSecure retorna true si se debe configurar la cookie como Secure.
// Debería ser true en entornos productivos con HTTPS o si se fuerza con env var.
func cookieSecure() bool {
	if os.Getenv("ENVIRONMENT") == "production" || os.Getenv("COOKIE_SECURE") == "true" {
		return true
	}
	// Permite false en entorno de desarrollo local sin HTTPS
	return false
}

// safeTokenPrefix retorna solo los primeros 10 caracteres del token para logs (evita exposición).
func safeTokenPrefix(token string) string {
	if len(token) > 10 {
		return token[:10] + "..."
	}
	return "***"
}
