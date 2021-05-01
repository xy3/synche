package ftp

import "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/handlers"

// Auth is used to implement ftp Auth interface
type Auth struct{}

// CheckPasswd is used to check the user whether is correct
func (a *Auth) CheckPasswd(username, password string) (correct bool, err error) {
	_, err = handlers.LoginUser(username, password)
	if err != nil {
		return false, err
	}
	return true, nil
}
