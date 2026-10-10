package internal_containers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// --------------------------------------------------------- //
// Shared postgres container run, the options are extra container customizers

// hostPort is the port on the host, it binds the same numbered container port
func postgresRun(
	ctx context.Context,
	imageName string,
	containerName string,
	hostPort int,
	dbName string,
	dbUser string,
	dbPassword string,
	options ...testcontainers.ContainerCustomizer,
) (*postgres.PostgresContainer, error) {
	port, ok := network.PortFrom(uint16(hostPort), network.TCP)
	if !ok {
		return nil, fmt.Errorf("invalid host port %d", hostPort)
	}

	if len(dbName) <= 0 {
		dbName = "test_db_default"
		log.Printf("WARNING: dbName is empty, use default '%s'\n", dbName)
	}
	if len(dbUser) <= 0 {
		dbUser = "postgres"
		log.Printf("WARNING: dbUser is empty, use default '%s'\n", dbUser)
	}
	if len(dbPassword) <= 0 {
		dbPassword = "postgres"
		n := len(dbPassword)
		first := dbPassword[:1]
		last := dbPassword[n-1:]
		mid := dbPassword[n/2 : n/2+1]
		masked := first + "**" + mid + "**" + last
		log.Printf(
			"WARNING: dbPassword is empty, use default '%s'\n",
			string(fmt.Sprintf("%s", masked)),
		)
	}

	runOptions := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithName(containerName),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.PortBindings = network.PortMap{
				port: []network.PortBinding{{HostPort: strconv.Itoa(hostPort)}},
			}
		}),
		postgres.BasicWaitStrategies(),
	}
	runOptions = append(runOptions, options...)

	ctr, err := postgres.Run(ctx, imageName, runOptions...)
	if err != nil {
		return nil, fmt.Errorf("container for %s error: %w", imageName, err)
	}

	return ctr, nil
}

// --------------------------------------------------------- //

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

// Postgres container start on a fixed host port
func ContainerPostgresStart(
	ctx context.Context,
	imageName string,
	containerName string,
	hostPort int,
	dbName string,
	dbUser string,
	dbPassword string,
) (*postgres.PostgresContainer, error) {
	return postgresRun(
		ctx, imageName, containerName, hostPort,
		dbName, dbUser, dbPassword,
	)
}

// Postgres container start on a fixed host port, the sql files run in order on first start
func ContainerPostgresStartWithInit(
	ctx context.Context,
	imageName string,
	containerName string,
	hostPort int,
	dbName string,
	dbUser string,
	dbPassword string,
	sqlFiles ...string,
) (*postgres.PostgresContainer, error) {
	return postgresRun(
		ctx,
		imageName,
		containerName,
		hostPort,
		dbName,
		dbUser,
		dbPassword,
		postgres.WithOrderedInitScripts(sqlFiles...),
	)
}
