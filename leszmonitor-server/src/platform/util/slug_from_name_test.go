package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugFromName_Init(t *testing.T) {
	name := "Sample Name"
	var slug SlugFromName
	slug.Init(name)

	expectedID := "sample-name"
	assert.Equal(t, expectedID, slug.Slug)
	assert.Equal(t, name, slug.Name)
}
