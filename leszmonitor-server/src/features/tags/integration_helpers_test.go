package tags_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
)

func setupTagIntegrationTest(
	t *testing.T,
) (context.Context, *tags.TagService, *db.Client, *users.User) {
	ctx, _, user := testsupport.Setup(t)

	tagService := tags.NewTagService(tags.TagServiceDeps{
		DB: db.Get(),
	})

	return ctx, tagService, db.Get(), user
}
