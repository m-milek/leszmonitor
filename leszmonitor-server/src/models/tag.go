package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/util"
)

type Tag struct {
	util.Timestamps
	ID          uuid.UUID `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Description string    `json:"description" db:"description"`
	ColorHex    string    `json:"colorHex"    db:"color_hex"`
}

func (t *Tag) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if t.ColorHex == "" {
		return fmt.Errorf("tag colorHex cannot be empty")
	}

	isValidHexColor := regexp.MustCompile(`^#[0-9a-f]{6}$`).MatchString(t.ColorHex)
	if !isValidHexColor {
		return fmt.Errorf("tag colorHex must be a valid hex color code")
	}

	return nil
}

func (t *Tag) Normalize() {
	t.Name = strings.TrimSpace(t.Name)
	t.ColorHex = strings.TrimSpace(t.ColorHex)
	t.ColorHex = strings.ToLower(t.ColorHex)

	// Convert shorthand hex color to full format (e.g., #abc -> #aabbcc)
	if len(t.ColorHex) == 4 {
		t.ColorHex = fmt.Sprintf("#%c%c%c%c%c%c",
			t.ColorHex[1], t.ColorHex[1],
			t.ColorHex[2], t.ColorHex[2],
			t.ColorHex[3], t.ColorHex[3])
	}
}
