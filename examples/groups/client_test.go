package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	skip "github.com/zackarysantana/goskip"
	"github.com/zackarysantana/goskip/examples"
)

func waitForIO() {
	time.Sleep(100 * time.Millisecond)
}

func TestClient(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := examples.StartSkipContainer(ctx, "examples/groups/skip.ts")
	require.NoError(t, err)
	defer shutdown()

	controlClient := skip.NewControlClient(os.Getenv("SKIP_CONTROL_URL"))
	streamClient := skip.NewStreamClient(os.Getenv("SKIP_STREAM_URL"))

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
	waitForIO()

	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionData{
		{
			Key: 2,
			Values: skip.Values(
				UsersValue{
					Name:    "Carol",
					Active:  true,
					Friends: []int{0, 1},
				},
			),
		},
	})
	require.NoError(t, err)
	waitForIO()

	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionData{
		{
			Key: 1,
			Values: skip.Values(
				UsersValue{
					Name:    "Alice",
					Active:  false,
					Friends: []int{0, 2},
				},
			),
		},
	})
	require.NoError(t, err)
	waitForIO()

	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionData{
		{
			Key: 0,
			Values: skip.Values(
				UsersValue{
					Name:    "Bob",
					Active:  true,
					Friends: []int{1, 2, 3},
				},
			),
		},
	})
	require.NoError(t, err)
	waitForIO()

	err = controlClient.UpdateInputCollection(ctx, "groups", []skip.CollectionData{
		{
			Key: 1002,
			Values: skip.Values(
				GroupsValue{
					Name:    "Group 2",
					Members: []int{0, 3},
				},
			),
		},
	})
	require.NoError(t, err)
	waitForIO()

	streamCancel()
	waitForIO()
	require.NoError(t, streamErr)
	require.True(t, streamingComplete)
}
