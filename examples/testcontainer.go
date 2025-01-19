package examples

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartSkipContainer(ctx context.Context, path string) (func(), error) {
	start := time.Now()
	fmt.Println("Starting skip container...")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	req := testcontainers.ContainerRequest{
		Image:        "lidtop/goskip",
		ExposedPorts: []string{"8080", "8081"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Skip control service listening on port 8081"),
			wait.ForListeningPort("8081"),
			wait.ForLog("Skip streaming service listening on port 8080"),
			wait.ForListeningPort("8080"),
		),
		Env: map[string]string{},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            file,
				ContainerFilePath: "/app/skip.ts",
			},
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("starting container: %v", err)
	}
	cleanup := func() {
		start := time.Now()
		fmt.Println("Cleaning up skip container...")
		err := container.Terminate(ctx)
		if err != nil {
			fmt.Printf("Error cleaning up container: %v\n", err)
		}
		fmt.Printf("Skip container cleaned up (%v).\n", time.Since(start).Round(time.Millisecond))
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting container host: %v", err)
	}

	controlPort, err := container.MappedPort(ctx, "8081")
	if err != nil {
		return nil, fmt.Errorf("getting control service port: %v", err)
	}

	streamPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, fmt.Errorf("getting stream service port: %v", err)
	}

	err = os.Setenv("SKIP_CONTROL_URL", fmt.Sprintf("http://%s:%s/v1", host, controlPort.Port()))
	if err != nil {
		return nil, fmt.Errorf("setting SKIP_CONTROL_URL: %v", err)
	}

	err = os.Setenv("SKIP_STREAM_URL", fmt.Sprintf("http://%s:%s/v1", host, streamPort.Port()))
	if err != nil {
		return nil, fmt.Errorf("setting SKIP_STREAM_URL: %v", err)
	}

	fmt.Printf("Skip container started (%v).\n", time.Since(start).Round(time.Millisecond))
	return cleanup, nil
}
