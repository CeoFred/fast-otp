# Go FastOTP

FastOTP is a Go library for interacting with the FastOTP API to generate one-time passwords (OTPs) easily.

## Installation

```bash
go get -u github.com/CeoFred/fast-otp@v1.0.2
```

## Usage

### Basic Example

```go

package main

import (
	"fmt"
	"log"

	"github.com/CeoFred/fast-otp"
)

func main() {
	// Replace "your_api_key" with your actual FastOTP API key
	apiKey := "your_api_key"

	// Create an instance of FastOtp
	client := fastotp.NewFastOTP(apiKey)

	// Create context for library functions
	ctx := context.Background()

	// Define OTP generation payload
	payload := fastotp.GenerateOTPPayload{
    Delivery: OtpDelivery{
      "email": "example@example.com",
    },
		Identifier:  "user123",
		TokenLength: 6,
		Type:        "numeric",
		Validity:    120,
	}

	// Generate OTP
	otp, err := client.GenerateOTP(ctx, payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated OTP: %s\n", otp)

  // Validate OTP
  otp, err = client.ValidateOTP(
	 ctx,
     ValidateOTPPayload{
        Identifier: "user123",
        Token: "123456",
    }
  )
	if err != nil {
		log.Fatal(err)
	}

  fmt.Printf("Otp validation status: %s\n", otp.Status)
}
```

## API Documentation

For detailed information about the FastOTP API and available endpoints, refer to the [official API documentation](https://api.fastotp.co/docs).

## Configuration

- `APIKey`: Your FastOTP API key.

## Contributing

If you'd like to contribute to this project, please follow the guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the FastOTP team for providing the awesome OTP generation service.
