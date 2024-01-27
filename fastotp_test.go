package fastotp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	httpclient "github.com/CeoFred/fast-otp/lib"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/stretchr/testify.v1/require"
)

var (
	mockedResponse = `{
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
	}`

	mockedValidationResponse = `{
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
	}`

	mockedError = `{
		"message": "something isn't right"
	}`

	mockedErrorWithSomeError = `{
		"message": "something isn't right"
		"errors":[
			"some error"
		]
	}`

	mockAPIKey = "test_api_key"
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	code = m.Run()
}

func TestMarshalling(t *testing.T) {
	var resp *OTPResponse
	err := json.NewDecoder(bytes.NewBuffer([]byte(mockedResponse))).Decode(&resp)
	require.NoError(t, err)

	fmt.Printf("\n\nID is: %s\n\n", resp.OTP.ID)
	assert.Equal(t, "test@example.com", resp.OTP.DeliveryDetails.Email)
	assert.Equal(t, OTPStatusPending, resp.OTP.Status)
	assert.Equal(t, OTPTypeAlphaNumeric, resp.OTP.Type)
	assert.Equal(t, "9b202659-fee7-46ab-836b-cdd310c4f327", resp.OTP.ID)
	assert.Equal(t, "test_identifier", resp.OTP.Identifier)
}

func TestGenerateOTP(t *testing.T) {
	reset := mockHttpRequest(http.StatusOK, http.MethodPost, "/generate", mockedResponse)
	defer reset()
	fastOtp := NewFastOTP(mockAPIKey)

	delivery := OTPDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        OTPTypeAlphaNumeric,
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	require.NoError(t, err)
	require.NotNil(t, otp)

	assert.Equal(t, "test@example.com", otp.DeliveryDetails.Email)
	assert.Equal(t, "test_identifier", otp.Identifier)
	assert.Equal(t, OTPStatusPending, otp.Status)
	assert.Equal(t, OTPTypeAlphaNumeric, otp.Type)
	assert.Equal(t, "9b202659-fee7-46ab-836b-cdd310c4f327", otp.ID)
}

func TestGenerateOTP_500Response(t *testing.T) {
	reset := mockErrorHttpRequest(500, http.MethodPost, "/generate", "")
	defer reset()

	fastOtp := NewFastOTP(mockAPIKey)

	delivery := OTPDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        OTPTypeAlphaNumeric,
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)
}

func TestGenerateOTP_400Response(t *testing.T) {
	reset := mockErrorHttpRequest(http.StatusBadRequest, http.MethodPost, "/generate", mockedError)
	defer reset()

	fastOtp := NewFastOTP(mockAPIKey)

	delivery := OTPDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        OTPTypeAlphaNumeric,
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)
	assert.Contains(t, err.Error(), "something isn't right")
}

func TestGenerateOTP_401ErrorResponse(t *testing.T) {
	fastOtp := &FastOTP{
		client: mockedHTTPClient{
			GetFunc: func(id string) (*http.Response, error) {
				var otpResponse *OTPResponse
				err := json.NewDecoder(bytes.NewBuffer([]byte(mockedValidationResponse))).Decode(&otpResponse)
				if err != nil {
					log.Fatal(err)
				}
				return httpmock.NewJsonResponse(401,
					&ErrorResponse{
						Errors:  map[string][]string{"error": {"something isn't right"}},
						Message: "something isn't right"})
			},
			PostFunc: func(endpoint string, payload interface{}) (*http.Response, error) {
				return nil, errors.New("something isn't right")
			},
		},
	}
	delivery := OTPDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        OTPTypeAlphaNumeric,
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	require.Error(t, err)
	require.Nil(t, otp)
}

func TestGenerateOTP_401MoreErrorResponse(t *testing.T) {
	fastOtp := &FastOTP{
		client: mockedHTTPClient{
			PostFunc: func(endpoint string, payload interface{}) (*http.Response, error) {
				return httpmock.NewJsonResponse(401,
					&ErrorResponse{
						Errors:  map[string][]string{"error": {"something isn't right"}},
						Message: "something isn't right"})
			},
		},
	}
	delivery := OTPDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        OTPTypeAlphaNumeric,
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	require.Error(t, err)
	require.Nil(t, otp)
}

func TestGenerateOTP_401Response(t *testing.T) {
	reset := mockErrorHttpRequest(http.StatusBadRequest, http.MethodPost, "/generate", mockedErrorWithSomeError)
	defer reset()

	fastOtp := NewFastOTP(mockAPIKey)

	delivery := OTPDelivery{
		"email": "test@example.com",
	}

	payload := GenerateOTPPayload{
		Delivery:    delivery,
		Identifier:  "test_identifier",
		TokenLength: 6,
		Type:        OTPTypeAlphaNumeric,
		Validity:    120,
	}

	otp, err := fastOtp.GenerateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)
	assert.Contains(t, err.Error(), "something isn't right")

}

func TestValidateOTP(t *testing.T) {
	reset := mockHttpRequest(http.StatusOK, http.MethodPost, "/validate", mockedValidationResponse)
	defer reset()

	fastOtp := NewFastOTP(mockAPIKey)

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	require.NoError(t, err)
	require.NotNil(t, otp)
	assert.Equal(t, "test_identifier", otp.Identifier)
	assert.Equal(t, OTPStatusValidated, otp.Status)
}

func TestValidateOTP_500Response(t *testing.T) {
	reset := mockErrorHttpRequest(500, http.MethodPost, "/validate", "")
	defer reset()

	fastOtp := NewFastOTP(mockAPIKey)

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)

}

func TestValidateOTP_400Response(t *testing.T) {
	reset := mockErrorHttpRequest(http.StatusBadRequest, http.MethodPost, "/validate", mockedError)
	defer reset()
	fastOtp := NewFastOTP(mockAPIKey)

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)
}

func TestValidateOTP_401Response(t *testing.T) {
	reset := mockErrorHttpRequest(http.StatusUnauthorized, http.MethodPost, "/validate", mockedErrorWithSomeError)
	defer reset()

	fastOtp := NewFastOTP(mockAPIKey)

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)
}

func TestValidateOTP_401MoreErrorResponse(t *testing.T) {
	fastOtp := &FastOTP{
		client: mockedHTTPClient{
			PostFunc: func(endpoint string, payload interface{}) (*http.Response, error) {
				return httpmock.NewJsonResponse(401,
					&ErrorResponse{
						Errors:  map[string][]string{"error": {"something isn't right"}},
						Message: "something isn't right"})
			},
		},
	}

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	require.Error(t, err)
	assert.Nil(t, otp)
}

func TestValidateOTP_POSTError(t *testing.T) {
	fastOtp := &FastOTP{
		client: mockedHTTPClient{
			PostFunc: func(endpoint string, payload interface{}) (*http.Response, error) {
				return nil, errors.New("some error")
			},
		},
	}

	payload := ValidateOTPPayload{
		Identifier: "test_identifier",
		Token:      "123456",
	}

	otp, err := fastOtp.ValidateOTP(payload)
	require.EqualError(t, err, "some error")
	assert.Nil(t, otp)
}

func TestFastOtp_GetOtp(t *testing.T) {
	fastOtp := &FastOTP{
		APIKey:  "",
		BaseURL: "",
		client: mockedHTTPClient{
			GetFunc: func(id string) (*http.Response, error) {
				var otpResponse *OTPResponse
				err := json.NewDecoder(bytes.NewBuffer([]byte(mockedValidationResponse))).Decode(&otpResponse)
				if err != nil {
					log.Fatal(err)
				}
				return httpmock.NewJsonResponse(200, otpResponse)
			},
		},
	}

	otp, err := fastOtp.GetOtp("test")
	require.NoError(t, err)
	require.NotNil(t, otp)

	//assert.Equal(t, "123456", otp.T)
	assert.Equal(t, OTPStatusValidated, otp.Status)
}

func TestFastOtp_GetOtpWithError(t *testing.T) {
	fastOtp := &FastOTP{
		client: mockedHTTPClient{
			GetFunc: func(id string) (*http.Response, error) {
				var otpResponse *OTPResponse
				err := json.NewDecoder(bytes.NewBuffer([]byte(mockedValidationResponse))).Decode(&otpResponse)
				if err != nil {
					log.Fatal(err)
				}
				return httpmock.NewJsonResponse(401, &ErrorResponse{Message: "something isn't right"})
			},
		},
	}

	otp, err := fastOtp.GetOtp("test")
	require.Error(t, err)
	require.Nil(t, otp)
}

func TestFastOtp_GetOtpWithMoreError(t *testing.T) {
	fastOtp := &FastOTP{
		client: mockedHTTPClient{
			GetFunc: func(id string) (*http.Response, error) {
				var otpResponse *OTPResponse
				err := json.NewDecoder(bytes.NewBuffer([]byte(mockedValidationResponse))).Decode(&otpResponse)
				if err != nil {
					log.Fatal(err)
				}
				return httpmock.NewJsonResponse(401,
					&ErrorResponse{
						Errors:  map[string][]string{"error": {"something isn't right"}},
						Message: "something isn't right"})
			},
		},
	}

	otp, err := fastOtp.GetOtp("test")
	require.Error(t, err)
	require.Nil(t, otp)
}

func mockHttpRequest(code int, method, path, response string) func() {
	httpmock.ActivateNonDefault(httpclient.FastOTPClient)
	httpmock.RegisterResponder(method, path, func(req *http.Request) (*http.Response, error) {
		// this can be used for method/path specific assertions
		// if req.Method == http.MethodPost {
		// 	body, err := io.ReadAll(req.Body)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }

		var otpResponse *OTPResponse
		err := json.NewDecoder(bytes.NewBuffer([]byte(response))).Decode(&otpResponse)
		if err != nil {
			log.Fatalln(err)
		}

		return httpmock.NewJsonResponse(code, otpResponse)
	})
	return httpmock.DeactivateAndReset
}

func mockErrorHttpRequest(code int, method, path, response string) func() {
	httpmock.ActivateNonDefault(httpclient.FastOTPClient)
	httpmock.RegisterResponder(method, path, func(req *http.Request) (*http.Response, error) {
		// this can be used for method/path specific assertions
		// if req.Method == http.MethodPost {
		// 	body, err := io.ReadAll(req.Body)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }

		var (
			errorResponse ErrorResponse
		)
		err := json.NewDecoder(bytes.NewBuffer([]byte(mockedError))).Decode(&errorResponse)
		if err != nil {
			return httpmock.NewStringResponse(http.StatusInternalServerError, mockedError), nil
		}

		return httpmock.NewJsonResponse(code, errorResponse)
	})
	return httpmock.DeactivateAndReset
}
