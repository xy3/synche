package config

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
)

// Client is used to configure what API Client Synche uses. This is useful for unit tests.
var Client = apiclient.Default

func ConfigureClient() error {
	transport := httptransport.New(Config.Server.Host, Config.Server.BasePath, apiclient.DefaultSchemes)
	Client = apiclient.New(transport, strfmt.Default)
	return nil
}
