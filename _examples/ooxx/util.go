package ooxx

import (
	"encoding/json"
)

func UnmarshalJSON(data []byte) ([]OOXXModel, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(data, dat); err != nil {
		return nil, err
	}
	return nil, nil
}
