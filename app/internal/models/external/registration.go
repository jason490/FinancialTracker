package external

// RegistrationConfigResponse describes whether invite codes are required to register.
type RegistrationConfigResponse struct {
	RegistrationCodeRequired bool  `json:"registration_code_required"`
	CodeExpiresInSeconds     int64 `json:"code_expires_in_seconds"`
}

// CreateRegistrationCodeResponse returns a newly issued invite code once.
type CreateRegistrationCodeResponse struct {
	Code      string `json:"code"`
	ExpiresAt int64  `json:"expires_at"`
}
