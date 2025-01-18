package skip

// SkipValue represents a value that will be used in a Skip collection.
type SkipValue interface {
	// Value returns the underlying value.
	Value() interface{}
}

type skipValueImpl struct {
	v interface{}
}

func (s *skipValueImpl) Value() interface{} {
	return s.v
}

func Values[T any](v ...T) []SkipValue {
	values := make([]SkipValue, len(v))
	for i, value := range v {
		values[i] = &skipValueImpl{v: value}
	}
	return values
}
