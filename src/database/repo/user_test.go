package repo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gorm.io/gorm"
	"testing"
)

func TestNewUser(t *testing.T) {
	db, down := database.NewTxForTest(t)
	defer down(t)

	type args struct {
		email    string
		password string
		name     *string
		picture  *string
		db       *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create user",
			args:    args{database.TestUser.Email, database.TestUser.PlaintextPassword, &database.TestUser.Name, &database.TestUser.Picture, db},
			wantErr: false,
		},
		{
			name:    "Create user with blank email",
			args:    args{"", database.TestUser.PlaintextPassword, &database.TestUser.Name, &database.TestUser.Picture, db},
			wantErr: true,
		},
		{
			name:    "Create user with blank password",
			args:    args{database.TestUser.Email, "", &database.TestUser.Name, &database.TestUser.Picture, db},
			wantErr: true,
		},
		{
			name:    "Create user with only email and password",
			args:    args{database.TestUser.Email, database.TestUser.PlaintextPassword, nil, nil, db},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := NewUser(tt.args.email, tt.args.password, tt.args.name, tt.args.picture, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				assert.Nil(t, gotUser)
			} else {
				assert.NotNil(t, gotUser.Password)
				assert.NotNil(t, gotUser.ID)
				assert.NotNil(t, gotUser.CreatedAt)
				assert.Equal(t, tt.args.email, gotUser.Email)
				assert.Equal(t, *tt.args.name, gotUser.Name)
				assert.Equal(t, *tt.args.picture, gotUser.Picture)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, down := database.NewTxForTest(t)
	defer down(t)
	testUser, err := NewUser(database.TestUser.Email, database.TestUser.PlaintextPassword, &database.TestUser.Name, &database.TestUser.Picture, db)
	assert.NoError(t, err)

	type args struct {
		email string
		db    *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.User
		wantErr bool
	}{
		{"Get an existing user", args{testUser.Email, db}, testUser, false},
		{"Get a non-existing user", args{"non.existing@user.com", db}, &schema.User{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := GetUserByEmail(tt.args.email, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				assert.Nil(t, gotUser)
			} else {
				require.NotNil(t, gotUser)
				gotUser.CreatedAt = tt.want.CreatedAt
				gotUser.UpdatedAt = tt.want.UpdatedAt
				assert.EqualValues(t, tt.want, gotUser)
			}
		})
	}
}
