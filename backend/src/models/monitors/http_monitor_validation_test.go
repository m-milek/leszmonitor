package monitors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpConfigValidate(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		config := HttpConfig{
			Method:              "GET",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		assert.NoError(t, config.validate())
	})

	t.Run("Empty URL", func(t *testing.T) {
		config := HttpConfig{
			Method:              "GET",
			URL:                 "",
			ExpectedStatusCodes: []int{200},
		}
		assert.Error(t, config.validate())
		assert.Contains(t, config.validate().Error(), "URL cannot be empty")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		config := HttpConfig{
			Method:              "GET",
			URL:                 "invalid-url",
			ExpectedStatusCodes: []int{200},
		}
		assert.Error(t, config.validate())
		assert.Contains(t, config.validate().Error(), "invalid URL")
	})

	t.Run("Empty Method", func(t *testing.T) {
		config := HttpConfig{
			Method:              "",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		assert.Error(t, config.validate())
		assert.Contains(t, config.validate().Error(), "method cannot be empty")
	})

	t.Run("Invalid Method", func(t *testing.T) {
		config := HttpConfig{
			Method:              "INVALID",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		assert.Error(t, config.validate())
		assert.Contains(t, config.validate().Error(), "invalid HTTP method")
	})

	t.Run("Invalid Body Regex", func(t *testing.T) {
		config := HttpConfig{
			Method:              "GET",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
			ExpectedBodyRegex:   "[invalid-regex",
		}
		assert.Error(t, config.validate())
		assert.Contains(t, config.validate().Error(), "invalid body regex")
	})
}
