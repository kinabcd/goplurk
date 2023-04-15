package goplurk

import (
	"encoding/json"
)

type APIUsers struct {
	client *Client
}

func (u *APIUsers) Me() (*User, error) {
	resBytes, err := u.client.Engine.CallAPI("/APP/Users/me", map[string]string{})
	if err != nil {
		return nil, err
	}
	var userDate = &User{}
	if err := json.Unmarshal(resBytes, userDate); err != nil {
		return nil, err
	}
	return userDate, nil
}
