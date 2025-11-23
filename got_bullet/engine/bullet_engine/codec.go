package bullet_engine

import "encoding/json"

// The JSON encoder/decoder interface
type Codec[T any] interface {
	Encode(value T) (string, error)
	Decode(data string, value *T) error
}

// A simple JSON implementation
type JSONCodec[T any] struct{}

func (j *JSONCodec[T]) Encode(value T) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (j *JSONCodec[T]) Decode(data string, value *T) error {
	return json.Unmarshal([]byte(data), value)
}
