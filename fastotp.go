package fastotp

import (
	"encoding/json"
	"fmt"

	"github.com/CeoFred/fast-otp/lib"

	"net/http"
)

type FastOtp struct {
	APIKey  *string
	BaseURL string
}

type ErrorResponse struct {
	Errors  map[string][]string `json:"errors"`
	Message string              `json:"message"`
}

type GenerateOTPPayload struct {
	Delivery struct {
		Email string `json:"email"`
	} `json:"delivery"`
	Identifier  string `json:"identifier"`
	TokenLength int    `json:"token_length"`
	Type        string `json:"type"`
	Validity    int    `json:"validity"`
}

func (f *FastOtp) NewFastOTP(apiKey string) *FastOtp {
	return &FastOtp{APIKey: &apiKey, BaseURL: "https://api.fastotp.co"}
}

func (f *FastOtp) GenerateOTP(payload GenerateOTPPayload) (*GenerateOTPPayload, error) {

	cl := httpclient.NewAPIClient(f.BaseURL, *f.APIKey)
	resp, err := cl.Post("/generate", payload)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		if len(errorResponse.Errors) > 0 {
			return nil, formatValidationError(errorResponse.Errors)
		}

		return nil, fmt.Errorf("API error: %s", errorResponse.Message)
	}
	var otpResponse GenerateOTPPayload
	if err := json.NewDecoder(resp.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse, nil
}

func formatValidationError(errors map[string][]string) error {
	var errorMessage string
	for field, fieldErrors := range errors {
		for _, err := range fieldErrors {
			errorMessage += fmt.Sprintf("%s: %s\n", field, err)
		}
	}
	return fmt.Errorf("validation errors:\n%s", errorMessage)
}
