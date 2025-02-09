package utils

import (
	"net/http"
	"strconv"
)

type PaginationParameter struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

type PaginatedHandler struct {
	handlerFunc func(http.ResponseWriter, *http.Request, int, int, string, string)
}

func NewPaginatedHandler(handlerFunc func(http.ResponseWriter, *http.Request, int, int, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = DefaultPage
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		if pageSize == 0 {
			pageSize = DefaultPageSize
		}
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		handlerFunc(w, r, page, pageSize, sortBy, sortOrder)
	}
}
