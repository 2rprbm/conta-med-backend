package domain

import (
	"regexp"
	"strings"
)

// Brazilian states mapping
var brazilianStates = map[string]string{
	"AC": "Acre",
	"AL": "Alagoas",
	"AP": "Amapá",
	"AM": "Amazonas",
	"BA": "Bahia",
	"CE": "Ceará",
	"DF": "Distrito Federal",
	"ES": "Espírito Santo",
	"GO": "Goiás",
	"MA": "Maranhão",
	"MT": "Mato Grosso",
	"MS": "Mato Grosso do Sul",
	"MG": "Minas Gerais",
	"PA": "Pará",
	"PB": "Paraíba",
	"PR": "Paraná",
	"PE": "Pernambuco",
	"PI": "Piauí",
	"RJ": "Rio de Janeiro",
	"RN": "Rio Grande do Norte",
	"RS": "Rio Grande do Sul",
	"RO": "Rondônia",
	"RR": "Roraima",
	"SC": "Santa Catarina",
	"SP": "São Paulo",
	"SE": "Sergipe",
	"TO": "Tocantins",
}

// ValidatePhoneNumber validates Brazilian phone number format
func ValidatePhoneNumber(phoneNumber string) bool {
	// Remove all non-digit characters
	cleanPhone := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")
	
	// Brazilian phone number patterns:
	// - Mobile: +55 11 9XXXX-XXXX (13 digits with country code)
	// - Landline: +55 11 XXXX-XXXX (12 digits with country code)
	// - Without country code: 11 9XXXX-XXXX (11 digits mobile) or 11 XXXX-XXXX (10 digits landline)
	
	if len(cleanPhone) == 13 && strings.HasPrefix(cleanPhone, "55") {
		// Mobile with country code +55
		return true
	}
	
	if len(cleanPhone) == 12 && strings.HasPrefix(cleanPhone, "55") {
		// Landline with country code +55
		return true
	}
	
	if len(cleanPhone) == 11 || len(cleanPhone) == 10 {
		// Without country code
		return true
	}
	
	return false
}

// ValidateMainMenuOption validates main menu selection
func ValidateMainMenuOption(option string) bool {
	normalizedOption := strings.TrimSpace(option)
	validOptions := []string{"1", "2", "3", "4"}
	
	for _, valid := range validOptions {
		if normalizedOption == valid {
			return true
		}
	}
	
	return false
}

// ValidateCRMOption validates CRM selection
func ValidateCRMOption(option string) bool {
	normalizedOption := strings.TrimSpace(option)
	validOptions := []string{"1", "2"}
	
	for _, valid := range validOptions {
		if normalizedOption == valid {
			return true
		}
	}
	
	return false
}

// ValidateBrazilianState validates Brazilian state abbreviation
func ValidateBrazilianState(state string) bool {
	normalizedState := strings.ToUpper(strings.TrimSpace(state))
	
	// Check if it's a valid state abbreviation
	if _, exists := brazilianStates[normalizedState]; exists {
		return true
	}
	
	// Check if it's a full state name
	normalizedStateLower := strings.ToLower(strings.TrimSpace(state))
	for _, stateName := range brazilianStates {
		if strings.ToLower(stateName) == normalizedStateLower {
			return true
		}
	}
	
	return false
}

// NormalizeBrazilianState normalizes Brazilian state to abbreviation format
func NormalizeBrazilianState(state string) string {
	normalizedState := strings.ToUpper(strings.TrimSpace(state))
	
	// If it's already an abbreviation
	if _, exists := brazilianStates[normalizedState]; exists {
		return normalizedState
	}
	
	// If it's a full state name, find the abbreviation
	normalizedStateLower := strings.ToLower(strings.TrimSpace(state))
	for abbr, stateName := range brazilianStates {
		if strings.ToLower(stateName) == normalizedStateLower {
			return abbr
		}
	}
	
	// If not found, return as-is
	return normalizedState
}

// ValidateCityName validates city name (basic validation)
func ValidateCityName(city string) bool {
	normalizedCity := strings.TrimSpace(city)
	
	// Basic validation: not empty and reasonable length
	if len(normalizedCity) < 2 || len(normalizedCity) > 100 {
		return false
	}
	
	// Check if contains only letters, spaces, apostrophes, and hyphens
	cityPattern := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'\-]+$`)
	return cityPattern.MatchString(normalizedCity)
}

// GetBrazilianStates returns all Brazilian states
func GetBrazilianStates() map[string]string {
	return brazilianStates
}

// GetMainMenuOptions returns valid main menu options with descriptions
func GetMainMenuOptions() map[string]string {
	return map[string]string{
		"1": "Já tenho uma empresa médica constituída",
		"2": "Quero abrir uma empresa",
		"3": "Gostaria de tirar dúvidas",
		"4": "Outros",
	}
}

// GetCRMOptions returns valid CRM options with descriptions
func GetCRMOptions() map[string]string {
	return map[string]string{
		"1": "Já tenho CRM",
		"2": "Ainda não possuo CRM",
	}
} 