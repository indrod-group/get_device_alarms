package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestUpdateDevice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	device := &Device{
		Imei:             "123456789012345",
		UserName:         "test_user",
		CarOwner:         nil,
		LicenseNumber:    nil,
		Vin:              nil,
		IsTrackingAlarms: false,
		LastTimeTracked:  0,
	}

	httpmock.RegisterResponder(
		"POST",
		DEVICES_API_URL,
		httpmock.NewStringResponder(201, `{"success": true}`),
	)

	err := device.UpdateDevice()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	httpmock.RegisterResponder(
		"POST",
		DEVICES_API_URL,
		httpmock.NewStringResponder(500, `{"success": false}`),
	)

	err = device.UpdateDevice()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGenerateURL(t *testing.T) {
	device := &Device{
		Imei:             "123456789012345",
		UserName:         "test_user",
		CarOwner:         nil,
		LicenseNumber:    nil,
		Vin:              nil,
		IsTrackingAlarms: false,
		LastTimeTracked:  0,
	}

	url := device.GenerateURL()
	if url == "" {
		t.Errorf("Expected a URL but got an empty string")
	}
}
