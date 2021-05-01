package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var generatedToken string

func TestGenerateToken(t *testing.T) {
	jwtWrapper := Service{
		SecretKey:       "secretKey",
		Issuer:          "Service",
		ExpirationHours: 24,
	}

	token, err := jwtWrapper.GenerateAccessToken(0, "jwt@email.com", "", "", "")
	assert.NoError(t, err)
	generatedToken = token
}

func TestValidateToken(t *testing.T) {
	jwtWrapper := Service{
		SecretKey: "secretKey",
		Issuer:    "Service",
	}

	claims, err := jwtWrapper.ValidateAccessToken(generatedToken)
	assert.NoError(t, err)

	assert.Equal(t, "jwt@email.com", claims.Email)
	assert.Equal(t, "Service", claims.Issuer)
}
