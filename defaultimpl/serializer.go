package defaultimpl

import "encoding/json"

type defaultSerializer struct{}

func (s *defaultSerializer) serialize(val interface{}) (string, error) {
	bytes, err := json.Marshal(val)
	return string(bytes), err
}

func (s *defaultSerializer) deserialize(val string, v interface{}) error {
	return json.Unmarshal([]byte(val), v)
}
