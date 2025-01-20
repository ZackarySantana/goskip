package examples

import (
	"context"
	"errors"
	"fmt"

	skip "github.com/zackarysantana/goskip"
)

type StreamEvent[K comparable, V comparable] struct {
	Type skip.StreamType
	Data []skip.CollectionValue[K, V]
}

func ExpectData[K comparable, V comparable](ctx context.Context, streamClient skip.StreamClient, uuid string, expected []StreamEvent[K, V]) error {
	i := 0
	err := streamClient.Stream(ctx, uuid, skip.ReadStream(func(event skip.StreamType, data []skip.CollectionValue[K, V]) error {
		if i >= len(expected) {
			return fmt.Errorf("extra event received: '%v' with data '%v'", event, data)
		}
		if event != expected[i].Type {
			return fmt.Errorf("event type mismatch received '%s' expected '%s'", event, expected[i].Type)
		}
		for j, v := range data {
			expectedCollection := expected[i].Data
			if j >= len(expectedCollection) {
				return fmt.Errorf("got collections '%v' expected '%v'", data, expectedCollection)
			}
			expectedData := expected[i].Data[j]
			if expectedData.Key != v.Key {
				return fmt.Errorf("got key '%v' expected '%v'", v.Key, expectedData.Key)
			}
			if len(v.Values) != len(expectedData.Values) {
				return fmt.Errorf("got values '%v' expected '%v'", v.Values, expectedData.Values)
			}
			for k, expected := range expectedData.Values {
				got := v.Values[k]
				if expected != got {
					return fmt.Errorf("got value '%v' at index '%d' expected '%v'", got, k, expected)
				}
			}
		}
		i++
		return nil
	}))
	// If there error is context cancelled, ignore it. That's the test finishing.
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("streaming data failed during event index '%d': %w", i, err)
	}

	return nil
}
