package skip

import (
	"encoding/json"
	"fmt"
)

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

type CollectionValue[K any, V any] struct {
	Key    K
	Values []V
}

func Read[K any, V any](data []byte) ([]CollectionValue[K, V], error) {
	var values []CollectionValue[K, V]

	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	return values, nil
}
