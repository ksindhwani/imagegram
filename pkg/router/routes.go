package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ksindhwani/imagegram/pkg/app"
)

func New(deps *app.Dependencies) (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/ping", PingHandler).Methods(http.MethodGet)
	return r, nil
}
