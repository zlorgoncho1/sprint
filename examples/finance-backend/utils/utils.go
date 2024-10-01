package utils

import "encoding/json"

func InterfaceToJSONObj(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var jsonObj map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonObj)
	if err != nil {
		return nil, err
	}

	return jsonObj, nil
}
