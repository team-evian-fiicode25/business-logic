package services

import (
	"business_logic/internal/config"
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"text/template"

	"gopkg.in/gomail.v2"
)

type MailService interface {
	CreateMail(mailReq *Mail) []byte
	SendMail(mailReq *Mail) error
	NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail
}

type MailType int

const (
	MailConfirmation MailType = iota + 1
	PassReset
	TwoFactorAuth
)

type MailData struct {
	Username          string
	Code              string
	Token             string
	ExpiresInMin      int
	BaseURL           string
	VerifyEndpoint    string
	ResetPassEndpoint string
}

type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	mtype   MailType
	data    *MailData
}

type SMTPMailService struct {
	smtpHost string
	smtpPort int
	username string
	password string
}

func NewSMTPMailService(smtpHost string, smtpPort int, username string, password string) *SMTPMailService {
	return &SMTPMailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

func (ms *SMTPMailService) CreateMail(mailReq *Mail) (*gomail.Message, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", mailReq.from)
	m.SetHeader("To", mailReq.to...)
	m.SetHeader("Subject", mailReq.subject)

	body, err := parseHTMLTemplate(mailReq.mtype, mailReq.data)
	if err != nil {
		return nil, err
	}

	m.SetBody("text/html", body)
	return m, nil
}

func (ms *SMTPMailService) SendMail(mailReq *Mail) error {
	m, err := ms.CreateMail(mailReq)
	if err != nil {
		log.Println("Failed to create email:", err)
		return err
	}

	d := gomail.NewDialer(ms.smtpHost, ms.smtpPort, ms.username, ms.password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully to", mailReq.to)
	return nil
}

func (ms *SMTPMailService) NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail {
	secConfig := config.GetSecurityConfig()
	apiConfig := config.GetAPIConfig()

	switch mailType {
	case MailConfirmation:
		data.ExpiresInMin = int(secConfig.EmailExpiration.Minutes())
	case PassReset:
		data.ExpiresInMin = int(secConfig.PasswordExpiration.Minutes())
	case TwoFactorAuth:
		data.ExpiresInMin = int(secConfig.TwoFactorExpiration.Minutes())
	}

	data.BaseURL = apiConfig.BaseURL
	data.VerifyEndpoint = apiConfig.VerifyEndpoint
	data.ResetPassEndpoint = apiConfig.ResetPassEndpoint

	return &Mail{
		from:    from,
		to:      to,
		subject: subject,
		mtype:   mailType,
		data:    data,
	}
}

func parseHTMLTemplate(mailType MailType, data *MailData) (string, error) {
	templateConfig := config.GetTemplateConfig()
	var templateName string

	switch mailType {
	case MailConfirmation:
		templateName = templateConfig.MailConfirmation
	case PassReset:
		templateName = templateConfig.PasswordReset
	case TwoFactorAuth:
		templateName = templateConfig.TwoFactorAuth
	default:
		return "", fmt.Errorf("Invalid mail type")
	}

	templatePath := filepath.Join(templateConfig.BasePath, templateName+templateConfig.TemplateExtension)

	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Failed to parse template file %s: %v", templatePath, err)
		return "", err
	}

	var body bytes.Buffer
	if err := templ.Execute(&body, data); err != nil {
		log.Printf("Error executing template %s: %v", templatePath, err)
		return "", err
	}

	return body.String(), nil
}
