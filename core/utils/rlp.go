package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeReceived(base64Str any, isJson bool) interface{} {

	decoded, err := base64.StdEncoding.DecodeString(base64Str.(string)) // Decode the base64 string
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return nil
	}

	var data string

	if err := rlp.DecodeBytes(decoded, &data); err != nil {
		fmt.Println("Error decoding :", err)
		fmt.Printf("Decoded bytes: %v\n", decoded)
		return nil
	}

	if isJson {
		var finalData interface{}
		json.Unmarshal([]byte(data), &finalData)
		fmt.Println("Received decoded:", finalData)
		return finalData
	}

	fmt.Println("Received decoded:", data)
	return data
}

func EncodeData(data interface{}, isJson bool) []byte {
	var err error
	if isJson {
		data, err = json.Marshal(data)
		if err != nil {
			fmt.Println("Error json marshaling :", err)
			return nil
		}
	}
	encodedData, err := rlp.EncodeToBytes(data)
	if err != nil {
		fmt.Println("Error encoding :", err)
		return nil
	}
	return encodedData
}
