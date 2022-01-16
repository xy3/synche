package ftp

import (
	"github.com/goftp/server"
	log "github.com/sirupsen/logrus"
	server2 "github.com/xy3/synche/src/server"
)

// Factory is an FTP driver factory, it is used to generate a driver when a new connection comes
type Factory struct {
	Logger *log.Logger
}

func (f *Factory) NewDriver() (server.Driver, error) {
	return &Driver{
		db:     server2.RequireNewConnection(),
		logger: f.Logger,
	}, nil
}
