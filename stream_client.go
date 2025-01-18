package skip

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmaxmax/go-sse"
)

type StreamType string

const (
	InitStreamType   StreamType = "init"
	UpdateStreamType StreamType = "update"
)

type StreamClient interface {
	// StreamData is a live data stream for a resource instance represented by the UUID.
	// Corresponds to the GET /v1/streams/:uuid endpoint.
	StreamData(ctx context.Context, uuid string, callback func(event StreamType, data []CollectionUpdate)) error
}

type streamingClientImpl struct {
	baseURL string
}

func NewStreamingClient(baseURL string) StreamClient {
	return &streamingClientImpl{baseURL: baseURL}
}

func (s *streamingClientImpl) StreamData(ctx context.Context, uuid string, callback func(event StreamType, data []CollectionUpdate)) error {
	url := fmt.Sprintf("%s/streams/%s", s.baseURL, uuid)
	resp, err := sendRequest(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("sending stream data request: %w", err)
	}
	defer resp.Body.Close()

	sse.Read(resp.Body, nil)(func(ev sse.Event, readErr error) bool {
		if readErr != nil {
			err = fmt.Errorf("reading stream data: %w", readErr)
			return false
		}
		if ev.Type != "init" && ev.Type != "update" {
			err = fmt.Errorf("unexpected event type: %s", ev.Type)
			return false
		}

		var data []CollectionUpdate
		err = json.Unmarshal([]byte(ev.Data), &data)
		if err != nil {
			err = fmt.Errorf("unmarshalling stream data: %w", err)
			return false
		}

		callback(StreamType(ev.Type), data)

		return true
	})

	return err
}
