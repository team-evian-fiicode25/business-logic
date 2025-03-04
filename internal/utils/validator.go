package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type GenericResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Error   string                 `json:"error,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type ValidationErrorType string

const (
	Required      ValidationErrorType = "required"
	InvalidFormat ValidationErrorType = "invalid_format"
	Custom        ValidationErrorType = "custom"
)

type ValidationError struct {
	Field    string                 `json:"field"`
	Type     ValidationErrorType    `json:"type"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type ValidationResult struct {
	errors []ValidationError
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		errors: []ValidationError{},
	}
}

func (vr *ValidationResult) AddError(field string, errType ValidationErrorType, message string) {
	vr.errors = append(vr.errors, ValidationError{
		Field:   field,
		Type:    errType,
		Message: message,
	})
}

func (vr *ValidationResult) AddErrorWithMetadata(field string, errType ValidationErrorType, message string, metadata map[string]interface{}) {
	vr.errors = append(vr.errors, ValidationError{
		Field:    field,
		Type:     errType,
		Message:  message,
		Metadata: metadata,
	})
}

func (vr *ValidationResult) HasErrors() bool {
	return len(vr.errors) > 0
}

func (vr *ValidationResult) Errors() []ValidationError {
	return vr.errors
}

func (vr *ValidationResult) ErrorsMap() map[string][]ValidationError {
	result := make(map[string][]ValidationError)
	for _, err := range vr.errors {
		result[err.Field] = append(result[err.Field], err)
	}
	return result
}

func (vr *ValidationResult) ErrorsByType(errType ValidationErrorType) []ValidationError {
	var filtered []ValidationError
	for _, err := range vr.errors {
		if err.Type == errType {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

func (vr *ValidationResult) ErrorMessages() []string {
	messages := make([]string, len(vr.errors))
	for i, err := range vr.errors {
		messages[i] = err.Message
	}
	return messages
}

func (vr *ValidationResult) FirstError() *ValidationError {
	if len(vr.errors) > 0 {
		return &vr.errors[0]
	}
	return nil
}

func (vr *ValidationResult) ToAPIResponse() *GenericResponse {
	if !vr.HasErrors() {
		return &GenericResponse{
			Status:  true,
			Message: "Validation successful",
		}
	}

	firstError := vr.FirstError()
	return &GenericResponse{
		Status:  false,
		Message: "Validation failed",
		Error:   firstError.Message,
		Details: map[string]interface{}{
			"errors":    vr.Errors(),
			"errorsMap": vr.ErrorsMap(),
		},
	}
}

type Validator struct {
	emailRegex *regexp.Regexp
}

func NewValidator() *Validator {
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`),
	}
}

func (v *Validator) ValidateEmail(email string) (bool, *ValidationError) {
	if !v.emailRegex.MatchString(strings.ToLower(email)) {
		return false, &ValidationError{
			Field:   "email",
			Type:    InvalidFormat,
			Message: "Invalid email format",
			Metadata: map[string]interface{}{
				"pattern": v.emailRegex.String(),
				"value":   email,
			},
		}
	}
	return true, nil
}

func (v *Validator) ValidateRequired(value string, field string) (bool, *ValidationError) {
	if strings.TrimSpace(value) == "" {
		return false, &ValidationError{
			Field:   field,
			Type:    Required,
			Message: fmt.Sprintf("%s is required", field),
			Metadata: map[string]interface{}{
				"provided": value,
			},
		}
	}
	return true, nil
}

func (v *Validator) Validate(i interface{}) *ValidationResult {
	result := NewValidationResult()

	switch obj := i.(type) {
	case VerificationDataValidator:
		if ok, err := v.ValidateRequired(obj.GetIdentifier(), "identifier"); !ok {
			result.AddErrorWithMetadata(
				err.Field,
				err.Type,
				err.Message,
				err.Metadata,
			)
		}

		if ok, err := v.ValidateRequired(obj.GetVerificationType(), "verification_type"); !ok {
			result.AddErrorWithMetadata(
				err.Field,
				err.Type,
				err.Message,
				err.Metadata,
			)
		}

		if strings.Contains(strings.ToLower(obj.GetVerificationType()), "email") {
			if ok, err := v.ValidateEmail(obj.GetIdentifier()); !ok {
				result.AddErrorWithMetadata(
					"identifier",
					err.Type,
					err.Message,
					err.Metadata,
				)
			}
		}

		if obj.RequiresCode() && obj.GetCode() == "" {
			result.AddError(
				"code",
				Required,
				"Verification code is required",
			)
		}

		if obj.RequiresToken() && obj.GetToken() == "" {
			result.AddError(
				"token",
				Required,
				"Verification token is required",
			)
		}
	}

	return result
}

type VerificationDataValidator interface {
	GetIdentifier() string
	GetVerificationType() string
	GetCode() string
	RequiresCode() bool
	GetToken() string
	RequiresToken() bool
}
