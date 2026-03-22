package util

import (
	"regexp"
	"strings"
)

var (
	phonePattern = regexp.MustCompile(`^01[016789]-?\d{3,4}-?\d{4}$`)
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func IsValidPhone(phone string) bool {
	return phonePattern.MatchString(strings.TrimSpace(phone))
}

func IsValidEmail(email string) bool {
	return emailPattern.MatchString(strings.TrimSpace(email))
}
