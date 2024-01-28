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
	baseURL = "https://api.fastotp.co"
)

// FastOTP is the main struct for the FastOtp package.
type FastOTP struct {
	apiKey  string
	baseURL string
	client  HttpClient
}

// ErrorResponse is the error struct for the FastOtp package.
type ErrorResponse struct {
	Errors  map[string][]string `json:"errors"`
	Message string              `json:"message"`
}

// OTP is the struct for the OTP object.
type OTP struct {
	CreatedAt       time.Time       `json:"created_at"`
	ExpiresAt       time.Time       `json:"expires_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeliveryDetails DeliveryDetails `json:"delivery_details"`
	ID              string          `json:"id"`
	Identifier      string          `json:"identifier"`
	Status          OTPStatus       `json:"status"`
	Type            OTPType         `json:"type"`
	DeliveryMethods []string        `json:"delivery_methods"`
}

// OTPResponse is the struct for the OTP response object.
type OTPResponse struct {
	OTP OTP `json:"otp"`
}

// DeliveryDetails is the struct for the DeliveryDetails object.
type DeliveryDetails struct {
	Email string `json:"email"`
}

// OTPDelivery is the struct for the OtpDelivery object.
type OTPDelivery map[string]string

// GenerateOTPPayload is the struct for the GenerateOTPPayload object.
type GenerateOTPPayload struct {
	Delivery    OTPDelivery `json:"delivery"`
	Identifier  string      `json:"identifier"`
	Type        OTPType     `json:"type"`
	TokenLength int         `json:"token_length"`
	Validity    int         `json:"validity"`
}

// ValidateOTPPayload is the struct for the ValidateOTPPayload object.
type ValidateOTPPayload struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

// NewFastOTP creates a new FastOtp instance.
func NewFastOTP(apiKey string) *FastOTP {
	return &FastOTP{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  httpclient.NewAPIClient(baseURL, apiKey),
	}
}

func (f *FastOTP) GenerateOTP(ctx context.Context, payload GenerateOTPPayload) (*OTP, error) {
	resp, err := f.client.Post(ctx, "/generate", payload)
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

func (f *FastOTP) ValidateOTP(ctx context.Context, payload ValidateOTPPayload) (*OTP, error) {
	resp, err := f.client.Post(ctx, "/validate", payload)
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

// GetOtp gets a new otp
func (f *FastOTP) GetOtp(ctx context.Context, id string) (*OTP, error) {
	resp, err := f.client.Get(ctx, id)
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

func formatValidationError(errors map[string][]string) error {
	var errorMessage string
	for field, fieldErrors := range errors {
		for _, err := range fieldErrors {
			errorMessage += fmt.Sprintf("%s: %s\n", field, err)
		}
	}
	return fmt.Errorf("validation errors:\n%s", errorMessage)
}
