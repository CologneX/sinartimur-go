package utils

import (
	"encoding/json"
	"net/http"
	"sinartimur-go/pkg/dto"
)

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ErrorJSON(w http.ResponseWriter, apiError *dto.APIError) {
	WriteJSON(w, apiError.StatusCode, map[string]interface{}{"error": apiError.Details})
}
