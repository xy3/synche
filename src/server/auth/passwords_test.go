package auth

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	type args struct {
		storedPassword string
		inputPassword  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"different passwords", args{"password_ONE", "password_TWO"}, true},
		{"same passwords", args{"same", "same"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPassword(tt.args.storedPassword, tt.args.inputPassword); (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"hashing a password", args{"testPassword21231"}, false},
		{"hashing a 1 character password", args{"1"}, false},
		{"hashing a number password", args{"123123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidateStrongPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"short password", args{"short"}, true},
		{"long but weak password", args{"aaaaaaaaaaaaaaaaa"}, true},
		{"short strong password", args{"s$h0RtSt4roN#!"}, false},
		{"long strong password", args{"L9ngeRESh&#horsStrongPASSw0rds$h0RtSt4roN#!"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateStrongPassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStrongPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
