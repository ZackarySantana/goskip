package skip

import (
	"encoding/json"
	"fmt"
)

// CollectionData represents a untyped set of collection values.
// This is used when posting new data to a collection.
type CollectionData struct {
	Key    interface{}
	Values []SkipValue
}

func (u *CollectionData) MarshalJSON() ([]byte, error) {
	values := make([]interface{}, len(u.Values))
	for i, value := range u.Values {
		values[i] = value.Value()
	}
	return json.Marshal([]interface{}{u.Key, values})
}

func (u *CollectionData) UnmarshalJSON(data []byte) error {
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

// CollectionValue is a typed set of collection values.
// This is used when reading data from a collection.
type CollectionValue[K any, V any] struct {
	Key    K
	Values []V
}

func (c *CollectionValue[K, V]) UnmarshalJSON(data []byte) error {
	var d []interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	if len(d) != 2 {
		return fmt.Errorf("invalid data length: expected 2, got %d", len(d))
	}
	var ok bool
	c.Key, ok = d[0].(K)
	if !ok {
		return fmt.Errorf("invalid data type '%T' for key, expected '%T'", d[0], c.Key)
	}
	values, ok := d[1].([]interface{})
	if !ok {
		return fmt.Errorf("invalid data type for '%T' values, expected []interface{}", d[1])
	}
	c.Values = make([]V, len(values))
	for i, value := range values {
		c.Values[i], ok = value.(V)
		if !ok {
			return fmt.Errorf("invalid data type '%T' for value, expected '%T'", value, c.Values[i])
		}
	}
	return nil
}

// ReadResourceSnapshot takes raw resource snapshot data and unmarshals it with the given types.
func ReadResourceSnapshot[K any, V any](data []byte, err error) ([]CollectionValue[K, V], error) {
	if err != nil {
		return nil, err
	}

	var values []CollectionValue[K, V]

	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	return values, nil
}

// ReadResourceKey takes raw resource key data and unmarshals it with the given type.
func ReadResourceKey[V any](data []byte, err error) ([]V, error) {
	var rawValues []interface{}
	if err := json.Unmarshal(data, &rawValues); err != nil {
		return nil, err
	}
	values := make([]V, len(rawValues))
	var ok bool
	for i, value := range rawValues {
		values[i], ok = value.(V)
		if !ok {
			return nil, fmt.Errorf("invalid data type '%T' for value, expected '%T'", value, values[i])
		}
	}
	return values, nil
}

func ReadStream[K any, V any](callback func(event StreamType, data []CollectionValue[K, V]) error) func(event StreamType, data []byte) error {
	return func(event StreamType, data []byte) error {
		snapshot, err := ReadResourceSnapshot[K, V](data, nil)
		if err != nil {
			return fmt.Errorf("reading resource snapshot: %w", err)
		}
		return callback(event, snapshot)
	}
}
