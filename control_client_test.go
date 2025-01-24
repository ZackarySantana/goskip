package skip_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
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
		if expected == nil {
			return actual == nil
		}
		if actual == nil {
			fmt.Println("actual is nil")
			return false
		}
		if actual.Method != expected.Method {
			fmt.Printf("actual.Method: %s, expected.Method: %s\n", actual.Method, expected.Method)
			return false
		}
		if actual.URL.String() != expected.URL.String() {
			fmt.Printf("actual.URL: %s, expected.URL: %s\n", actual.URL.String(), expected.URL.String())
			return false
		}
		for key, expectedValues := range expected.Header {
			if len(expectedValues) != len(actual.Header[key]) {
				fmt.Printf("actual.Header: %v, expected.Header: %v\n", actual.Header[key], expectedValues)
				return false
			}
			for _, expectedValue := range expectedValues {
				if !slices.Contains(actual.Header[key], expectedValue) {
					fmt.Printf("actual.Header: %v, expectedValue: %s\n", actual.Header[key], expectedValue)
					return false
				}
			}
		}
		if expected.Body == nil {
			if actual.Body != nil {
				fmt.Println("actual.Body is not nil, expected.Body is nil")
				return false
			}
		} else {
			if actual.Body == nil {
				fmt.Println("actual.Body is nil, expected.Body is not nil")
				return false
			}
			actualBody, err := io.ReadAll(actual.Body)
			if err != nil {
				fmt.Println("error reading actual body:", err)
				return false
			}
			expectedBody, err := io.ReadAll(expected.Body)
			if err != nil {
				fmt.Println("error reading expected body:", err)
				return false
			}
			if !bytes.Equal(actualBody, expectedBody) {
				fmt.Printf("actualBody: %s, expectedBody: %s\n", actualBody, expectedBody)
				return false
			}
		}
		if expected.Context() != nil {
			if actual.Context() == nil {
				fmt.Println("actual.Context is nil, expected.Context is not nil")
				return false
			}
		} else {
			if actual.Context() != expected.Context() {
				fmt.Printf("actual.Context: %v, expected.Context: %v\n", actual.Context(), expected.Context())
				return false
			}
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
			mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
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

	t.Run("GetResourceKey", func(t *testing.T) {
		t.Parallel()

		key := 123
		expectedResponse := []byte(`{"key":123,"value":"data"}`)

		t.Run("NoParams", func(t *testing.T) {
			t.Parallel()

			expectedRequestBody := map[string]interface{}{
				"key":    key,
				"params": nil,
			}
			expectedRequestBodyJSON, err := json.Marshal(expectedRequestBody)
			require.NoError(t, err)

			expectedRequest, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/snapshot/"+resourceName+"/lookup", bytes.NewReader(expectedRequestBodyJSON))
			require.NoError(t, err)
			expectedRequest.Header.Set("Content-Type", "application/json")

			mockHttpClient := mocks.NewMockhttpClient(t)
			mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
				Body:       io.NopCloser(bytes.NewBuffer(expectedResponse)),
				StatusCode: http.StatusOK,
			}, nil)

			client := skip.NewControlClient(baseURL, mockHttpClient)
			response, err := client.GetResourceKey(ctx, resourceName, key, nil)
			require.NoError(t, err)

			assert.Equal(t, expectedResponse, response)
		})

		t.Run("WithParams", func(t *testing.T) {
			t.Parallel()

			params := map[string]string{"filter": "active"}
			expectedRequestBody := map[string]interface{}{
				"key":    key,
				"params": params,
			}
			expectedRequestBodyJSON, err := json.Marshal(expectedRequestBody)
			require.NoError(t, err)

			expectedRequest, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/snapshot/"+resourceName+"/lookup", bytes.NewReader(expectedRequestBodyJSON))
			require.NoError(t, err)
			expectedRequest.Header.Set("Content-Type", "application/json")

			mockHttpClient := mocks.NewMockhttpClient(t)
			mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
				Body:       io.NopCloser(bytes.NewBuffer(expectedResponse)),
				StatusCode: http.StatusOK,
			}, nil)

			client := skip.NewControlClient(baseURL, mockHttpClient)
			response, err := client.GetResourceKey(ctx, resourceName, key, params)
			require.NoError(t, err)

			assert.Equal(t, expectedResponse, response)
		})
	})

	t.Run("CreateResourceInstance", func(t *testing.T) {
		t.Parallel()

		expectedUUID := "resource-uuid"
		expectedResponseBody := []byte(expectedUUID)

		t.Run("NoParams", func(t *testing.T) {
			t.Parallel()

			expectedRequest, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/streams/"+resourceName, nil)
			require.NoError(t, err)

			mockHttpClient := mocks.NewMockhttpClient(t)
			mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(expectedResponseBody)),
			}, nil)

			client := skip.NewControlClient(baseURL, mockHttpClient)
			uuid, err := client.CreateResourceInstance(ctx, resourceName, nil)
			require.NoError(t, err)

			assert.Equal(t, expectedUUID, uuid)
		})

		t.Run("WithParams", func(t *testing.T) {
			t.Parallel()

			params := map[string]interface{}{"config": "default"}
			expectedRequestBody, err := json.Marshal(params)
			require.NoError(t, err)

			expectedRequest, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/streams/"+resourceName, bytes.NewReader(expectedRequestBody))
			require.NoError(t, err)
			expectedRequest.Header.Set("Content-Type", "application/json")

			mockHttpClient := mocks.NewMockhttpClient(t)
			mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(expectedResponseBody)),
			}, nil)

			client := skip.NewControlClient(baseURL, mockHttpClient)
			uuid, err := client.CreateResourceInstance(ctx, resourceName, params)
			require.NoError(t, err)

			assert.Equal(t, expectedUUID, uuid)
		})
	})

	t.Run("UpdateInputCollection", func(t *testing.T) {
		t.Parallel()

		collection := "example-collection"
		updates := []skip.CollectionData{
			{Key: 1, Values: skip.Values[float64](10, 20)},
			{Key: 2, Values: skip.Values[float64](30, 40)},
		}
		expectedRequestBody, err := json.Marshal(updates)
		require.NoError(t, err)

		expectedRequest, err := http.NewRequestWithContext(ctx, "PATCH", baseURL+"/inputs/"+collection, bytes.NewReader(expectedRequestBody))
		require.NoError(t, err)
		expectedRequest.Header.Set("Content-Type", "application/json")

		mockHttpClient := mocks.NewMockhttpClient(t)
		mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		client := skip.NewControlClient(baseURL, mockHttpClient)
		err = client.UpdateInputCollection(ctx, collection, updates)
		require.NoError(t, err)
	})

	t.Run("DeleteResourceInstance", func(t *testing.T) {
		t.Parallel()

		uuid := "resource-uuid"
		expectedRequest, err := http.NewRequestWithContext(ctx, "DELETE", baseURL+"/streams/"+uuid, nil)
		require.NoError(t, err)

		mockHttpClient := mocks.NewMockhttpClient(t)
		mockHttpClient.EXPECT().Do(mock.MatchedBy(matchRequest(expectedRequest))).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		client := skip.NewControlClient(baseURL, mockHttpClient)
		err = client.DeleteResourceInstance(ctx, uuid)
		require.NoError(t, err)
	})
}
