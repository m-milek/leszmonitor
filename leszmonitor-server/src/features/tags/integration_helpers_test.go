package tags

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/require"
)

// setupTagIntegrationTest initializes a temporary SQLite DB, sets up the tag service and puts a test user in the context.
func setupTagIntegrationTest(t *testing.T) (context.Context, *TagService, *db.Client, *auth.UserClaims) {
	ctx := context.Background()

	dsn := "file:" + filepath.Join(t.TempDir(), "testdb.sqlite") + "?_pragma=foreign_keys(1)"
	database, err := db.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(database.Close)

	tagService := NewTagService(TagServiceDeps{
		DB: database,
	})

	owner := &auth.UserClaims{
		Username: "integration_user",
	}
	ctx = auth.SetUserInContext(ctx, owner)

	return ctx, tagService, database, owner
}
