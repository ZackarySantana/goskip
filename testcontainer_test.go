package skip_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	skip "github.com/zackarysantana/goskip"

	"github.com/testcontainers/testcontainers-go"
)

func TestSkipContainer(t *testing.T) {
	skipFile, err := os.Open("testcontainer.ts")
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
				skip.WithSkipFile(skipFile),
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			skipContainer, err := skip.Run(ctx, "lidtop/goskip", tc.opts...)
			testcontainers.CleanupContainer(t, skipContainer)
			if tc.shouldErr {
				require.Error(t, err)
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
