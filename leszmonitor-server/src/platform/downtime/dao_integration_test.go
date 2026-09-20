package downtime

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDowntimeIntegrationTest(t *testing.T) (context.Context, IAppDowntimeDAO) {
	ctx, _, _ := testsupport.Setup(t)
	return ctx, NewAppDowntimeDAO(db.Get().Querier())
}

func insertTestDowntime(ctx context.Context, t *testing.T, dao IAppDowntimeDAO, startedAt, endedAt time.Time) *AppDowntime {
	t.Helper()
	inserted, err := dao.InsertDowntime(ctx, &AppDowntime{
		ID:        uuid.New(),
		StartedAt: startedAt,
		EndedAt:   endedAt,
	})
	require.NoError(t, err)
	return inserted
}

func TestIntegration_AppDowntimeDAO_Heartbeat(t *testing.T) {
	t.Run("Returns ErrNotFound when no heartbeat was ever recorded", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)

		beat, err := dao.GetHeartbeat(ctx)

		require.ErrorIs(t, err, db.ErrNotFound)
		assert.Nil(t, beat)
	})

	t.Run("Round-trips the recorded beat time", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)

		now := time.Now().UTC()
		_, err := dao.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: now})
		require.NoError(t, err)

		beat, err := dao.GetHeartbeat(ctx)

		require.NoError(t, err)
		require.NotNil(t, beat)
		assert.WithinDuration(t, now, beat.BeatAt, time.Second)
	})

	t.Run("Overwrites the previous beat instead of appending", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)

		first := time.Now().UTC().Add(-time.Hour)
		_, err := dao.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: first})
		require.NoError(t, err)

		second := time.Now().UTC()
		_, err = dao.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: second})
		require.NoError(t, err)

		beat, err := dao.GetHeartbeat(ctx)
		require.NoError(t, err)
		assert.WithinDuration(t, second, beat.BeatAt, time.Second)

		var rowCount int
		err = db.Get().Querier().QueryRowxContext(ctx, `SELECT COUNT(*) FROM heartbeat`).Scan(&rowCount)
		require.NoError(t, err)
		assert.Equal(t, 1, rowCount)
	})
}

func TestIntegration_AppDowntimeDAO_GetAllDowntimes(t *testing.T) {
	t.Run("Returns an empty slice when no windows exist", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)

		now := time.Now().UTC()
		downtimes, err := dao.GetAllDowntimes(ctx, now.Add(-time.Hour), now.Add(time.Hour))

		require.NoError(t, err)
		assert.Empty(t, downtimes)
	})

	t.Run("Round-trips a stored window", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)

		now := time.Now().UTC()
		startedAt := now.Add(-30 * time.Minute)
		endedAt := now.Add(-20 * time.Minute)
		inserted := insertTestDowntime(ctx, t, dao, startedAt, endedAt)

		downtimes, err := dao.GetAllDowntimes(ctx, now.Add(-time.Hour), now)

		require.NoError(t, err)
		require.Len(t, downtimes, 1)
		assert.Equal(t, inserted.ID, downtimes[0].ID)
		assert.WithinDuration(t, startedAt, downtimes[0].StartedAt, time.Second)
		assert.WithinDuration(t, endedAt, downtimes[0].EndedAt, time.Second)
	})

	t.Run("Includes windows overlapping the query bounds and excludes disjoint ones", func(t *testing.T) {
		ctx, dao := setupDowntimeIntegrationTest(t)

		now := time.Now().UTC()
		from := now.Add(-time.Hour)
		to := now

		before := insertTestDowntime(ctx, t, dao, now.Add(-3*time.Hour), now.Add(-2*time.Hour))
		straddlesStart := insertTestDowntime(ctx, t, dao, now.Add(-90*time.Minute), now.Add(-30*time.Minute))
		contained := insertTestDowntime(ctx, t, dao, now.Add(-20*time.Minute), now.Add(-10*time.Minute))
		straddlesEnd := insertTestDowntime(ctx, t, dao, now.Add(-5*time.Minute), now.Add(5*time.Minute))
		after := insertTestDowntime(ctx, t, dao, now.Add(time.Hour), now.Add(2*time.Hour))

		downtimes, err := dao.GetAllDowntimes(ctx, from, to)

		require.NoError(t, err)
		returnedIDs := make([]uuid.UUID, 0, len(downtimes))
		for _, d := range downtimes {
			returnedIDs = append(returnedIDs, d.ID)
		}

		assert.Equal(t, []uuid.UUID{straddlesStart.ID, contained.ID, straddlesEnd.ID}, returnedIDs)
		assert.NotContains(t, returnedIDs, before.ID)
		assert.NotContains(t, returnedIDs, after.ID)
	})
}
