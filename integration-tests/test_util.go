package integration

import (
	"context"
	"fmt"
	"net"
	"registry-backend/ent"
	"registry-backend/ent/migrate"
	auth "registry-backend/server/middleware"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func createTestUser(ctx context.Context, client *ent.Client) *ent.User {
	return client.User.Create().
		SetID(uuid.New().String()).
		SetIsApproved(true).
		SetName("integration-test").
		SetEmail("integration-test@gmail.com").
		SaveX(ctx)
}

func decorateUserInContext(ctx context.Context, user *ent.User) context.Context {
	return context.WithValue(ctx, auth.UserContextKey, &auth.UserDetails{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func setupDB(t *testing.T, ctx context.Context) (*ent.Client, *postgres.PostgresContainer) {
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

	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed opening connection to postgres")
	}

	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true),
		migrate.WithDropColumn(true), migrate.WithDropIndex(true)); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed creating schema resources.")
	}
	return client, postgresContainer
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
