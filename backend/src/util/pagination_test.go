package util

import (
	"net/http"
	"testing"
)

func TestPaginationOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		perPage  int
		expected int
	}{
		{"first page", 1, 20, 0},
		{"second page", 2, 20, 20},
		{"tenth page with 10 items per page", 10, 10, 90},
		{"large page number", 100, 50, 4950},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, PerPage: tt.perPage}
			if got := p.Offset(); got != tt.expected {
				t.Errorf("Offset() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestPaginationValidate(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		perPage   int
		expectErr bool
	}{
		{"valid pagination", 1, 20, false},
		{"page 0 invalid", 0, 20, true},
		{"negative page invalid", -1, 20, true},
		{"per_page 0 invalid", 1, 0, true},
		{"negative per_page invalid", 1, -10, true},
		{"large valid values", 999, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, PerPage: tt.perPage}
			err := p.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("Validate() error = %v, expectErr = %v", err, tt.expectErr)
			}
		})
	}
}

func TestPaginationFromRequest(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expected  *Pagination
		expectErr bool
	}{
		{"no parameters uses defaults", "http://example.com/api", &Pagination{Page: 1, PerPage: 20}, false},
		{"page parameter only", "http://example.com/api?page=3", &Pagination{Page: 3, PerPage: 20}, false},
		{"per_page parameter only", "http://example.com/api?per_page=50", &Pagination{Page: 1, PerPage: 50}, false},
		{"both parameters", "http://example.com/api?page=5&per_page=30", &Pagination{Page: 5, PerPage: 30}, false},
		{"empty page parameter uses default", "http://example.com/api?page=", &Pagination{Page: 1, PerPage: 20}, false},
		{"empty per_page parameter uses default", "http://example.com/api?per_page=", &Pagination{Page: 1, PerPage: 20}, false},
		{"extra parameters ignored", "http://example.com/api?page=2&per_page=10&sort=name", &Pagination{Page: 2, PerPage: 10}, false},
		{"invalid page parameter", "http://example.com/api?page=abc", nil, true},
		{"invalid per_page parameter", "http://example.com/api?per_page=xyz", nil, true},
		{"floating point page", "http://example.com/api?page=2.5", nil, true},
		{"page 0 fails validation", "http://example.com/api?page=0", nil, true},
		{"per_page 0 fails validation", "http://example.com/api?per_page=0", nil, true},
		{"negative page fails validation", "http://example.com/api?page=-5", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)
			got, err := PaginationFromRequest(req)

			if (err != nil) != tt.expectErr {
				t.Errorf("PaginationFromRequest() error = %v, expectErr = %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr && (got.Page != tt.expected.Page || got.PerPage != tt.expected.PerPage) {
				t.Errorf("PaginationFromRequest() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
