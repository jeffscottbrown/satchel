//go:build !production

package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ConfigureRepositoryForTest sets the global Repository variable to the
// provided test repository. It also ensures that the original repository
// is restored after the test completes. This is useful for isolating
// tests and ensuring that they do not interfere with each other.
func ConfigureRepositoryForTest(t *testing.T, testRepo EmployeeRepository) {
	originalRepository := employeeRepository
	t.Cleanup(func() {
		employeeRepository = originalRepository
	})
	employeeRepository = testRepo
}

func RunTestsWithTestContainer(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	var err error
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start postgres container: %v\n", err)
		os.Exit(1)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get container host: %v\n", err)
		os.Exit(1)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get container port: %v\n", err)
		os.Exit(1)
	}

	os.Setenv("SATCHEL_DB_HOST", host)
	os.Setenv("SATCHEL_DB_PORT", port.Port())
	os.Setenv("SATCHEL_DB_USER", "testuser")
	os.Setenv("SATCHEL_DB_PASSWORD", "testpass")
	os.Setenv("SATCHEL_DB_NAME", "testdb")

	InitializeDatabase()

	code := m.Run()

	if postgresContainer != nil {
		_ = postgresContainer.Terminate(ctx)
	}

	os.Exit(code)
}
