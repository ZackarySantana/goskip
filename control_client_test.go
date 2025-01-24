package skip_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	skip "github.com/zackarysantana/goskip"
	"github.com/zackarysantana/goskip/internal/mocks"
)

func matchRequest(expected *http.Request) func(actual *http.Request) bool {
	return func(actual *http.Request) bool {
		if actual == nil {
			return false
		}
		if actual.Method != expected.Method {
			return false
		}
		if actual.URL.String() != expected.URL.String() {
			return false
		}
		if actual.Header.Get("Content-Type") != expected.Header.Get("Content-Type") {
			return false
		}
		if expected.Body == nil {
			return actual.Body == nil
		}
		actualBody, err := io.ReadAll(actual.Body)
		if err != nil {
			return false
		}
		expectedBody, err := io.ReadAll(expected.Body)
		if err != nil {
			return false
		}
		if !bytes.Equal(actualBody, expectedBody) {
			return false
		}
		if expected.Context() != nil {
			return actual.Context() == expected.Context()
		}
		if actual.Context() != nil {
			return false
		}
		return true
	}
}

func TestControlClient(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	baseURL := "http://localhost:8081/v1"
	resourceName := "rName"

	t.Run("GetResourceSnapshot", func(t *testing.T) {
		t.Parallel()

		expectedResponse := []skip.CollectionData{
			{
				Key:    1,
				Values: skip.Values[float64](20, 30),
			},
		}

		expectedResponseBody, err := json.Marshal(expectedResponse)
		require.NoError(t, err)

		t.Run("NoParams", func(t *testing.T) {
			t.Parallel()

			expectedRequest, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/snapshot/"+resourceName, nil)
			require.NoError(t, err)

			mockHttpClient := mocks.NewMockhttpClient(t)
			mockHttpClient.EXPECT().Do(expectedRequest).Return(&http.Response{
				Body:       io.NopCloser(bytes.NewBuffer(expectedResponseBody)),
				StatusCode: http.StatusOK,
			}, nil)

			client := skip.NewControlClient(baseURL, mockHttpClient)
			response, err := client.GetResourceSnapshot(ctx, resourceName, nil)
			require.NoError(t, err)

			assert.Equal(t, expectedResponseBody, response)
		})

		t.Run("WithParams", func(t *testing.T) {
			t.Parallel()

			expectedRequest, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/snapshot/"+resourceName, strings.NewReader("0"))
			require.NoError(t, err)
			expectedRequest.Header.Set("Content-Type", "application/json")

			mockHttpClient := mocks.NewMockhttpClient(t)
			mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
				Body:       io.NopCloser(bytes.NewBuffer(expectedResponseBody)),
				StatusCode: http.StatusOK,
			}, nil)

			client := skip.NewControlClient(baseURL, mockHttpClient)
			response, err := client.GetResourceSnapshot(ctx, resourceName, 0)
			require.NoError(t, err)

			assert.Equal(t, expectedResponseBody, response)
		})
	})
}
