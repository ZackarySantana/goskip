package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controlClient, streamClient, shutdown, err := CreateClients(ctx)
	require.NoError(t, err)
	defer shutdown()
	require.NotNil(t, controlClient)
	require.NotNil(t, streamClient)

}
