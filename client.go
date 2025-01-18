package skip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CollectionUpdate represents an update to a key-value pair in a collection.
type CollectionUpdate struct {
	Key    interface{}
	Values []SkipValue
}

func (u *CollectionUpdate) MarshalJSON() ([]byte, error) {
	values := make([]interface{}, len(u.Values))
	for i, value := range u.Values {
		values[i] = value.Value()
	}
	return json.Marshal([]interface{}{u.Key, values})
}

func (u *CollectionUpdate) UnmarshalJSON(data []byte) error {
	var d []interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	if len(d) != 2 {
		return fmt.Errorf("invalid data length: expected 2, got %d", len(d))
	}
	u.Key = d[0]
	values, ok := d[1].([]interface{})
	if !ok {
		return fmt.Errorf("invalid data type for values")
	}
	u.Values = make([]SkipValue, len(values))
	for i, value := range values {
		u.Values[i] = &skipValueImpl{v: value}
	}
	return nil
}

func sendRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var requestBody *strings.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		requestBody = strings.NewReader(string(data))
	} else {
		requestBody = strings.NewReader("")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	return client.Do(req)
}

func isSuccessStatus(statusCode int) bool {
	return statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices
}
