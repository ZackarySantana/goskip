package skip

import (
	"context"
	"fmt"
	"io"
)

// ControlClient defines access to Skip's Control API.
type ControlClient interface {
	// GetSnapshot retrieves a snapshot of the entire resource.
	// Corresponds to the POST /v1/snapshot/:resource endpoint.
	GetResourceSnapshot(ctx context.Context, resource string, params interface{}) ([]byte, error)

	// GetResourceKey retrieves the data associated with a specific key in a resource.
	// Corresponds to the POST /v1/snapshot/:resource/lookup endpoint.
	GetResourceKey(ctx context.Context, resource string, key interface{}, params interface{}) ([]byte, error)

	// UpdateInputCollection updates a collection of key-value pairs in the specified input collection.
	// Corresponds to the PATCH /v1/inputs/:collection endpoint.
	UpdateInputCollection(ctx context.Context, collection string, updates []CollectionData) error

	// CreateResourceInstance creates a new resource instance and returns its UUID.
	// Corresponds to the POST /v1/streams/:resource endpoint.
	CreateResourceInstance(ctx context.Context, resource string, params interface{}) (string, error)

	// DeleteResourceInstance deletes a resource instance by its UUID.
	// Corresponds to the DELETE /v1/streams/:uuid endpoint.
	DeleteResourceInstance(ctx context.Context, uuid string) error
}

type controlClientImpl struct {
	baseURL string
}

// NewControlClient creates a new instance of ControlClient.
func NewControlClient(baseURL string) ControlClient {
	return &controlClientImpl{baseURL: baseURL}
}

func (c *controlClientImpl) GetResourceSnapshot(ctx context.Context, resource string, params interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/snapshot/%s", c.baseURL, resource)
	resp, err := sendRequest(ctx, "POST", url, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource snapshot: %w", err)
	}
	defer resp.Body.Close()

	if isSuccessStatus(resp.StatusCode) {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *controlClientImpl) GetResourceKey(ctx context.Context, resource string, key interface{}, params interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/snapshot/%s/lookup", c.baseURL, resource)
	body := map[string]interface{}{
		"key":    key,
		"params": params,
	}
	resp, err := sendRequest(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource key: %w", err)
	}
	defer resp.Body.Close()

	if isSuccessStatus(resp.StatusCode) {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *controlClientImpl) UpdateInputCollection(ctx context.Context, collection string, updates []CollectionData) error {
	url := fmt.Sprintf("%s/inputs/%s", c.baseURL, collection)
	resp, err := sendRequest(ctx, "PATCH", url, updates)
	if err != nil {
		return fmt.Errorf("failed to update input collection: %w", err)
	}
	defer resp.Body.Close()

	if isSuccessStatus(resp.StatusCode) {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *controlClientImpl) CreateResourceInstance(ctx context.Context, resource string, params interface{}) (string, error) {
	url := fmt.Sprintf("%s/streams/%s", c.baseURL, resource)
	resp, err := sendRequest(ctx, "POST", url, params)
	if err != nil {
		return "", fmt.Errorf("failed to create resource: %w", err)
	}
	defer resp.Body.Close()

	if isSuccessStatus(resp.StatusCode) {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	uuid, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(uuid), nil
}

func (c *controlClientImpl) DeleteResourceInstance(ctx context.Context, uuid string) error {
	url := fmt.Sprintf("%s/streams/%s", c.baseURL, uuid)
	resp, err := sendRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	defer resp.Body.Close()

	if isSuccessStatus(resp.StatusCode) {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
