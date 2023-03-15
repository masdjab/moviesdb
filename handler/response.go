package handler

import (
	"encoding/json"
	"net/http"

	"moviesdb.com/model"
)

type MovieResponse struct {
	Success bool          `json:"success"`
	Error   string        `json:"error"`
	Data    interface{}   `json:"data"`
}

type OneRowResponse struct {
	Success bool          `json:"success"`
	Error   string        `json:"error"`
	Movie   model.Movie   `json:"movie"`
}

type MultipleRowsResponse struct {
	Success bool          `json:"success"`
	Error   string        `json:"error"`
	Movies  []model.Movie `json:"movies"`
}


func respond(w http.ResponseWriter, status int, responseBody interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseJson, _ := json.Marshal(responseBody)
	_, _ = w.Write(responseJson)
}

func notAuthorized(w http.ResponseWriter, reason string) {
	respond(w, http.StatusUnauthorized, MovieResponse{
		Success: false,
		Error:   reason,
	})
}

func badRequest(w http.ResponseWriter, reason string) {
	respond(w, http.StatusBadRequest, MovieResponse{
		Success: false,
		Error:   reason,
	})
}

func internalError(w http.ResponseWriter, reason string) {
	respond(w, http.StatusInternalServerError, MovieResponse{
		Success: false,
		Error:   reason,
	})
}

func success(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseBody := MovieResponse{
		Success: true,
		Error:   "",
		Data:    data,
	}
	responseJson, _ := json.Marshal(responseBody)
	_, _ = w.Write(responseJson)
}
