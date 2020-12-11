package utils

import (
	"encoding/json"
)

// UnmarshalInto given JSON like data
func UnmarshalInto(data interface{}, target interface{}) error {
	return json.Unmarshal([]byte(data.(string)), target)
}
