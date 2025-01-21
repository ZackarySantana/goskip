package skipcontainer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SkipContainer represents the Skip container type.
type SkipContainer struct {
	testcontainers.Container
	host        string
	controlPort string
	streamPort  string
}

// Run creates an instance of the Skip container type.
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*SkipContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        img,
		ExposedPorts: []string{"8080", "8081"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Skip control service listening on port 8081"),
			wait.ForListeningPort("8081"),
			wait.ForLog("Skip streaming service listening on port 8080"),
			wait.ForListeningPort("8080"),
		),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, err
		}
	}
	skipContainer := &SkipContainer{}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if container != nil {
		skipContainer.Container = container
	}

	return skipContainer, err
}

// WithSkipFile sets the skip file to be used in the container.
func WithSkipFile(reader io.Reader) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Files = append(req.Files, testcontainers.ContainerFile{
			Reader:            reader,
			ContainerFilePath: "/app/skip.ts",
		})

		return nil
	}
}

type File struct {
	Reader            io.Reader
	ContainerFilePath string
}

// WithFiles sets the files to be used in the container. This is used
// for Skip services that have multiple files.
// This requires one of the files to be placed at /app/skip.ts
func WithFiles(files ...File) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		containerFiles := make([]testcontainers.ContainerFile, len(files))
		for i, file := range files {
			containerFiles[i] = testcontainers.ContainerFile{
				Reader:            file.Reader,
				ContainerFilePath: file.ContainerFilePath,
			}
		}

		req.Files = append(req.Files, containerFiles...)
		return nil
	}
}

// WithDirectory adds all files in the specified directory to the container. It will add all files in the subdirectories as well.
// This requires one of the files to be named skip.ts at the root
// of the directory. Files are automatically placed in the /app directory.
func WithDirectory(dir string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		info, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("failed to access directory %q: %w", dir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("provided path %q is not a directory", dir)
		}

		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking the path %q: %w", path, err)
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %q: %w", path, err)
			}
			defer file.Close()

			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, file)
			if err != nil {
				return fmt.Errorf("failed to copy file %q into buffer: %w", path, err)
			}

			relativePath, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("failed to determine relative path for %q: %w", path, err)
			}

			req.Files = append(req.Files, testcontainers.ContainerFile{
				Reader:            bytes.NewReader(buffer.Bytes()),
				ContainerFilePath: filepath.Join("/app", relativePath),
			})

			return nil
		})
		if err != nil {
			return fmt.Errorf("error processing directory %q: %w", dir, err)
		}

		return nil
	}
}

func (c *SkipContainer) saveHost(ctx context.Context) error {
	if c.host == "" {
		return nil
	}
	host, err := c.Container.Host(ctx)
	if err != nil {
		return err
	}
	c.host = host
	return nil
}

// GetControlURL returns the control service for the Skip container.
// The URL will be in the format http://<host>:<port>/v1
func (c *SkipContainer) GetControlURL() (string, error) {
	err := c.saveHost(context.Background())
	if err != nil {
		return "", err
	}
	if c.controlPort == "" {
		port, err := c.Container.MappedPort(context.Background(), "8081")
		if err != nil {
			return "", err
		}
		c.controlPort = port.Port()
	}
	return fmt.Sprintf("http://%s:%s/v1", c.host, c.controlPort), nil
}

// GetStreamURL returns the stream service for the Skip container.
// The URL will be in the format http://<host>:<port>/v1
func (c *SkipContainer) GetStreamURL() (string, error) {
	err := c.saveHost(context.Background())
	if err != nil {
		return "", err
	}
	if c.streamPort == "" {
		port, err := c.Container.MappedPort(context.Background(), "8080")
		if err != nil {
			return "", err
		}
		c.streamPort = port.Port()
	}
	return fmt.Sprintf("http://%s:%s/v1", c.host, c.streamPort), nil
}
