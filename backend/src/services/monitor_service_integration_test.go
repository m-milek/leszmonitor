package services

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MonitorService_CreateMonitor(t *testing.T) {
	t.Run("Successfully creates a monitor", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		require.Nil(t, err)

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		payload := monitors.Monitor{
			Name:        "Ping API",
			Description: "Pings our main API every minute",
			Interval:    60,
			Type:        consts.HttpConfigType,
			ProbeConfig: "{}",
		}
		payload.GenerateSlug()

		resp, svcErr := monitorService.CreateMonitor(ctx, auth, payload)
		require.Nil(t, svcErr)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.MonitorID)

		// Verify the monitor exists and fields match
		monitorFromDB, svcErr := monitorService.GetMonitorByID(ctx, auth, resp.MonitorID)
		require.Nil(t, svcErr)
		assert.Equal(t, "Ping API", monitorFromDB.Name)
		assert.Equal(t, project.ID, monitorFromDB.ProjectID)
		assert.Equal(t, consts.HttpConfigType, monitorFromDB.Type)
		assert.Equal(t, "ping-api", monitorFromDB.Slug)

		// Verify audit log was created
		filter := security.AuditLogFilter{ProjectID: &project.ID}
		entries, dbErr := db.Get().AuditLog().GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == security.ActionCreateMonitor {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitorFromDB.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "Audit log for monitor creation not found")
	})

	t.Run("Fails with 400 Bad Request for invalid monitor payload", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		payload := monitors.Monitor{
			Name:     "",  // Empty name makes it invalid
			Interval: -10, // Invalid interval
			Type:     consts.HttpConfigType,
		}

		resp, svcErr := monitorService.CreateMonitor(ctx, auth, payload)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
		assert.Nil(t, resp)
	})
}

func TestIntegration_MonitorService_DeleteMonitor(t *testing.T) {
	t.Run("Successfully deletes a monitor and records audit log", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		require.Nil(t, err)

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		// Insert a monitor to delete
		monitor := insertTestMonitor(t, ctx, project.ID)

		// Delete it
		svcErr := monitorService.DeleteMonitor(ctx, auth, monitor.ID.String())
		require.Nil(t, svcErr)

		// Verify it's gone
		_, getErr := monitorService.GetMonitorByID(ctx, auth, monitor.ID.String())
		require.NotNil(t, getErr)
		assert.Equal(t, http.StatusNotFound, getErr.Code)

		// Verify audit log was created
		filter := security.AuditLogFilter{ProjectID: &project.ID}
		entries, dbErr := db.Get().AuditLog().GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == security.ActionDeleteMonitor {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitor.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "Audit log for monitor deletion not found")
	})

	t.Run("Fails with 404 when deleting nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		// Use a random UUID
		svcErr := monitorService.DeleteMonitor(ctx, auth, uuid.New().String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

}

func TestIntegration_MonitorService_GetMonitorsByProjectID(t *testing.T) {
	t.Run("Successfully retrieves monitors for a project", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		require.Nil(t, err)

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		// Insert 2 monitors
		m1 := insertTestMonitor(t, ctx, project.ID)
		m2 := insertTestMonitor(t, ctx, project.ID)

		monitorsList, svcErr := monitorService.GetMonitorsByProjectID(ctx, auth)
		require.Nil(t, svcErr)
		require.Len(t, monitorsList, 2)

		foundM1, foundM2 := false, false
		for _, m := range monitorsList {
			if m.ID == m1.ID {
				foundM1 = true
			}
			if m.ID == m2.ID {
				foundM2 = true
			}
		}
		assert.True(t, foundM1)
		assert.True(t, foundM2)
	})

	t.Run("Returns empty list if project has no monitors", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Empty Project"})
		require.Nil(t, err)

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitorsList, svcErr := monitorService.GetMonitorsByProjectID(ctx, auth)
		require.Nil(t, svcErr)
		assert.Empty(t, monitorsList)
	})

}

func TestIntegration_MonitorService_GetMonitorByID(t *testing.T) {
	t.Run("Successfully retrieves a monitor by ID", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, auth, monitor.ID.String())
		require.Nil(t, svcErr)
		require.NotNil(t, retrieved)
		assert.Equal(t, monitor.ID, retrieved.ID)
		assert.Equal(t, monitor.Name, retrieved.Name)
	})

	t.Run("Fails with 404 for nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, auth, uuid.New().String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
		assert.Nil(t, retrieved)
	})

	t.Run("Fails with 400 for invalid UUID format", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		retrieved, svcErr := monitorService.GetMonitorByID(ctx, auth, "not-a-uuid")
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
		assert.Nil(t, retrieved)
	})

}

func TestIntegration_MonitorService_UpdateMonitor(t *testing.T) {
	t.Run("Successfully updates a monitor and records audit log", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)

		// Modify it
		monitor.Name = "Updated Name"
		monitor.Interval = 120

		svcErr := monitorService.UpdateMonitor(ctx, auth, *monitor)
		require.Nil(t, svcErr)

		// Verify update in DB
		updatedMonitor, _ := monitorService.GetMonitorByID(ctx, auth, monitor.ID.String())
		assert.Equal(t, "Updated Name", updatedMonitor.Name)
		assert.Equal(t, 120, updatedMonitor.Interval)

		// Verify audit log
		filter := security.AuditLogFilter{ProjectID: &project.ID}
		entries, dbErr := db.Get().AuditLog().GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == security.ActionUpdateMonitor {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, monitor.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "Audit log for monitor update not found")
	})

	t.Run("Preserves project ID and state even if explicitly changed in payload", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)
		originalState := monitor.State

		// Try to illegally change project ID and state
		monitor.ProjectID = uuid.New()
		monitor.State = monitors.MonitorStateStopped

		svcErr := monitorService.UpdateMonitor(ctx, auth, *monitor)
		require.Nil(t, svcErr)

		updatedMonitor, _ := monitorService.GetMonitorByID(ctx, auth, monitor.ID.String())
		assert.Equal(t, project.ID, updatedMonitor.ProjectID, "Project ID should not have been updated")
		assert.Equal(t, originalState, updatedMonitor.State, "State should not have been updated via UpdateMonitor")
	})

	t.Run("Fails with 404 for nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		fakeMonitor := monitors.Monitor{
			ID:          uuid.New(),
			ProjectID:   project.ID,
			Name:        "Fake",
			Description: "Fake",
			Interval:    60,
			Type:        consts.HttpConfigType,
			ProbeConfig: "{}",
		}

		svcErr := monitorService.UpdateMonitor(ctx, auth, fakeMonitor)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Fails with 400 for invalid configuration", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)
		monitor.Name = "" // Invalid

		svcErr := monitorService.UpdateMonitor(ctx, auth, *monitor)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

}

func TestIntegration_MonitorService_GetMonitorBySlugByProject(t *testing.T) {
	t.Run("Successfully retrieves a monitor by slug", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)

		retrieved, svcErr := monitorService.GetMonitorBySlugByProject(ctx, auth, monitor.Slug)
		require.Nil(t, svcErr)
		require.NotNil(t, retrieved)
		assert.Equal(t, monitor.ID, retrieved.ID)
		assert.Equal(t, monitor.Slug, retrieved.Slug)
	})

	t.Run("Fails with 404 for nonexistent slug", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		retrieved, svcErr := monitorService.GetMonitorBySlugByProject(ctx, auth, "does-not-exist")
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
		assert.Nil(t, retrieved)
	})

}

func TestIntegration_MonitorService_UpdateMonitorStateByID(t *testing.T) {
	t.Run("Successfully updates a monitor state", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, auth, monitor.ID, monitors.MonitorStateStopped)
		require.Nil(t, svcErr)

		retrieved, _ := monitorService.GetMonitorByID(ctx, auth, monitor.ID.String())
		assert.Equal(t, monitors.MonitorStateStopped, retrieved.State)
	})

	t.Run("Returns nil if state is already the desired state", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, auth, monitor.ID, monitors.MonitorStateStopped)
		require.Nil(t, svcErr)

		svcErr = monitorService.UpdateMonitorStateByID(ctx, auth, monitor.ID, monitors.MonitorStateStopped)
		require.Nil(t, svcErr)

		retrieved, _ := monitorService.GetMonitorByID(ctx, auth, monitor.ID.String())
		assert.Equal(t, monitors.MonitorStateStopped, retrieved.State)
	})

	t.Run("Fails with 400 for invalid monitor state", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		monitor := insertTestMonitor(t, ctx, project.ID)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, auth, monitor.ID, monitors.MonitorState("invalid_state"))
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

	t.Run("Fails with 404 for nonexistent monitor", func(t *testing.T) {
		ctx, monitorService, projectService, _, owner := setupMonitorIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "My Real Project"})
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  owner.Username,
		}

		svcErr := monitorService.UpdateMonitorStateByID(ctx, auth, uuid.New(), monitors.MonitorStateStopped)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

}
