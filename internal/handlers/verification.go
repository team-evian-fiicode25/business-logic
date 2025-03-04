package handlers

import (
	"business_logic/internal/config"
	"business_logic/internal/models"
	"business_logic/internal/services"
	"business_logic/internal/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ContextKey string

const (
	VerificationDataKey ContextKey = "verification_data"
	AccessTokenKey      ContextKey = "access_token"
)

type VerificationHandler struct {
	mailService     *services.SMTPMailService
	consumerService *services.ConsumerService
	fromAddress     string
	validator       *utils.Validator
	logger          *log.Logger
}

func NewVerificationHandler(mailService *services.SMTPMailService, consumerService *services.ConsumerService) *VerificationHandler {
	mailConfig := config.GetMailConfig()

	// TODO (mihaescuvlad): Discuss logging
	logger := log.New(log.Writer(), "[VERIFICATION] ", log.LstdFlags)

	return &VerificationHandler{
		mailService:     mailService,
		consumerService: consumerService,
		fromAddress:     mailConfig.FromAddr,
		validator:       utils.NewValidator(),
		logger:          logger,
	}
}

func (h *VerificationHandler) MiddlewareValidateVerificationData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.logger.Print("Validating verification data")

		verificationReq := &models.VerificationRequest{}

		if err := utils.FromJSON(verificationReq, r.Body); err != nil {
			h.logger.Printf("Deserialization of verification data failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJSON(&models.GenericResponse{
				Status:  false,
				Message: "Invalid request format",
				Error:   err.Error(),
			}, w)
			return
		}

		validationResult := h.validator.Validate(verificationReq)
		if validationResult.HasErrors() {
			h.logger.Printf("Validation of verification data failed: %v", validationResult.ErrorMessages())
			w.WriteHeader(http.StatusBadRequest)

			apiResponse := validationResult.ToAPIResponse()
			utils.ToJSON(apiResponse, w)
			return
		}

		ctx := context.WithValue(r.Context(), VerificationDataKey, *verificationReq)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *VerificationHandler) MiddlewareValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.logger.Print("Validating access token")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.logger.Print("Authorization header is missing")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ToJSON(&models.GenericResponse{
				Status:  false,
				Message: "Authentication required",
				Error:   "Authorization header is missing",
			}, w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.logger.Print("Invalid Authorization header format")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ToJSON(&models.GenericResponse{
				Status:  false,
				Message: "Invalid authentication format",
				Error:   "Authorization header format must be Bearer <token>",
			}, w)
			return
		}

		token := parts[1]

		// TODO (mihaescuvlad): Call C# Auth to validate token
		if token == "" {
			h.logger.Print("Empty token provided")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ToJSON(&models.GenericResponse{
				Status:  false,
				Message: "Invalid authentication",
				Error:   "Empty token provided",
			}, w)
			return
		}

		ctx := context.WithValue(r.Context(), AccessTokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetVerificationDataFromContext(ctx context.Context) (models.VerificationRequest, bool) {
	verificationData, ok := ctx.Value(VerificationDataKey).(models.VerificationRequest)
	return verificationData, ok
}

func GetAccessTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AccessTokenKey).(string)
	return token, ok
}

func (h *VerificationHandler) VerifyMail(w http.ResponseWriter, r *http.Request) {
	verificationReq, ok := GetVerificationDataFromContext(r.Context())
	if !ok {
		h.logger.Print("Verification data not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Internal server error",
			Error:   "Verification data not found in context",
		}, w)
		return
	}

	h.logger.Printf("Processing verification for: %s, type: %s",
		verificationReq.Identifier, verificationReq.VerificationType)

	switch verificationReq.VerificationType {
	case models.EmailVerification:
		if verificationReq.Token != "" {
			h.logger.Printf("Processing email verification with token: %s...",
				verificationReq.Token[:8])
			// TODO (mihaescuvlad): Call C# Auth to verify email token
		} else if verificationReq.Code != "" {
			h.logger.Printf("Processing email verification with manual code: %s", verificationReq.Code)
			// TODO (mihaescuvlad): Verify 2FA code (probably in C#)
		} else {
			h.logger.Print("Neither token nor code provided for email verification")
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJSON(&models.GenericResponse{
				Status:  false,
				Message: "Validation failed",
				Error:   "Either token or code must be provided",
			}, w)
			return
		}

	case models.TwoFactorEmail:
		if verificationReq.Code == "" {
			h.logger.Print("No code provided for 2FA verification")
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJSON(&models.GenericResponse{
				Status:  false,
				Message: "Validation failed",
				Error:   "Code is required for 2FA",
			}, w)
			return
		}

		h.logger.Printf("Processing 2FA verification with code: %s", verificationReq.Code)
		// TODO (mihaescuvlad): Verify the 2FA code

	default:
		h.logger.Printf("Invalid verification type: %s", verificationReq.VerificationType)
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Invalid verification type",
			Error:   "Unsupported verification type",
		}, w)
		return
	}

	response := models.VerificationResponse{
		Success: true,
		Message: "Verification successful",
	}

	utils.ToJSON(response, w)
}

func (h *VerificationHandler) VerifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	verificationReq, ok := GetVerificationDataFromContext(r.Context())
	if !ok {
		h.logger.Print("Verification data not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Internal server error",
			Error:   "Verification data not found in context",
		}, w)
		return
	}

	h.logger.Printf("Processing password reset verification for: %s", verificationReq.Identifier)

	// TODO (mihaescuvlad): Check password reset code

	response := models.VerificationResponse{
		Success: true,
		Message: "Password reset code verified",
	}

	utils.ToJSON(response, w)
}

func (h *VerificationHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.logger.Print("Processing verification email request")

	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	if err := utils.FromJSON(&req, r.Body); err != nil {
		h.logger.Printf("Failed to parse request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Invalid request format",
			Error:   err.Error(),
		}, w)
		return
	}

	if valid, err := h.validator.ValidateEmail(req.Email); !valid {
		h.logger.Printf("Invalid email format: %s, error: %v", req.Email, err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Validation failed",
			Error:   "Invalid email format",
		}, w)
		return
	}

	secConfig := config.GetSecurityConfig()

	token, err := generateBase64URLToken(secConfig.TokenLengthBytes)
	if err != nil {
		h.logger.Printf("Failed to generate verification token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Failed to generate verification token",
			Error:   err.Error(),
		}, w)
		return
	}

	shortCode := generateNumericCode(secConfig.ShortCodeLength)

	h.logger.Printf("Sending verification email to %s with token: %s...",
		req.Email, token[:8])

	// TODO (mihaescuvlad): Store token

	data := &services.MailData{
		Username: req.Username,
		Code:     shortCode,
		Token:    token,
	}

	mail := h.mailService.NewMail(
		h.fromAddress,
		[]string{req.Email},
		"Email Verification",
		services.MailConfirmation,
		data,
	)

	if err := h.mailService.SendMail(mail); err != nil {
		h.logger.Printf("Failed to send verification email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Failed to send email",
			Error:   err.Error(),
		}, w)
		return
	}

	response := models.VerificationResponse{
		Success:      true,
		Message:      "Verification email sent",
		ExpiresInSec: int(secConfig.EmailExpiration.Seconds()),
	}

	utils.ToJSON(response, w)
}

func generateBase64URLToken(lengthInBytes int) (string, error) {
	token := make([]byte, lengthInBytes)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	base64Token := base64.StdEncoding.EncodeToString(token)
	base64Token = strings.TrimRight(base64Token, "=")
	base64Token = strings.Replace(base64Token, "+", "-", -1)
	base64Token = strings.Replace(base64Token, "/", "_", -1)

	return base64Token, nil
}

func generateNumericCode(length int) string {
	const digits = "0123456789"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}

	for i, b := range bytes {
		bytes[i] = digits[b%byte(len(digits))]
	}
	return string(bytes)
}

func (h *VerificationHandler) GeneratePassResetCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.logger.Print("Processing password reset code request")

	token, ok := GetAccessTokenFromContext(r.Context())
	if !ok {
		h.logger.Print("Access token not found in context, this should not happen because we have middleware")
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Internal server error",
			Error:   "Access token not found in context",
		}, w)
		return
	}

	h.logger.Printf("User with token %s requested password reset", token[:8])

	email := r.URL.Query().Get("email")
	if email == "" {
		h.logger.Print("Email parameter missing")
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Missing parameter",
			Error:   "Email is required",
		}, w)
		return
	}

	if valid, err := h.validator.ValidateEmail(email); !valid {
		h.logger.Printf("Invalid email format: %s, error: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Validation failed",
			Error:   "Invalid email format",
		}, w)
		return
	}

	// TODO (mihaescuvlad): I think this needs to be moved to C#. If not moved to C#
	// use apiConfig := config.GetAPIConfig() and generate the URL
	// verificationURL := apiConfig.BaseURL + apiConfig.VerifyEndpoint + "?token=" + token
	secConfig := config.GetSecurityConfig()

	code := generateNumericCode(secConfig.ShortCodeLength)
	resetToken, err := generateBase64URLToken(secConfig.TokenLengthBytes)
	if err != nil {
		h.logger.Printf("Failed to generate reset token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Failed to generate reset token",
			Error:   err.Error(),
		}, w)
		return
	}

	verificationData := &models.VerificationData{
		Token:            resetToken,
		Code:             code,
		ExpiresAt:        time.Now().Add(secConfig.PasswordExpiration),
		CreatedAt:        time.Now(),
		VerificationType: models.PasswordReset,
		UserID:           "",
		IsVerified:       false,
		Attempts:         0,
	}

	// TODO (mihaescuvlad): Store the verification data with expiration time
	_ = verificationData

	h.logger.Printf("Generated verification data for %s: code=%s, token=%s...",
		email, code, resetToken[:8])

	data := &services.MailData{
		Username: "",
		Code:     code,
		Token:    resetToken,
	}

	mail := h.mailService.NewMail(
		h.fromAddress,
		[]string{email},
		"Password Reset Request",
		services.PassReset,
		data,
	)

	if err := h.mailService.SendMail(mail); err != nil {
		h.logger.Printf("Failed to send password reset email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Failed to send email",
			Error:   err.Error(),
		}, w)
		return
	}

	response := models.VerificationResponse{
		Success:      true,
		Message:      "Password reset code sent",
		ExpiresInSec: int(secConfig.PasswordExpiration.Seconds()),
	}

	utils.ToJSON(response, w)
}

func (h *VerificationHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.logger.Print("Processing password reset request")

	token, ok := GetAccessTokenFromContext(r.Context())
	if !ok {
		h.logger.Print("Access token not found in context, this should not happen because we have middleware")
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Internal server error",
			Error:   "Access token not found in context",
		}, w)
		return
	}

	h.logger.Printf("User with token %s is resetting password", token[:8])

	var req struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}

	if err := utils.FromJSON(&req, r.Body); err != nil {
		h.logger.Printf("Failed to parse request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Invalid request format",
			Error:   err.Error(),
		}, w)
		return
	}

	if valid, err := h.validator.ValidateEmail(req.Email); !valid {
		h.logger.Printf("Invalid email format: %s, error: %v", req.Email, err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Validation failed",
			Error:   "Invalid email format",
		}, w)
		return
	}

	if req.Code == "" || req.NewPassword == "" {
		h.logger.Print("Missing required fields")
		w.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&models.GenericResponse{
			Status:  false,
			Message: "Missing required fields",
			Error:   "Code and new password are required",
		}, w)
		return
	}

	h.logger.Printf("Password reset successful for %s", req.Email)
	response := models.GenericResponse{
		Status:  true,
		Message: "Password has been reset successfully",
	}

	utils.ToJSON(response, w)
}
