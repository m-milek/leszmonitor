package util

import (
	"fmt"
	"net/http"
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p Pagination) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be at least 1")
	}
	if p.PerPage < 1 {
		return fmt.Errorf("per_page must be at least 1")
	}
	return nil
}

func PaginationFromRequest(r *http.Request) (*Pagination, error) {
	query := r.URL.Query()

	page := 1
	perPage := 20

	if pageStr := query.Get("page"); pageStr != "" {
		_, err := fmt.Sscanf(pageStr, "%d", &page)
		if err != nil {
			return nil, fmt.Errorf("invalid page parameter: %w", err)
		}
	}

	if perPageStr := query.Get("per_page"); perPageStr != "" {
		_, err := fmt.Sscanf(perPageStr, "%d", &perPage)
		if err != nil {
			return nil, fmt.Errorf("invalid per_page parameter: %w", err)
		}
	}

	pagination := &Pagination{
		Page:    page,
		PerPage: perPage,
	}

	if err := pagination.Validate(); err != nil {
		return nil, err
	}
	return pagination, nil
}
