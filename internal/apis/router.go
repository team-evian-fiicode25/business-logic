package apis

import (
	"business_logic/internal/config"
	"business_logic/internal/handlers"
	"business_logic/internal/services"
	"net/http"
	"strings"
)

func SetupRoutes() {
	config.LoadEnv()

	mux := http.NewServeMux()

	consumerService := services.NewConsumerService("http://localhost:8080/graphql")

	mailConfig := config.GetMailConfig()

	smtpMailService := services.NewSMTPMailService(
		mailConfig.SMTPHost,
		mailConfig.SMTPPort,
		mailConfig.Username,
		mailConfig.Password,
	)

	consumerAuthHandler := handlers.NewConsumerAuthHandler(consumerService)
	verificationHandler := handlers.NewVerificationHandler(smtpMailService, consumerService)

	registerAuthRoutes(mux, "/api/consumer", consumerAuthHandler)

	registerMailRoutes(mux, verificationHandler)

	http.ListenAndServe(":8000", mux)
}

func registerAuthRoutes(mux *http.ServeMux, basePath string, handler handlers.AuthHandler) {
	mux.HandleFunc(basePath+"/auth/register", handler.Register)
	mux.HandleFunc(basePath+"/auth/login", handler.Login)
}

func registerMailRoutes(mux *http.ServeMux, handler *handlers.VerificationHandler) {
	mux.HandleFunc("/api/verify/mail", withMiddleware(
		handler.VerifyMail,
		handler.MiddlewareValidateVerificationData,
	))
	mux.HandleFunc("/api/verify/password-reset", withMiddleware(
		handler.VerifyPasswordReset,
		handler.MiddlewareValidateVerificationData,
	))

	mux.HandleFunc("/api/mail/send-verification", handler.SendVerificationEmail)

	mux.HandleFunc("/api/mail/get-password-reset-code", withMiddleware(
		handler.GeneratePassResetCode,
		handler.MiddlewareValidateAccessToken,
	))

	mux.HandleFunc("/api/user/reset-password", withMiddleware(
		handler.ResetPassword,
		handler.MiddlewareValidateAccessToken,
	))
}

func withMiddleware(h http.HandlerFunc, middleware func(http.Handler) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (strings.HasPrefix(r.URL.Path, "/api/verify/") && r.Method != http.MethodPost) ||
			(strings.HasPrefix(r.URL.Path, "/api/mail/get-password-reset-code") && r.Method != http.MethodGet) ||
			(strings.HasPrefix(r.URL.Path, "/api/user/reset-password") && r.Method != http.MethodPut) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler := middleware(http.HandlerFunc(h))
		handler.ServeHTTP(w, r)
	}
}
