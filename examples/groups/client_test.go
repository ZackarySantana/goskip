package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	skip "github.com/zackarysantana/goskip"
	"github.com/zackarysantana/goskip/examples"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := examples.StartSkipContainer(ctx, "skip.ts")
	require.NoError(t, err)
	defer shutdown()

	controlClient := skip.NewControlClient(os.Getenv("SKIP_CONTROL_URL"), nil)
	streamClient := skip.NewStreamClient(os.Getenv("SKIP_STREAM_URL"), nil)

	uuid, err := controlClient.CreateResourceInstance(ctx, "active_friends", 0)
	require.NoError(t, err)
	require.NotEmpty(t, uuid)

	var streamErr error
	streamingComplete := false
	streamCtx, streamCancel := context.WithCancel(ctx)
	defer streamCancel()
	expectedStreamData := []examples.StreamEvent[float64, float64]{
		{
			Type: skip.InitStreamType,
			Data: []skip.CollectionValue[float64, float64]{
				{
					Key:    1001,
					Values: []float64{1},
				},
			},
		},
		{
			Type: skip.UpdateStreamType,
			Data: []skip.CollectionValue[float64, float64]{
				{
					Key:    1001,
					Values: []float64{1, 2},
				},
				{
					Key:    1002,
					Values: []float64{2},
				},
			},
		},
		{
			Type: skip.UpdateStreamType,
			Data: []skip.CollectionValue[float64, float64]{
				{
					Key:    1001,
					Values: []float64{2},
				},
			},
		},
		{
			Type: skip.UpdateStreamType,
			Data: []skip.CollectionValue[float64, float64]{
				{
					Key:    1001,
					Values: []float64{2, 3},
				},
				{
					Key:    1002,
					Values: []float64{2},
				},
			},
		},
		{
			Type: skip.UpdateStreamType,
			Data: []skip.CollectionValue[float64, float64]{
				{
					Key:    1002,
					Values: []float64{3},
				},
			},
		},
	}
	go func() {
		streamErr = examples.ExpectData(streamCtx, streamClient, uuid, expectedStreamData)
		streamingComplete = true
	}()
	examples.WaitForIO()

	examples.UpdateInputCollection(ctx, t, controlClient, "users", usersUpdate1)
	examples.UpdateInputCollection(ctx, t, controlClient, "users", usersUpdate2)
	examples.UpdateInputCollection(ctx, t, controlClient, "users", usersUpdate3)
	examples.UpdateInputCollection(ctx, t, controlClient, "groups", groupsUpdate1)

	streamCancel()
	examples.WaitForIO()
	require.NoError(t, streamErr)
	require.True(t, streamingComplete)
}
