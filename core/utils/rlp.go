package utils

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeReceived(endoded any, isJson bool) (interface{}, error) {

	var data string

	if err := rlp.DecodeBytes(endoded.([]byte), &data); err != nil {
		return nil, fmt.Errorf("error decoding RLP bytes: %v", err)
	}

	if isJson {
		var finalData interface{}
		if err := json.Unmarshal([]byte(data), &finalData); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
		}
		return finalData, nil
	}

	return data, nil
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
