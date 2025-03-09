package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sinartimur-go/pkg/dto"
)

type PaginationResponse struct {
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	TotalItems  int         `json:"total_items"`
	Items       interface{} `json:"items"`
}

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WritePaginationJSON(w http.ResponseWriter, code int, currentPage, totalItems int, pageSize int, items interface{}) {
	totalPages := (totalItems + pageSize - 1) / pageSize
	WriteJSON(w, code, PaginationResponse{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		Items:       items,
	})
}

func ErrorJSON(w http.ResponseWriter, apiError *dto.APIError) {
	WriteJSON(w, apiError.StatusCode, map[string]interface{}{"error": apiError.Details})
}

// ToJSON converts an interface to a JSON string
func ToJSON(data interface{}) (string, *dto.APIError) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Gagal mengonversi data",
			},
		}
	}
	return string(bytes), nil
}

func TransformRoles(roles []interface{}) ([]*string, error) {
	var result []*string
	for _, role := range roles {
		strRole, ok := role.(string)
		if !ok {
			return nil, fmt.Errorf("role is not a string: %v", role)
		}
		result = append(result, &strRole)
	}
	return result, nil
}

func WriteMessage(message interface{}) interface{} {
	return map[string]string{"message": message.(string)}
}
