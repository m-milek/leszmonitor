package monitors

import (
	"encoding/json"
	"testing"

	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalConfigFromBytes(t *testing.T) {
	t.Run("HTTP Config", func(t *testing.T) {
		config := HttpConfig{
			Method: "GET",
			URL:    "http://example.com",
		}
		bytes, _ := json.Marshal(config)

		parsed, err := UnmarshalConfigFromBytes(shared.HttpConfigType, bytes)
		assert.NoError(t, err)
		assert.IsType(t, &HttpConfig{}, parsed)
		assert.Equal(t, config.URL, parsed.(*HttpConfig).URL)
	})

	t.Run("TCP Config", func(t *testing.T) {
		config := TCPConfig{
			Host: "example.com",
			Port: 80,
		}
		bytes, _ := json.Marshal(config)

		parsed, err := UnmarshalConfigFromBytes(shared.TCPConfigType, bytes)
		assert.NoError(t, err)
		assert.IsType(t, &TCPConfig{}, parsed)
		assert.Equal(t, config.Host, parsed.(*TCPConfig).Host)
	})

	t.Run("Unknown Config", func(t *testing.T) {
		_, err := UnmarshalConfigFromBytes("unknown", []byte("{}"))
		assert.Error(t, err)
	})
}
