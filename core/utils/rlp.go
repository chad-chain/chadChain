package utils

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

// DecodeData decodes the data into the entity
// Entity should be a pointer to the struct to which the data should be decoded
func DecodeData(data []byte, entity interface{}) error {
	err := rlp.DecodeBytes(data, entity)
	if err != nil {
		return fmt.Errorf("error decoding: %v", err)
	}
	return nil
}

func EncodeData(data interface{}, isJson bool) ([]byte, error) {
	var err error
	if isJson {
		data, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error json marshaling: %v", err)
		}
	}
	encodedData, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding: %v", err)
	}
	return encodedData, nil
}
