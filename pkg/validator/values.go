package validator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/eduardo-works/writeit/pkg/project"
)

func ValidateNotEmpty(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", project.ErrEmptyValue
	}

	return value, nil
}

func ValidatePositive(value string) int {
	number, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || number <= 0 {
		return 0
	}

	return number
}

func ValidateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return "", fmt.Errorf("invalid email address: %s", email)
	}

	return email, nil
}

func ValidatePhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return phone, nil
	}

	phoneRegex := regexp.MustCompile(`^(\+\d{1,3})?(\d{9})$`)

	phone = strings.Join(strings.Fields(phone), "")

	match := phoneRegex.FindStringSubmatch(phone)
	if match == nil {
		return "", fmt.Errorf("invalid phone number: %s", phone)
	}

	number := match[2]
	formatted := fmt.Sprintf("%s %s %s", number[:3], number[3:6], number[6:])

	if match[1] != "" {
		return match[1] + " " + formatted, nil
	}

	return formatted, nil
}
