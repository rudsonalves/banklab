package domain

import (
	"regexp"
	"strings"
)

var cpfRegex = regexp.MustCompile(`^\d{11}$`)

// ValidateCPF validates a CPF string according to Brazilian rules.
// Returns true if the CPF is valid, false otherwise.
//
// A valid CPF must:
// 1. Have exactly 11 digits
// 2. Not be a known invalid CPF (all same digits)
// 3. Have correct check digits (first and second verification digits)
func ValidateCPF(cpf string) bool {
	// Remove common formatting characters
	normalized := normalizeCPF(cpf)

	// Check format: must be exactly 11 digits
	if !cpfRegex.MatchString(normalized) {
		return false
	}

	// CPFs with all same digits are invalid
	if isAllSameDigits(normalized) {
		return false
	}

	// Validate first check digit
	firstDigit := calculateCheckDigit(normalized[:9], 10)
	if normalized[9:10] != string(rune('0'+firstDigit)) {
		return false
	}

	// Validate second check digit
	secondDigit := calculateCheckDigit(normalized[:10], 11)
	if normalized[10:11] != string(rune('0'+secondDigit)) {
		return false
	}

	return true
}

// normalizeCPF removes common formatting characters from a CPF string.
func normalizeCPF(cpf string) string {
	// Remove spaces, dots, hyphens, and slashes
	normalized := strings.NewReplacer(
		" ", "",
		".", "",
		"-", "",
		"/", "",
	).Replace(cpf)
	return strings.TrimSpace(normalized)
}

// isAllSameDigits returns true if all digits in the CPF are the same.
func isAllSameDigits(cpf string) bool {
	if len(cpf) == 0 {
		return false
	}
	firstChar := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != firstChar {
			return false
		}
	}
	return true
}

// calculateCheckDigit calculates the check digit for a CPF.
// length should be 10 for the first digit or 11 for the second digit.
func calculateCheckDigit(digits string, length int) int {
	sum := 0
	for i, ch := range digits {
		digit := int(ch - '0')
		sum += digit * (length - i)
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}
