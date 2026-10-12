package util

import "time"

// Timestamps holds createdAt and updatedAt timestamps for database models.
type Timestamps struct {
	CreatedAt time.Time  `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}
