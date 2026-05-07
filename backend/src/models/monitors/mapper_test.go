package monitors

import (
	"encoding/json"
	"testing"

	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalConfigFromBytes(t *testing.T) {
	t.Run("HTTP Config", func(t *testing.T) {
		config := HttpProbe{
			Method: "GET",
			URL:    "http://example.com",
		}
		bytes, _ := json.Marshal(config)

		parsed, err := UnmarshalProbeFromBytes(shared.HttpConfigType, bytes)
		assert.NoError(t, err)
		assert.IsType(t, &HttpProbe{}, parsed)
		assert.Equal(t, config.URL, parsed.(*HttpProbe).URL)
	})

	t.Run("TCP Config", func(t *testing.T) {
		config := TCPProbe{
			Host: "example.com",
			Port: 80,
		}
		bytes, _ := json.Marshal(config)

		parsed, err := UnmarshalProbeFromBytes(shared.TCPConfigType, bytes)
		assert.NoError(t, err)
		assert.IsType(t, &TCPProbe{}, parsed)
		assert.Equal(t, config.Host, parsed.(*TCPProbe).Host)
	})

	t.Run("Unknown Config", func(t *testing.T) {
		_, err := UnmarshalProbeFromBytes("unknown", []byte("{}"))
		assert.Error(t, err)
	})
}
