package goplurk

import (
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
	responses := &Responses{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Responses/get", body, responses); err != nil {
		return nil, err
	} else {
		return responses, nil
	}
}
func (u *APIResponses) ResponseAdd(plurkId int64, qualifier string, content string) (*Response, error) {
	var body = map[string]string{}
	body["plurk_id"] = strconv.FormatInt(plurkId, 10)
	body["qualifier"] = qualifier
	body["content"] = content
	response := &Response{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Responses/responseAdd", body, response); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (u *APIResponses) ResponseDelete(responseId int64, plurkId int64) error {
	_, err := u.client.Engine.CallAPI("/APP/Responses/responseDelete", map[string]string{
		"response_id": strconv.FormatInt(responseId, 10),
		"plurk_id":    strconv.FormatInt(plurkId, 10),
	})
	if err != nil {
		return fmt.Errorf("failed to delete response: %v", err)
	} else {
		return nil
	}
}
