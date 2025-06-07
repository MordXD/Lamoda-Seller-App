package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a validation error with field and message
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		return v[0].Error()
	}
	
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// Add adds a new validation error
func (v *ValidationErrors) Add(field, message string) {
	*v = append(*v, ValidationError{Field: field, Message: message})
}

// PasswordValidator contains password validation configuration
type PasswordValidator struct {
	MinLength        int
	MaxLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireDigit     bool
	RequireSpecial   bool
	ForbiddenWords   []string
}

// DefaultPasswordValidator returns a password validator with secure defaults
func DefaultPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:        8,
		MaxLength:        128,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireDigit:     true,
		RequireSpecial:   true,
		ForbiddenWords: []string{
			"password", "123456", "qwerty", "admin", "user", "login",
			"lamoda", "seller", "app", "test", "demo",
		},
	}
}

// ValidatePassword validates a password against the configured rules
func (pv *PasswordValidator) ValidatePassword(password string) ValidationErrors {
	var errors ValidationErrors

	// Check length
	if len(password) < pv.MinLength {
		errors.Add("password", fmt.Sprintf("Password must be at least %d characters long", pv.MinLength))
	}
	if len(password) > pv.MaxLength {
		errors.Add("password", fmt.Sprintf("Password must not exceed %d characters", pv.MaxLength))
	}

	// Check character requirements
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if pv.RequireUppercase && !hasUpper {
		errors.Add("password", "Password must contain at least one uppercase letter")
	}
	if pv.RequireLowercase && !hasLower {
		errors.Add("password", "Password must contain at least one lowercase letter")
	}
	if pv.RequireDigit && !hasDigit {
		errors.Add("password", "Password must contain at least one digit")
	}
	if pv.RequireSpecial && !hasSpecial {
		errors.Add("password", "Password must contain at least one special character")
	}

	// Check for forbidden words
	passwordLower := strings.ToLower(password)
	for _, word := range pv.ForbiddenWords {
		if strings.Contains(passwordLower, strings.ToLower(word)) {
			errors.Add("password", fmt.Sprintf("Password cannot contain the word '%s'", word))
			break
		}
	}

	// Check for common patterns
	if isCommonPattern(password) {
		errors.Add("password", "Password cannot be a common pattern (e.g., 12345678, abcdefgh)")
	}

	return errors
}

// EmailValidator contains email validation configuration
type EmailValidator struct {
	MaxLength       int
	AllowedDomains  []string
	BlockedDomains  []string
	RequireMX       bool
}

// DefaultEmailValidator returns an email validator with sensible defaults
func DefaultEmailValidator() *EmailValidator {
	return &EmailValidator{
		MaxLength: 255,
		BlockedDomains: []string{
			"tempmail.org", "10minutemail.com", "guerrillamail.com",
			"mailinator.com", "yopmail.com", "temp-mail.org",
		},
		RequireMX: false, // Set to true in production if you want MX record validation
	}
}

// ValidateEmail validates an email address
func (ev *EmailValidator) ValidateEmail(email string) ValidationErrors {
	var errors ValidationErrors

	// Basic validation
	if email == "" {
		errors.Add("email", "Email is required")
		return errors
	}

	// Length check
	if len(email) > ev.MaxLength {
		errors.Add("email", fmt.Sprintf("Email must not exceed %d characters", ev.MaxLength))
	}

	// Parse email
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		errors.Add("email", "Invalid email format")
		return errors
	}

	email = parsedEmail.Address

	// Additional format validation
	if !isValidEmailFormat(email) {
		errors.Add("email", "Invalid email format")
		return errors
	}

	// Extract domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		errors.Add("email", "Invalid email format")
		return errors
	}
	domain := strings.ToLower(parts[1])

	// Check allowed domains (if specified)
	if len(ev.AllowedDomains) > 0 {
		allowed := false
		for _, allowedDomain := range ev.AllowedDomains {
			if domain == strings.ToLower(allowedDomain) {
				allowed = true
				break
			}
		}
		if !allowed {
			errors.Add("email", "Email domain is not allowed")
		}
	}

	// Check blocked domains
	for _, blockedDomain := range ev.BlockedDomains {
		if domain == strings.ToLower(blockedDomain) {
			errors.Add("email", "Email domain is not allowed")
			break
		}
	}

	return errors
}

// UserValidator combines password and email validation
type UserValidator struct {
	PasswordValidator *PasswordValidator
	EmailValidator    *EmailValidator
}

// NewUserValidator creates a new user validator with default settings
func NewUserValidator() *UserValidator {
	return &UserValidator{
		PasswordValidator: DefaultPasswordValidator(),
		EmailValidator:    DefaultEmailValidator(),
	}
}

// ValidateRegistration validates user registration data
func (uv *UserValidator) ValidateRegistration(name, email, password string) ValidationErrors {
	var errors ValidationErrors

	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		errors.Add("name", "Name is required")
	} else if len(name) < 2 {
		errors.Add("name", "Name must be at least 2 characters long")
	} else if len(name) > 100 {
		errors.Add("name", "Name must not exceed 100 characters")
	} else if !isValidName(name) {
		errors.Add("name", "Name contains invalid characters")
	}

	// Validate email
	emailErrors := uv.EmailValidator.ValidateEmail(email)
	errors = append(errors, emailErrors...)

	// Validate password
	passwordErrors := uv.PasswordValidator.ValidatePassword(password)
	errors = append(errors, passwordErrors...)

	return errors
}

// ValidateProfileUpdate validates profile update data
func (uv *UserValidator) ValidateProfileUpdate(name, email string) ValidationErrors {
	var errors ValidationErrors

	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		errors.Add("name", "Name is required")
	} else if len(name) < 2 {
		errors.Add("name", "Name must be at least 2 characters long")
	} else if len(name) > 100 {
		errors.Add("name", "Name must not exceed 100 characters")
	} else if !isValidName(name) {
		errors.Add("name", "Name contains invalid characters")
	}

	// Validate email
	emailErrors := uv.EmailValidator.ValidateEmail(email)
	errors = append(errors, emailErrors...)

	return errors
}

// Helper functions

func isValidEmailFormat(email string) bool {
	// More strict email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidName(name string) bool {
	// Allow letters, spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-'\.]+$`)
	return nameRegex.MatchString(name)
}

func isCommonPattern(password string) bool {
	commonPatterns := []string{
		"12345678", "87654321", "abcdefgh", "qwertyui",
		"asdfghjk", "zxcvbnm", "password", "Password",
	}
	
	for _, pattern := range commonPatterns {
		if strings.Contains(strings.ToLower(password), strings.ToLower(pattern)) {
			return true
		}
	}
	
	// Check for sequential characters
	if hasSequentialChars(password, 4) {
		return true
	}
	
	// Check for repeated characters
	if hasRepeatedChars(password, 4) {
		return true
	}
	
	return false
}

func hasSequentialChars(s string, minLength int) bool {
	if len(s) < minLength {
		return false
	}
	
	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1]+1 || s[i] == s[i-1]-1 {
			count++
			if count >= minLength {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

func hasRepeatedChars(s string, minLength int) bool {
	if len(s) < minLength {
		return false
	}
	
	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= minLength {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}