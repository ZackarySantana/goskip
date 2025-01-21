package skipcontainer_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	skip "github.com/zackarysantana/goskip"
	"github.com/zackarysantana/goskip/skipcontainer"

	"github.com/testcontainers/testcontainers-go"
)

func TestSkipContainer(t *testing.T) {
	singleFile, err := os.Open("single/skip.ts")
	require.NoError(t, err)

	multiFile1, err := os.Open("multi/skip.ts")
	require.NoError(t, err)
	multiFile2, err := os.Open("multi/helpers.ts")
	require.NoError(t, err)

	mulitfileErr, err := os.Open("multi/skip.ts")
	require.NoError(t, err)

	for _, tc := range []struct {
		name      string
		shouldErr bool
		opts      []testcontainers.ContainerCustomizer
	}{
		{
			name:      "No Skip File",
			shouldErr: true,
			opts:      []testcontainers.ContainerCustomizer{},
		},
		{
			name: "With Skip File",
			opts: []testcontainers.ContainerCustomizer{
				skipcontainer.WithSkipFile(singleFile),
			},
		},
		{
			name: "With Skip Files",
			opts: []testcontainers.ContainerCustomizer{
				skipcontainer.WithFiles(
					testcontainers.ContainerFile{
						Reader:            multiFile1,
						ContainerFilePath: "/app/skip.ts",
					},
					testcontainers.ContainerFile{
						Reader:            multiFile2,
						ContainerFilePath: "/app/helpers.ts",
					},
				),
			},
		},
		{
			name:      "With Incomplete Skip Files",
			shouldErr: true,
			opts: []testcontainers.ContainerCustomizer{
				skipcontainer.WithFiles(
					testcontainers.ContainerFile{
						Reader:            mulitfileErr,
						ContainerFilePath: "/app/skip.ts",
					},
				),
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			skipContainer, err := skipcontainer.Run(ctx, "lidtop/goskip", tc.opts...)
			testcontainers.CleanupContainer(t, skipContainer)
			if tc.shouldErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			controlService, err := skipContainer.GetControlURL()
			require.NoError(t, err)

			controlClient := skip.NewControlClient(controlService)

			uuid, err := controlClient.CreateResourceInstance(ctx, "active_friends", 0)
			require.NoError(t, err)
			assert.NotEmpty(t, uuid)
		})
	}
}
