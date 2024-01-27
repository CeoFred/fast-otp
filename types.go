package fastotp

type (

	// OTPType otp types
	OTPType string

	// OTPStatus status of otp
	OTPStatus string
)

const (
	OTPTypeUnknown OTPType = "unknown"
	// OTPTypeNumeric numbers only OTP
	OTPTypeNumeric OTPType = "numeric"
	// OTPTypeAlpha alphabet only OTP
	OTPTypeAlpha OTPType = "alpha"
	// OTPTypeAlphaNumeric combination of numbers and alphabet OTP
	OTPTypeAlphaNumeric OTPType = "alpha_numeric"

	// OTPStatusPending pending otp status
	OTPStatusPending OTPStatus = "pending"
	// OTPStatusValidated validated otp status
	OTPStatusValidated OTPStatus = "validated"
)

// String returns the string value of OTPType
func (o OTPType) String() string {
	return string(o)
}

func (o OTPStatus) String() string {
	return string(o)
}
