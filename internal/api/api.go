package api

import "net/http"

type API interface {
	CheckHealth(mux *http.ServeMux)
}
