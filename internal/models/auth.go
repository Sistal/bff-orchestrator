package models

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type AuthValidateResponse struct {
	Valid     bool   `json:"valid"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiresAt int64  `json:"expires_at"`
}

type AuthMeResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
