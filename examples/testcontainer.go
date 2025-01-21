package examples

import (
	"context"
	"fmt"
	"os"
	"time"

	skip "github.com/zackarysantana/goskip"
)

// StartSkipContainer creates a Skip container and returns a cleanup function to terminate it.
// The given path is the path to the skip.ts file to be used in the container.
// This is not suitable for production use but is useful for testing and development.
func StartSkipContainer(ctx context.Context, path string) (func(), error) {
	start := time.Now()
	fmt.Println("Starting skip container...")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	container, err := skip.Run(ctx, "lidtop/goskip", skip.WithSkipFile(file))
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

	controlURL, err := container.GetControlURL()
	if err != nil {
		return nil, fmt.Errorf("getting control url: %v", err)
	}

	streamURL, err := container.GetStreamURL()
	if err != nil {
		return nil, fmt.Errorf("getting stream url: %v", err)
	}

	err = os.Setenv("SKIP_CONTROL_URL", controlURL)
	if err != nil {
		return nil, fmt.Errorf("setting SKIP_CONTROL_URL: %v", err)
	}

	err = os.Setenv("SKIP_STREAM_URL", streamURL)
	if err != nil {
		return nil, fmt.Errorf("setting SKIP_STREAM_URL: %v", err)
	}

	fmt.Printf("Skip container started (%v).\n", time.Since(start).Round(time.Millisecond))
	return cleanup, nil
}
