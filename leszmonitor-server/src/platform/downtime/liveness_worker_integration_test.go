package downtime

import (
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_HeartbeatWorker_Beat(t *testing.T) {
	t.Run("Records the first heartbeat without reporting downtime", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)
		worker := NewHeartbeatWorker(db.Get())

		require.NoError(t, worker.beat(ctx))

		beat, err := dao.GetHeartbeat(ctx)
		require.NoError(t, err)
		assert.WithinDuration(t, time.Now().UTC(), beat.BeatAt, 5*time.Second)

		downtimes, err := dao.GetAllDowntimes(ctx, time.Now().UTC().Add(-time.Hour), time.Now().UTC().Add(time.Hour))
		require.NoError(t, err)
		assert.Empty(t, downtimes)
	})

	t.Run("Does not report downtime when the previous beat is within the threshold", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)
		worker := NewHeartbeatWorker(db.Get())

		recent := time.Now().UTC().Add(-durationBetweenBeats)
		_, err := dao.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: recent})
		require.NoError(t, err)

		require.NoError(t, worker.beat(ctx))

		downtimes, err := dao.GetAllDowntimes(ctx, time.Now().UTC().Add(-time.Hour), time.Now().UTC().Add(time.Hour))
		require.NoError(t, err)
		assert.Empty(t, downtimes)
	})

	t.Run("Records a downtime window when the previous beat is older than the threshold", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)
		worker := NewHeartbeatWorker(db.Get())

		stale := time.Now().UTC().Add(-time.Hour)
		_, err := dao.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: stale})
		require.NoError(t, err)

		beatAt := time.Now().UTC()
		require.NoError(t, worker.beat(ctx))

		downtimes, err := dao.GetAllDowntimes(ctx, stale.Add(-time.Hour), beatAt.Add(time.Hour))
		require.NoError(t, err)
		require.Len(t, downtimes, 1)
		assert.WithinDuration(t, stale, downtimes[0].StartedAt, time.Second)
		assert.WithinDuration(t, beatAt, downtimes[0].EndedAt, 5*time.Second)

		beat, err := dao.GetHeartbeat(ctx)
		require.NoError(t, err)
		assert.WithinDuration(t, beatAt, beat.BeatAt, 5*time.Second)
	})

	t.Run("Does not record the same gap twice on a subsequent beat", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)
		worker := NewHeartbeatWorker(db.Get())

		stale := time.Now().UTC().Add(-time.Hour)
		_, err := dao.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: stale})
		require.NoError(t, err)

		require.NoError(t, worker.beat(ctx))
		require.NoError(t, worker.beat(ctx))

		downtimes, err := dao.GetAllDowntimes(ctx, stale.Add(-time.Hour), time.Now().UTC().Add(time.Hour))
		require.NoError(t, err)
		assert.Len(t, downtimes, 1)
	})
}
