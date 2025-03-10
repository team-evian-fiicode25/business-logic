package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type SecurityConfig struct {
	TokenLengthBytes    int           `json:"token_length_bytes"`
	ShortCodeLength     int           `json:"short_code_length"`
	EmailExpiration     time.Duration `json:"email_expiration"`
	PasswordExpiration  time.Duration `json:"password_expiration"`
	TwoFactorExpiration time.Duration `json:"two_factor_expiration"`
}

type TemplateConfig struct {
	BasePath          string `json:"base_path"`
	MailConfirmation  string `json:"mail_confirmation"`
	PasswordReset     string `json:"password_reset"`
	TwoFactorAuth     string `json:"two_factor_auth"`
	TemplateExtension string `json:"template_extension"`
}

type APIConfig struct {
	BaseURL           string `json:"base_url"`
	VerifyEndpoint    string `json:"verify_endpoint"`
	ResetPassEndpoint string `json:"reset_pass_endpoint"`
}

type AppConfig struct {
	Security  SecurityConfig `json:"security"`
	Templates TemplateConfig `json:"templates"`
	API       APIConfig      `json:"api"`
}

func init() {
	LoadEnv()
	defaultConfig = GetDefaultConfig()
}

func LoadEnv() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("Warning: Could not determine working directory:", err)
		return
	}

	envPath := filepath.Join(workDir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		envPath = filepath.Join(filepath.Dir(workDir), ".env")
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			log.Println("Warning: .env file not found, using default or existing environment variables")
			return
		}
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}
}

func GetDefaultConfig() AppConfig {
	return AppConfig{
		Security: SecurityConfig{
			TokenLengthBytes:    getEnvInt("SECURITY_TOKEN_LENGTH", 128),
			ShortCodeLength:     getEnvInt("SECURITY_CODE_LENGTH", 6),
			EmailExpiration:     time.Duration(getEnvInt("SECURITY_EMAIL_EXPIRATION", 1800)) * time.Second,   // 30 min
			PasswordExpiration:  time.Duration(getEnvInt("SECURITY_PASSWORD_EXPIRATION", 900)) * time.Second, // 15 min
			TwoFactorExpiration: time.Duration(getEnvInt("SECURITY_2FA_EXPIRATION", 300)) * time.Second,      // 5 min
		},
		Templates: TemplateConfig{
			BasePath:          getEnvStr("TEMPLATE_BASE_PATH", "templates"),
			MailConfirmation:  "mail_confirmation",
			PasswordReset:     "pass_reset",
			TwoFactorAuth:     "2fa_email",
			TemplateExtension: getEnvStr("TEMPLATE_EXTENSION", ".html"),
		},
		API: APIConfig{
			BaseURL:           getEnvStr("API_BASE_URL", ""),
			VerifyEndpoint:    getEnvStr("API_VERIFY_ENDPOINT", "/verify"),
			ResetPassEndpoint: getEnvStr("API_RESET_PASS_ENDPOINT", "/reset-password"),
		},
	}
}

func getEnvStr(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

var defaultConfig = GetDefaultConfig()

func GetConfig() AppConfig {
	return defaultConfig
}

func GetSecurityConfig() SecurityConfig {
	return defaultConfig.Security
}

func GetTemplateConfig() TemplateConfig {
	return defaultConfig.Templates
}

func GetAPIConfig() APIConfig {
	return defaultConfig.API
}
