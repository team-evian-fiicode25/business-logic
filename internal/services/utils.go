package services

import "regexp"

// TODO (mihaescuvlad): Use validator instead of this function
func isValidEmail(identifier string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(identifier)
}
