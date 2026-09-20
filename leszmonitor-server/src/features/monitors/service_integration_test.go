package monitors_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MonitorService_GetAllMonitors(t *testing.T) {
	t.Run("Returns every monitor in the instance", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		m1 := insertTestMonitor(ctx, t)
		m2 := insertTestMonitor(ctx, t)

		all, svcErr := monitorService.GetAllMonitors(ctx)
		require.Nil(t, svcErr)

		ids := make(map[string]bool)
		for _, m := range all {
			ids[m.ID.String()] = true
		}
		assert.True(t, ids[m1.ID.String()])
		assert.True(t, ids[m2.ID.String()])
	})

	t.Run("Returns an empty list when there are no monitors", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		all, svcErr := monitorService.GetAllMonitors(ctx)
		require.Nil(t, svcErr)
		assert.Empty(t, all)
	})
}

func TestIntegration_MonitorService_CreateMonitor(t *testing.T) {
	t.Run("Successfully creates a monitor", func(t *testing.T) {
		ctx, monitorService, _, owner := setupMonitorIntegrationTest(t)

		payload := monitors.Monitor{
			Name:        "Ping API",
			Description: "Pings our main API every minute",
			Interval:    60,
			Type:        kind.HTTPConfigType,
			ProbeConfig: "{}",
		}
		payload.GenerateSlug()

		resp, svcErr := monitorService.CreateMonitor(ctx, payload)
		require.Nil(t, svcErr)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.MonitorID)

		// Verify the monitor exists and fields match
		monitorFromDB, svcErr := monitorService.GetMonitorByID(ctx, resp.MonitorID)
		require.Nil(t, svcErr)
		assert.Equal(t, "Ping API", monitorFromDB.Name)
		assert.Equal(t, kind.HTTPConfigType, monitorFromDB.Type)
		assert.Equal(t, "ping-api", monitorFromDB.Slug)

		// Verify audit log was created
		filter := audit.AuditLogFilter{ResourceID: &monitorFromDB.ID}
		entries, dbErr := audit.NewAuditLogDAO(db.Get().Querier()).GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionCreateMonitor {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitorFromDB.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "Audit log for monitor creation not found")
	})

	t.Run("Fails with 400 Bad Request for invalid monitor payload", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		payload := monitors.Monitor{
			Name:     "",  // Empty name makes it invalid
			Interval: -10, // Invalid interval
			Type:     kind.HTTPConfigType,
		}

		resp, svcErr := monitorService.CreateMonitor(ctx, payload)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
		assert.Nil(t, resp)
	})
}

func TestIntegration_MonitorService_DeleteMonitor(t *testing.T) {
	t.Run("Successfully deletes a monitor and records audit log", func(t *testing.T) {
		ctx, monitorService, _, owner := setupMonitorIntegrationTest(t)

		// Insert a monitor to delete
		monitor := insertTestMonitor(ctx, t)

		// Delete it
		svcErr := monitorService.DeleteMonitor(ctx, monitor.ID.String())
		require.Nil(t, svcErr)

		// Verify it's gone
		_, getErr := monitorService.GetMonitorByID(ctx, monitor.ID.String())
		require.NotNil(t, getErr)
		assert.Equal(t, http.StatusNotFound, getErr.Code)

		// Verify audit log was created
		filter := audit.AuditLogFilter{ResourceID: &monitor.ID}
		entries, dbErr := audit.NewAuditLogDAO(db.Get().Querier()).GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionDeleteMonitor {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitor.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "Audit log for monitor deletion not found")
	})

	t.Run("Fails with 404 when deleting nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		// Use a random UUID
		svcErr := monitorService.DeleteMonitor(ctx, uuid.New().String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})
}

func TestIntegration_MonitorService_GetMonitorByID(t *testing.T) {
	t.Run("Successfully retrieves a monitor by ID", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, monitor.ID.String())
		require.Nil(t, svcErr)
		require.NotNil(t, retrieved)
		assert.Equal(t, monitor.ID, retrieved.ID)
		assert.Equal(t, monitor.Name, retrieved.Name)
	})

	t.Run("Fails with 404 for nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, uuid.New().String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
		assert.Nil(t, retrieved)
	})

	t.Run("Successfully retrieves a monitor by slug when the id does not parse as a UUID", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, monitor.Slug)
		require.Nil(t, svcErr)
		require.NotNil(t, retrieved)
		assert.Equal(t, monitor.ID, retrieved.ID)
		assert.Equal(t, monitor.Slug, retrieved.Slug)
	})

	t.Run("Fails with 404 for a non-UUID id that doesn't match any slug", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, "not-a-uuid")
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
		assert.Nil(t, retrieved)
	})
}

func TestIntegration_MonitorService_UpdateMonitor(t *testing.T) {
	t.Run("Successfully updates a monitor and records audit log", func(t *testing.T) {
		ctx, monitorService, _, owner := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		// Modify it
		monitor.Name = "Updated Name"
		monitor.Interval = 120

		svcErr := monitorService.UpdateMonitor(ctx, *monitor)
		require.Nil(t, svcErr)

		// Verify update in DB
		updatedMonitor, _ := monitorService.GetMonitorByID(ctx, monitor.ID.String())
		assert.Equal(t, "Updated Name", updatedMonitor.Name)
		assert.Equal(t, 120, updatedMonitor.Interval)

		// Verify audit log
		filter := audit.AuditLogFilter{ResourceID: &monitor.ID}
		entries, dbErr := audit.NewAuditLogDAO(db.Get().Querier()).GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionUpdateMonitor {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitor.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "Audit log for monitor update not found")
	})

	t.Run("Preserves state even if explicitly changed in payload", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)
		originalState := monitor.RunState

		// Try to illegally change state
		monitor.RunState = monitors.MonitorStateStopped

		svcErr := monitorService.UpdateMonitor(ctx, *monitor)
		require.Nil(t, svcErr)

		updatedMonitor, _ := monitorService.GetMonitorByID(ctx, monitor.ID.String())
		assert.Equal(t, originalState, updatedMonitor.RunState, "State should not have been updated via UpdateMonitor")
	})

	t.Run("Fails with 404 for nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		fakeMonitor := monitors.Monitor{
			ID:          uuid.New(),
			Name:        "Fake",
			Description: "Fake",
			Interval:    60,
			Type:        kind.HTTPConfigType,
			ProbeConfig: "{}",
		}

		svcErr := monitorService.UpdateMonitor(ctx, fakeMonitor)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Fails with 400 for invalid configuration", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)
		monitor.Name = "" // Invalid

		svcErr := monitorService.UpdateMonitor(ctx, *monitor)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})
}

func TestIntegration_MonitorService_UpdateMonitorStateByID(t *testing.T) {
	t.Run("Successfully updates a monitor state", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, monitor.ID, monitors.MonitorStateStopped)
		require.Nil(t, svcErr)

		retrieved, _ := monitorService.GetMonitorByID(ctx, monitor.ID.String())
		assert.Equal(t, monitors.MonitorStateStopped, retrieved.RunState)
	})

	t.Run("Returns nil if state is already the desired state", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, monitor.ID, monitors.MonitorStateStopped)
		require.Nil(t, svcErr)

		svcErr = monitorService.UpdateMonitorStateByID(ctx, monitor.ID, monitors.MonitorStateStopped)
		require.Nil(t, svcErr)

		retrieved, _ := monitorService.GetMonitorByID(ctx, monitor.ID.String())
		assert.Equal(t, monitors.MonitorStateStopped, retrieved.RunState)
	})

	t.Run("Fails with 400 for invalid monitor state", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, monitor.ID, monitors.MonitorRunState("invalid_state"))
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

	t.Run("Fails with 404 for nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, uuid.New(), monitors.MonitorStateStopped)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})
}

func TestIntegration_MonitorService_RunMonitorManuallyByID(t *testing.T) {
	t.Run("Successfully runs a monitor and records audit log", func(t *testing.T) {
		ctx, monitorService, _, owner := setupMonitorIntegrationTest(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		payload := monitors.Monitor{
			Name:        "Manual Run Target",
			Description: "Monitor run manually in tests",
			Interval:    60,
			Type:        kind.HTTPConfigType,
			ProbeConfig: fmt.Sprintf(`{"method":"GET","url":%q,"expectedStatusCodes":[200]}`, server.URL),
		}
		payload.GenerateSlug()

		created, svcErr := monitorService.CreateMonitor(ctx, payload)
		require.Nil(t, svcErr)

		monitor, svcErr := monitorService.GetMonitorByID(ctx, created.MonitorID)
		require.Nil(t, svcErr)

		svcErr = monitorService.RunMonitorManuallyByID(ctx, monitor.ID)
		require.Nil(t, svcErr)

		filter := audit.AuditLogFilter{ResourceID: &monitor.ID}
		entries, dbErr := audit.NewAuditLogDAO(db.Get().Querier()).GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionRunMonitorManually {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitor.ID.String(), entry.ResourceID.String())
				assert.True(t, entry.IsSuccess)
				assert.NotNil(t, entry.After)
				break
			}
		}
		assert.True(t, found, "Audit log for manual monitor run not found")
	})
}
