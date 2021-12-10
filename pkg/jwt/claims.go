package jwt

// StandardClaims is a structured version of Claims sections, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1.
type StandardClaims struct {
	ID        string `json:"jti,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Audience  string `json:"aud,omitempty"`
	IssuedAt  *Time  `json:"iat,omitempty"`
	ExpiresAt *Time  `json:"exp,omitempty"`
	NotBefore *Time  `json:"nbf,omitempty"`
}

func (s StandardClaims) Valid(at *Time) error {
	if s.ExpiresAt != nil && at.After(s.ExpiresAt.Time) {
		return ErrExpired
	}

	if s.NotBefore != nil && at.Before(s.NotBefore.Time) {
		return ErrNotBefore
	}

	return nil
}
