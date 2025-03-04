package config

import (
	"os"
	"strconv"
)

type MailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	FromAddr string
}

func GetMailConfig() MailConfig {
	config := MailConfig{
		SMTPHost: getEnvStr("MAIL_SMTP_HOST", ""),
		SMTPPort: 587,
		Username: getEnvStr("MAIL_USERNAME", ""),
		Password: getEnvStr("MAIL_PASSWORD", ""),
		FromAddr: getEnvStr("MAIL_FROM_ADDRESS", ""),
	}

	if portStr := os.Getenv("MAIL_SMTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.SMTPPort = port
		}
	}

	return config
}
