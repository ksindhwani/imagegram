package httputils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func GetRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	return body, err
}

func WriteErrorResponse(w http.ResponseWriter, err Error) {
	errorTrace := map[string]string{
		"message": err.Message,
		"error":   err.Err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(errorTrace)
}

func GetUrlParam(r *http.Request, paramName string) (string, error) {
	vars := mux.Vars(r)
	param, ok := vars[paramName]
	if !ok {
		return "", errors.New("param " + paramName + " is missing in url")
	}
	return param, nil
}
