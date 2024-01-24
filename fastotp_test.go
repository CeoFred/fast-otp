package fastotp

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/generate" {
			t.Errorf("Expected path /generate, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"otp": {
				"created_at": "2024-01-19T00:24:06.000000Z",
				"delivery_details": {
					"email": "test@example.com"
				},
				"delivery_methods": [
					"email"
				],
				"expires_at": "2024-01-19T17:04:06.000000Z",
				"id": "9b202659-fee7-46ab-836b-cdd310c4f327",
				"identifier": "test_identifier",
				"status": "pending",
				"type": "alpha_numeric",
				"updated_at": "2024-01-19T00:24:06.000000Z"
			}
		}`))

		body, _ := io.ReadAll(r.Body)
		fmt.Println("Request Body:", string(body))
	}))
	defer server.Close()

	fastOtp := &FastOtp{APIKey: new(string), BaseURL: server.URL}

	delivery := OtpDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        "alpha_numeric",
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Println("Generated OTP:", otp)

	if otp.DeliveryDetails.Email != "test@example.com" {
		t.Errorf("Expected email: %s, got %s", "test@example.com", otp.DeliveryDetails.Email)
	}
	if otp.Identifier != "test_identifier" {
		t.Errorf("Expected identifier: %s, got %s", "test_identifier", otp.Identifier)
	}
}

func TestValidateOTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/validate" {
			t.Errorf("Expected path /validate, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"otp": {
				"created_at": "2024-01-19T00:24:06.000000Z",
				"delivery_details": {
					"email": "test@example.com"
				},
				"delivery_methods": [
					"email"
				],
				"expires_at": "2024-01-19T17:04:06.000000Z",
				"id": "9b202659-fee7-46ab-836b-cdd310c4f327",
				"identifier": "test_identifier",
				"status": "validated",
				"type": "alpha_numeric",
				"updated_at": "2024-01-19T00:24:06.000000Z"
			}
		}`))

		body, _ := io.ReadAll(r.Body)
		fmt.Println("Request Body:", string(body))
	}))
	defer server.Close()

	fastOtp := &FastOtp{APIKey: new(string), BaseURL: server.URL}

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if otp.Status != "validated" {
		t.Errorf("Expected OTP status to be: %s, got %s", "validated", otp.DeliveryDetails.Email)
	}
}
