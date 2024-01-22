package fastotp

import (
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
			"delivery": {
				"email": "test@example.com"
			},
			"identifier": "test_identifier",
			"token_length": 6,
			"type": "numeric",
			"validity": 120
		}`))
	}))
	defer server.Close()

	fastOtp := &FastOtp{APIKey: new(string), BaseURL: server.URL}

	payload := GenerateOTPPayload{
		Delivery: struct {
			Email string `json:"email"`
		}{
			Email: "test@example.com",
		},
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        "numeric",
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if otp.Delivery.Email != "test@example.com" {
		t.Errorf("Expected email: %s, got %s", "test@example.com", otp.Delivery.Email)
	}
	if otp.Identifier != "test_identifier" {
		t.Errorf("Expected identifier: %s, got %s", "test_identifier", otp.Identifier)
	}
}
