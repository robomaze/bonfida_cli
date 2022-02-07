package common

import (
	"encoding/json"
)

type BonfidaResponse struct {
	Success bool   `json:"success"`
	Data    []byte `json:"data"`
}

func (r *BonfidaResponse) UnmarshalJSON(data []byte) error {
	resMap := struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}{}

	err := json.Unmarshal(data, &resMap)
	if err != nil {
		return err
	}

	r.Success = resMap.Success

	resData, err := json.Marshal(resMap.Data)
	if err != nil {
		return err
	}

	r.Data = resData

	return nil
}
