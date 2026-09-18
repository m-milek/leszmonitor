package users_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/internal/testsupport"
)

func setupIntegrationTest(t *testing.T) (context.Context, *users.UserService, *users.User) {
	return testsupport.Setup(t)
}
