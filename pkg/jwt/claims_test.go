package jwt

import (
	"testing"
	"time"
)

func TestStandardClaims_Verify(t *testing.T) {
	tests := []struct {
		desc       string
		claims     StandardClaims
		verifyTime *Time
		err        error
	}{
		{
			desc:       "empty claims should be always valid",
			claims:     StandardClaims{},
			verifyTime: NewTime(time.Now()),
			err:        nil,
		},
		{
			desc: "claims not activated yet",
			claims: StandardClaims{
				NotBefore: NewTime(time.Now().Add(time.Hour)),
			},
			verifyTime: NewTime(time.Now()),
			err:        ErrNotBefore,
		},
		{
			desc: "claims expired",
			claims: StandardClaims{
				ExpiresAt: NewTime(time.Now().Add(-time.Hour)),
			},
			verifyTime: NewTime(time.Now()),
			err:        ErrExpired,
		},
		{
			desc: "complete claims also valid",
			claims: StandardClaims{
				ID:        "123",
				Issuer:    "just for func",
				Subject:   "subject",
				Audience:  "service",
				IssuedAt:  NewTime(time.Now()),
				ExpiresAt: NewTime(time.Now().Add(time.Hour)),
				NotBefore: NewTime(time.Now()),
			},
			verifyTime: NewTime(time.Now()),
			err:        nil,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.desc, func(t *testing.T) {
			if err := tc.claims.Valid(tc.verifyTime); err != tc.err {
				t.Errorf("expecting %v but got %v", tc.err, err)
			}
		})
	}

}
