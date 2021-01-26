package upload

import (
	"net/http"
)

type ServerResponse struct {
	Status string
	Header http.Header
	Body   string
}
