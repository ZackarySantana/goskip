package skip

import "fmt"

// skipValue represents a value that will be used in a Skip collection.
type skipValue struct {
	v interface{}
}

func (s *skipValue) Value() interface{} {
	return s.v
}

func (s *skipValue) String() string {
	if s == nil || s.v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", s.v)
}

func Values[T any](v ...T) []skipValue {
	values := make([]skipValue, len(v))
	for i, value := range v {
		values[i] = skipValue{v: value}
	}
	return values
}
