package apiclient

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"time"
)

// Client is used to configure what API Client Synche uses. This is useful for unit tests.
var Client *Synche

func ConfigureClient(host, basePath string) {
	httptransport.DefaultTimeout = 100 * time.Second
	Client = New(httptransport.New(host, basePath, DefaultSchemes), strfmt.Default)
}
