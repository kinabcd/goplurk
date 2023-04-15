package goplurk

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type APIResponses struct {
	client *Client
}

func (u *APIResponses) Get(plurkId int64, fromResoponse int64, count int64) (*Responses, error) {
	var body = map[string]string{}
	body["plurk_id"] = strconv.FormatInt(plurkId, 10)
	body["from_response"] = strconv.FormatInt(fromResoponse, 10)
	body["count"] = strconv.FormatInt(count, 10)
	res, err := u.client.Engine.CallAPI("/APP/Responses/get", body)
	if err != nil {
		return nil, err
	}
	responses := Responses{}
	if err := json.Unmarshal(res, &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal responses: %v, %s", err, string(res))
	}

	return &responses, nil
}
func (u *APIResponses) ResponseAdd(plurkId int64, qualifier string, content string) (*Response, error) {
	if qualifier == "" {
		qualifier = ":"
	}
	if content == "" {
		return nil, fmt.Errorf("content can not be empty")
	}
	var body = map[string]string{}
	body["plurk_id"] = strconv.FormatInt(plurkId, 10)
	body["qualifier"] = qualifier
	body["content"] = content
	res, err := u.client.Engine.CallAPI("/APP/Responses/responseAdd", body)
	if err != nil {
		return nil, err
	}
	response := Response{}
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v, %s", err, string(res))
	}
	return &response, nil
}

func (u *APIResponses) ResponseDelete(responseId int64, plurkId int64) error {
	_, err := u.client.Engine.CallAPI("/APP/Responses/responseDelete", map[string]string{
		"response_id": strconv.FormatInt(responseId, 10),
		"plurk_id":    strconv.FormatInt(plurkId, 10),
	})
	if err != nil {
		return fmt.Errorf("failed to delete response: %v", err)
	}
	return nil

}
