package ftp

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/handlers"
	"gorm.io/gorm"
)

// Auth is used to implement ftp Auth interface
type Auth struct {
	testDB *gorm.DB
}

// CheckPasswd is used to check the user whether is correct
func (a *Auth) CheckPasswd(username, password string) (correct bool, err error) {
	if a.testDB != nil {
		_, err = handlers.LoginUser(username, password, a.testDB)
	} else {
		_, err = handlers.LoginUser(username, password, database.DB)
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
