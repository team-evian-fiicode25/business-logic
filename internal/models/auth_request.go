package models

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthResponse struct {
	Session_id string `json:"session_id,omitempty"`
	Error      string `json:"error,omitempty"`
}
