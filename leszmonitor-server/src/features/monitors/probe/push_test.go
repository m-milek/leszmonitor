package probe

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushProbe_Validate(t *testing.T) {
	t.Run("Negative grace period", func(t *testing.T) {
		probe := &PushProbe{GracePeriodSeconds: -1}
		require.Error(t, probe.Validate())
	})

	t.Run("Zero grace period", func(t *testing.T) {
		probe := &PushProbe{GracePeriodSeconds: 0}
		require.NoError(t, probe.Validate())
	})

	t.Run("Positive grace period", func(t *testing.T) {
		probe := &PushProbe{GracePeriodSeconds: 30}
		require.NoError(t, probe.Validate())
	})
}

func TestPushProbe_Run(t *testing.T) {
	probe := &PushProbe{GracePeriodSeconds: 30}

	result, err := probe.Run(context.Background(), uuid.New())
	require.Error(t, err)
	assert.Nil(t, result)
}
