package internal_containers

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Default container start
func ContainerStart(ctx context.Context, imageName string, containerName string) (*testcontainers.DockerContainer, error) {
	ctr, err := testcontainers.Run(
		ctx,
		containerName,
		testcontainers.WithName(imageName),
	)
	if err != nil {
		return nil, fmt.Errorf("container for %s error: %w", imageName, err)
	}

	return ctr, nil
}

// Postgres container start
func ContainerPostgresStart(ctx context.Context, imageName string, containerName string) (*postgres.PostgresContainer, error) {
	ctr, err := postgres.Run(
		ctx,
		imageName,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithName(containerName),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("container for %s error: %w", imageName, err)
	}

	return ctr, nil
}
