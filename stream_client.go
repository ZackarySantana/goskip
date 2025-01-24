package skip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tmaxmax/go-sse"
)

type StreamType string

const (
	InitStreamType   StreamType = "init"
	UpdateStreamType StreamType = "update"
)

// StreamClient defines access to Skip's Stream API.
type StreamClient interface {
	// Stream is a live data stream for a resource instance represented by the UUID.
	// Corresponds to the GET /v1/streams/:uuid endpoint.
	Stream(ctx context.Context, uuid string, callback func(event StreamType, data []byte) error) error
}

type streamClientImpl struct {
	baseURL    string
	httpClient httpClient
}

// NewStreamClient creates a new instance of StreamClient.
func NewStreamClient(baseURL string, httpClient httpClient) StreamClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &streamClientImpl{baseURL: baseURL, httpClient: httpClient}
}

func (s *streamClientImpl) Stream(ctx context.Context, uuid string, callback func(event StreamType, data []byte) error) error {
	url := fmt.Sprintf("%s/streams/%s", s.baseURL, uuid)
	resp, err := sendRequest(ctx, s.httpClient, "GET", url, nil)
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

		err = callback(StreamType(ev.Type), []byte(ev.Data))

		return err == nil
	})

	return err
}
