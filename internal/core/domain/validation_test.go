package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePhoneNumber(t *testing.T) {
	t.Run("should return true when phone number has valid format with country code", func(t *testing.T) {
		// arrange
		phoneNumbers := []string{
			"+5511999999999",
			"5511999999999",
			"+551199999999",
			"551199999999",
		}

		// act & assert
		for _, phone := range phoneNumbers {
			assert.True(t, ValidatePhoneNumber(phone), "Phone number %s should be valid", phone)
		}
	})

	t.Run("should return true when phone number has valid format without country code", func(t *testing.T) {
		// arrange
		phoneNumbers := []string{
			"11999999999",
			"11999999999",
		}

		// act & assert
		for _, phone := range phoneNumbers {
			assert.True(t, ValidatePhoneNumber(phone), "Phone number %s should be valid", phone)
		}
	})

	t.Run("should return false when phone number has invalid format", func(t *testing.T) {
		// arrange
		phoneNumbers := []string{
			"123",
			"12345678901234",
			"abc123456789",
			"",
			"55119999999999", // too many digits
		}

		// act & assert
		for _, phone := range phoneNumbers {
			assert.False(t, ValidatePhoneNumber(phone), "Phone number %s should be invalid", phone)
		}
	})
}

func TestValidateMainMenuOption(t *testing.T) {
	t.Run("should return true when option is valid", func(t *testing.T) {
		// arrange
		validOptions := []string{"1", "2", "3", "4", " 1 ", " 2 "}

		// act & assert
		for _, option := range validOptions {
			assert.True(t, ValidateMainMenuOption(option), "Option %s should be valid", option)
		}
	})

	t.Run("should return false when option is invalid", func(t *testing.T) {
		// arrange
		invalidOptions := []string{"0", "5", "a", "", "12", "1a"}

		// act & assert
		for _, option := range invalidOptions {
			assert.False(t, ValidateMainMenuOption(option), "Option %s should be invalid", option)
		}
	})
}

func TestValidateCRMOption(t *testing.T) {
	t.Run("should return true when CRM option is valid", func(t *testing.T) {
		// arrange
		validOptions := []string{"1", "2", " 1 ", " 2 "}

		// act & assert
		for _, option := range validOptions {
			assert.True(t, ValidateCRMOption(option), "CRM option %s should be valid", option)
		}
	})

	t.Run("should return false when CRM option is invalid", func(t *testing.T) {
		// arrange
		invalidOptions := []string{"0", "3", "a", "", "12", "1a"}

		// act & assert
		for _, option := range invalidOptions {
			assert.False(t, ValidateCRMOption(option), "CRM option %s should be invalid", option)
		}
	})
}

func TestValidateBrazilianState(t *testing.T) {
	t.Run("should return true when state abbreviation is valid", func(t *testing.T) {
		// arrange
		validStates := []string{"SP", "RJ", "MG", "RS", "PR", "sp", "rj", " SP ", " rj "}

		// act & assert
		for _, state := range validStates {
			assert.True(t, ValidateBrazilianState(state), "State %s should be valid", state)
		}
	})

	t.Run("should return true when full state name is valid", func(t *testing.T) {
		// arrange
		validStates := []string{"São Paulo", "Rio de Janeiro", "Minas Gerais", "são paulo", " Rio de Janeiro "}

		// act & assert
		for _, state := range validStates {
			assert.True(t, ValidateBrazilianState(state), "State %s should be valid", state)
		}
	})

	t.Run("should return false when state is invalid", func(t *testing.T) {
		// arrange
		invalidStates := []string{"XX", "ZZ", "Invalid State", "", "123", "SP1"}

		// act & assert
		for _, state := range invalidStates {
			assert.False(t, ValidateBrazilianState(state), "State %s should be invalid", state)
		}
	})
}

func TestNormalizeBrazilianState(t *testing.T) {
	t.Run("should return abbreviation when state abbreviation is provided", func(t *testing.T) {
		// arrange & act & assert
		assert.Equal(t, "SP", NormalizeBrazilianState("SP"))
		assert.Equal(t, "RJ", NormalizeBrazilianState("rj"))
		assert.Equal(t, "MG", NormalizeBrazilianState(" mg "))
	})

	t.Run("should return abbreviation when full state name is provided", func(t *testing.T) {
		// arrange & act & assert
		assert.Equal(t, "SP", NormalizeBrazilianState("São Paulo"))
		assert.Equal(t, "RJ", NormalizeBrazilianState("rio de janeiro"))
		assert.Equal(t, "MG", NormalizeBrazilianState(" Minas Gerais "))
	})

	t.Run("should return input as-is when state is not found", func(t *testing.T) {
		// arrange & act & assert
		assert.Equal(t, "INVALID", NormalizeBrazilianState("Invalid"))
		assert.Equal(t, "XX", NormalizeBrazilianState("XX"))
	})
}

func TestValidateCityName(t *testing.T) {
	t.Run("should return true when city name is valid", func(t *testing.T) {
		// arrange
		validCities := []string{
			"São Paulo",
			"Rio de Janeiro",
			"Belo Horizonte",
			"Porto Alegre",
			"Ribeirão Preto",
			"São José dos Campos",
			"Feira de Santana",
			"Campos dos Goytacazes",
		}

		// act & assert
		for _, city := range validCities {
			assert.True(t, ValidateCityName(city), "City %s should be valid", city)
		}
	})

	t.Run("should return false when city name is invalid", func(t *testing.T) {
		// arrange
		invalidCities := []string{
			"", // empty
			"A", // too short
			"São Paulo123", // contains numbers
			"City@Name", // contains special characters
			"A very long city name that exceeds the maximum allowed length for a city name in this validation function", // too long
		}

		// act & assert
		for _, city := range invalidCities {
			assert.False(t, ValidateCityName(city), "City %s should be invalid", city)
		}
	})
}

func TestGetBrazilianStates(t *testing.T) {
	t.Run("should return all Brazilian states", func(t *testing.T) {
		// act
		states := GetBrazilianStates()

		// assert
		assert.NotEmpty(t, states)
		assert.Equal(t, "São Paulo", states["SP"])
		assert.Equal(t, "Rio de Janeiro", states["RJ"])
		assert.Equal(t, "Minas Gerais", states["MG"])
		assert.Len(t, states, 27) // 26 states + 1 federal district
	})
}

func TestGetMainMenuOptions(t *testing.T) {
	t.Run("should return all main menu options", func(t *testing.T) {
		// act
		options := GetMainMenuOptions()

		// assert
		assert.NotEmpty(t, options)
		assert.Equal(t, "Já tenho uma empresa médica constituída", options["1"])
		assert.Equal(t, "Quero abrir uma empresa", options["2"])
		assert.Equal(t, "Gostaria de tirar dúvidas", options["3"])
		assert.Equal(t, "Outros", options["4"])
		assert.Len(t, options, 4)
	})
}

func TestGetCRMOptions(t *testing.T) {
	t.Run("should return all CRM options", func(t *testing.T) {
		// act
		options := GetCRMOptions()

		// assert
		assert.NotEmpty(t, options)
		assert.Equal(t, "Já tenho CRM", options["1"])
		assert.Equal(t, "Ainda não possuo CRM", options["2"])
		assert.Len(t, options, 2)
	})
}