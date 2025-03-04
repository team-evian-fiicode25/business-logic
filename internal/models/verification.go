package models

import (
	"time"
)

type VerificationType string

const (
	EmailVerification VerificationType = "email_verification"
	TwoFactorEmail    VerificationType = "2fa_email"
	TwoFactorPhone    VerificationType = "2fa_phone"
	PasswordReset     VerificationType = "password_reset"
)

type VerificationData struct {
	Token            string           `json:"token"`
	Code             string           `json:"code"`
	ExpiresAt        time.Time        `json:"expires_at"`
	CreatedAt        time.Time        `json:"created_at"`
	VerificationType VerificationType `json:"verification_type"`
	UserID           string           `json:"user_id"`
	IsVerified       bool             `json:"is_verified"`
	Attempts         int              `json:"attempts,omitempty"`
}

type VerificationRequest struct {
	VerificationType VerificationType `json:"verification_type"`
	Identifier       string           `json:"identifier"`
	Code             string           `json:"code,omitempty"`
	Token            string           `json:"token,omitempty"`
}

func (v *VerificationRequest) GetIdentifier() string {
	return v.Identifier
}

func (v *VerificationRequest) GetVerificationType() string {
	return string(v.VerificationType)
}

func (v *VerificationRequest) GetCode() string {
	return v.Code
}

func (v *VerificationRequest) GetToken() string {
	return v.Token
}

func (v *VerificationRequest) RequiresToken() bool {
	switch v.VerificationType {
	case EmailVerification:
		return v.Code == ""
	case PasswordReset:
		return true
	default:
		return false
	}
}

func (v *VerificationRequest) RequiresCode() bool {
	switch v.VerificationType {
	case TwoFactorEmail, TwoFactorPhone:
		return true
	case EmailVerification:
		return v.Token == ""
	case PasswordReset:
		return true
	default:
		return false
	}
}

type VerificationResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
	ExpiresInSec int    `json:"expires_in_sec,omitempty"`
}

type GenericResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}
