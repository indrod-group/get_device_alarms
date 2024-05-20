package main

import (
	"testing"
)

func TestGetPhoneNumbersFromAPI(t *testing.T) {
	imei := "860419050021378"

	// Llamar a la función con la URL de la API
	phoneNumbers, err := GetPhoneNumbersFromAPI(imei)
	if err != nil {
		t.Fatalf("GetPhoneNumbersFromAPI failed: %v", err)
	}

	// Comprobar que la función devolvió los números de teléfono
	if len(phoneNumbers) == 1 {
		t.Fatalf("No phone numbers returned")
	}
}
