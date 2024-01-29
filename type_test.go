package fastotp

import "testing"

func TestOTPType_String(t *testing.T) {
	tests := []struct {
		name string
		o    OTPType
		want string
	}{
		{
			name: "TestOTPType_String",
			o:    OTPTypeNumeric,
			want: "numeric",
		},
		{
			name: "TestOTPType_String",
			o:    OTPTypeAlpha,
			want: "alpha",
		},
		{
			name: "TestOTPType_String",
			o:    OTPTypeAlphaNumeric,
			want: "alpha_numeric",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("OTPType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOTPStatus_String(t *testing.T) {
	tests := []struct {
		name string
		o    OTPStatus
		want string
	}{
		{
			name: "TestOTPStatus_String",
			o:    OTPStatusPending,
			want: "pending",
		},
		{
			name: "TestOTPStatus_String",
			o:    OTPStatusValidated,
			want: "validated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("OTPStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
