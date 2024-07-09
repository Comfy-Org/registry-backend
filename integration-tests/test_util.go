package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"registry-backend/drip"
	auth "registry-backend/server/middleware/authentication"
	"runtime"
	"strings"

	"registry-backend/ent"
	"registry-backend/ent/migrate"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

func createTestUser(ctx context.Context, client *ent.Client) *ent.User {
	return client.User.Create().
		SetID(uuid.New().String()).
		SetIsApproved(true).
		SetName("integration-test").
		SetEmail("integration-test@gmail.com").
		SaveX(ctx)
}

func createAdminUser(ctx context.Context, client *ent.Client) *ent.User {
	return client.User.Create().
		SetID(uuid.New().String()).
		SetIsApproved(true).
		SetIsAdmin(true).
		SetName("admin").
		SetEmail("admin@gmail.com").
		SaveX(ctx)
}

func decorateUserInContext(ctx context.Context, user *ent.User) context.Context {
	return context.WithValue(ctx, auth.UserContextKey, &auth.UserDetails{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func setupDB(t *testing.T, ctx context.Context) (client *ent.Client, cleanup func()) {
	// Define Postgres container request
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}
	println("Postgres container started")

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get the host: %s", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get the mapped port: %s", err)
	}
	waitPortOpen(t, host, port.Port(), time.Minute)
	databaseURL := fmt.Sprintf("postgres://postgres:password@%s:%s/postgres?sslmode=disable", host, port.Port())

	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}

	client, err = ent.Open("postgres", databaseURL)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed opening connection to postgres")
	}

	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true),
		migrate.WithDropColumn(true), migrate.WithDropIndex(true)); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed creating schema resources.")
		println("Failed to create schema")

	}
	println("Schema created")

	cleanup = func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Ctx(ctx).Error().Msgf("failed to terminate container: %s", err)
		}
	}
	return
}

func waitPortOpen(t *testing.T, host string, port string, timeout time.Duration) {
	tc := time.After(timeout)
	w, m := 500*time.Microsecond, 32*time.Second
	for {
		select {
		case <-tc:
			t.Errorf("timeout waiting to connect to '%s:%s'", host, port)
		default:
		}

		conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			t.Logf("error connecting to '%s:%s' : %s", host, port, err)
			if w < m {
				w *= 2
			}
			<-time.After(w)
			continue
		}

		conn.Close()
		return
	}
}

func withMiddleware[R any, S any](mw drip.StrictMiddlewareFunc, h func(ctx context.Context, req R) (res S, err error)) func(ctx context.Context, req R) (res S, err error) {
	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return h(ctx.Request().Context(), request.(R))
	}

	nameA := strings.Split(runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(), ".")
	nameA = strings.Split(nameA[len(nameA)-1], "-")
	opname := nameA[0]

	return func(ctx context.Context, req R) (res S, err error) {
		fakeReq := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		fakeRes := httptest.NewRecorder()
		fakeCtx := echo.New().NewContext(fakeReq, fakeRes)

		f := mw(handler, opname)
		r, err := f(fakeCtx, req)
		if r == nil {
			return *new(S), err
		}
		return r.(S), err
	}
}
