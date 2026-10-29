package validator

import (
	"errors"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestValidateNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"normal string", "Hello", "Hello", false},
		{"leading spaces", "  Hello", "Hello", false},
		{"trailing spaces", "Hello  ", "Hello", false},
		{"surrounding spaces", "  Hello  ", "Hello", false},
		{"internal spaces", "Hello World", "Hello World", false},
		{"empty string", "", "", true},
		{"spaces only", "   ", "", true},
		{"tabs and newlines", "\t\n ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateNotEmpty(tt.input)

			if tt.wantErr {
				if !errors.Is(err, project.ErrEmptyValue) {
					t.Fatalf("expected ErrEmptyValue, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"positive integer", "10", 10},
		{"one", "1", 1},
		{"leading spaces", " 12", 12},
		{"trailing spaces", "12 ", 12},
		{"leading zeros", "0012", 12},
		{"explicit plus", "+12", 12},
		{"zero", "0", 0},
		{"negative integer", "-5", 0},
		{"decimal number", "1.1", 0},
		{"text", "abc", 0},
		{"empty string", "", 0},
		{"spaces only", "   ", 0},
		{"overflow", "999999999999999999999999999999", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePositive(tt.input)

			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple email", "user@example.com", "user@example.com", false},
		{"email with dot", "john.doe@example.com", "john.doe@example.com", false},
		{"email with plus", "user+tag@example.com", "user+tag@example.com", false},
		{"email with subdomain", "user@mail.example.com", "user@mail.example.com", false},
		{"surrounding spaces", "  user@example.com  ", "user@example.com", false},
		{"missing at", "userexample.com", "", true},
		{"missing domain", "user@", "", true},
		{"missing local part", "@example.com", "", true},
		{"empty string", "", "", true},
		{"spaces only", "   ", "", true},
		{"display name", "John <user@example.com>", "", true},
		{"multiple addresses", "a@example.com, b@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateEmail(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"national number", "222333444", "222 333 444", false},
		{"national formatted", "222 333 444", "222 333 444", false},
		{"international compact", "+351222333444", "+351 222 333 444", false},
		{"international separated", "+351 222333444", "+351 222 333 444", false},
		{"international formatted", "+351 222 333 444", "+351 222 333 444", false},
		{"one digit country code", "+1222333444", "+1 222 333 444", false},
		{"two digit country code", "+44222333444", "+44 222 333 444", false},
		{"surrounding spaces", "  +351 222 333 444  ", "+351 222 333 444", false},
		{"tabs", "+351\t222\t333\t444", "+351 222 333 444", false},
		{"empty string", "", "", false},
		{"spaces only", "   ", "", false},
		{"too short", "12345", "", true},
		{"too long", "1234567890", "", true},
		{"missing plus", "351222333444", "", true},
		{"letters", "abcdefghi", "", true},
		{"hyphens", "222-333-444", "", true},
		{"parentheses", "(222)333444", "", true},
		{"country code too long", "+1234222333444", "", true},
		{"country code without number", "+351", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePhone(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
