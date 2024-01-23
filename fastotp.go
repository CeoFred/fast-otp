package fastotp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	// "io/ioutil"

	"github.com/CeoFred/fast-otp/lib"
)

const (
	BaseUrl = "https://api.fastotp.co"
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
	DeliveryDetails DeliveryDetails `json:"delivery_details"`
	DeliveryMethods []string        `json:"delivery_methods"`
	ExpiresAt       time.Time       `json:"expires_at"`
	ID              string          `json:"id"`
	Identifier      string          `json:"identifier"`
	Status          string          `json:"status"`
	Type            string          `json:"type"`
	UpdatedAt       time.Time       `json:"updated_at"`
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
	TokenLength int         `json:"token_length"`
	Type        string      `json:"type"`
	Validity    int         `json:"validity"`
}

func NewFastOTP(apiKey string) *FastOtp {
	return &FastOtp{APIKey: &apiKey, BaseURL: BaseUrl}
}

func (f *FastOtp) GenerateOTP(payload GenerateOTPPayload) (*OTP, error) {

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

	var otpResponse OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	fmt.Println(otpResponse, "<<< ID")

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
