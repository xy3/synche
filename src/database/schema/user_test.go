package schema

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"testing"
)

func TestUser_ValidateForRegistration(t *testing.T) {
	validHash := hash.MD5HashString("123")

	type fields struct {
		Email     string
		Password  string
		Name      string
		TokenHash string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"invalid email", fields{"invalid@", "Password1234!", "Test", validHash}, true},
		{"invalid password", fields{"valid@mail.com", "short", "Test", validHash}, true},
		{"invalid token hash", fields{"valid@mail.com", "Password1234!", "Test", "1234"}, true},
		{"invalid name", fields{"valid@mail.com", "Password1234!", "n", validHash}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				Name:      tt.fields.Name,
				TokenHash: tt.fields.TokenHash,
			}
			if err := user.ValidateForRegistration(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateForRegistration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isEmailValid(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "mail@test.com", true},
		{"no TLD", "mail@", false},
		{"no email user", "@gmail.com", false},
		{"too short email", "123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmailValid(tt.email)
			assert.Equal(t, tt.want, got)
		})
	}
}
