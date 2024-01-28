package fastotp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpclient "github.com/CeoFred/fast-otp/lib"
)

const (
	BaseURL = "https://api.fastotp.co"
)

type FastOtp struct {
	APIKey  *string
	BaseURL string
}

type ErrorResponse struct {
	Errors  map[string][]string `json:"errors"`
	Message string              `json:"message"`
}

type OTP struct {
	CreatedAt       time.Time       `json:"created_at"`
	ExpiresAt       time.Time       `json:"expires_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeliveryDetails DeliveryDetails `json:"delivery_details"`
	ID              string          `json:"id"`
	Identifier      string          `json:"identifier"`
	Status          string          `json:"status"`
	Type            string          `json:"type"`
	DeliveryMethods []string        `json:"delivery_methods"`
}

type OTPResponse struct {
	OTP OTP `json:"otp"`
}

type DeliveryDetails struct {
	Email string `json:"email"`
}

type OtpDelivery map[string]string

type GenerateOTPPayload struct {
	Delivery    OtpDelivery `json:"delivery"`
	Identifier  string      `json:"identifier"`
	Type        string      `json:"type"`
	TokenLength int         `json:"token_length"`
	Validity    int         `json:"validity"`
}

type ValidateOTPPayload struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

func NewFastOTP(apiKey string) *FastOtp {
	return &FastOtp{APIKey: &apiKey, BaseURL: BaseURL}
}

func (f *FastOtp) GenerateOTP(ctx context.Context, payload GenerateOTPPayload) (*OTP, error) {
	cl := httpclient.NewAPIClient(f.BaseURL, *f.APIKey, ctx)
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

	var otpResponse OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse.OTP, nil
}

func (f *FastOtp) ValidateOTP(ctx context.Context, payload ValidateOTPPayload) (*OTP, error) {
	cl := httpclient.NewAPIClient(f.BaseURL, *f.APIKey, ctx)
	resp, err := cl.Post("/validate", payload)
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

	var otpResponse OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse.OTP, nil
}

func (f *FastOtp) GetOtp(ctx context.Context, id string) (*OTP, error) {
	cl := httpclient.NewAPIClient(f.BaseURL, *f.APIKey, ctx)
	resp, err := cl.Get(id)
	if err != nil {
		fmt.Println("got here")
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

	var otpResponse OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse.OTP, nil
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
