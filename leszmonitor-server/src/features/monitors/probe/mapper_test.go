package probe

import (
	"encoding/json"
	"testing"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("HTTP Config", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "GET",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		bytes, _ := json.Marshal(config)

		parsed, err := Parse[*HTTPProbe](kind.HTTPConfigType, string(bytes))
		require.NoError(t, err)
		assert.Equal(t, config.URL, parsed.URL)
	})

	t.Run("TCP Config", func(t *testing.T) {
		config := TCPProbe{
			Host:     "example.com",
			Port:     80,
			Protocol: "tcp",
			Timeout:  1000,
		}
		bytes, _ := json.Marshal(config)

		parsed, err := Parse[*TCPProbe](kind.TCPConfigType, string(bytes))
		require.NoError(t, err)
		assert.Equal(t, config.Host, parsed.Host)
	})

	t.Run("Push Config", func(t *testing.T) {
		config := PushProbe{
			GracePeriodSeconds: 30,
		}
		bytes, _ := json.Marshal(config)

		parsed, err := Parse[*PushProbe](kind.PushConfigType, string(bytes))
		require.NoError(t, err)
		assert.Equal(t, config.GracePeriodSeconds, parsed.GracePeriodSeconds)
	})

	t.Run("Interface Type", func(t *testing.T) {
		parsed, err := Parse[Probe](kind.PushConfigType, `{"gracePeriodSeconds":30}`)
		require.NoError(t, err)
		assert.IsType(t, &PushProbe{}, parsed)
	})

	t.Run("Unknown Config", func(t *testing.T) {
		_, err := Parse[Probe]("unknown", "{}")
		require.Error(t, err)
	})

	t.Run("Unknown Field", func(t *testing.T) {
		_, err := Parse[Probe](kind.PushConfigType, `{"gracePeriodSeconds":30,"gracePeriod":30}`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unknown field "gracePeriod"`)
	})

	t.Run("Mismatched Type", func(t *testing.T) {
		_, err := Parse[*PushProbe](
			kind.HTTPConfigType,
			`{"method":"GET","url":"http://example.com","expectedStatusCodes":[200]}`,
		)
		require.Error(t, err)
		assert.NotErrorIs(t, err, ErrInvalidProbeConfig)
	})

	t.Run("Invalid Config", func(t *testing.T) {
		_, err := Parse[*PushProbe](kind.PushConfigType, `{"gracePeriodSeconds":-1}`)
		require.ErrorIs(t, err, ErrInvalidProbeConfig)
	})
}
