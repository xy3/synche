package ftp

import (
	"github.com/goftp/server"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
)

// Factory is an FTP driver factory, it is used to generate a driver when a new connection comes
type Factory struct {
	Logger *log.Logger
}

func (f *Factory) NewDriver() (server.Driver, error) {
	return &Driver{
		db:     database.RequireNewConnection(),
		logger: f.Logger,
	}, nil
}
