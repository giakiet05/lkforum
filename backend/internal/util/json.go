package util

import "encoding/json"

func DecodeJson[T any](data interface{}, out *T) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, out)
}
