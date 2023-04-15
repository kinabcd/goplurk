package goplurk

import (
	"encoding/json"
	"strconv"
)

type APIProfile struct {
	client *Client
}

func (u *APIProfile) GetOwnProfile() (*Profile, error) {
	resBytes, err := u.client.Engine.CallAPI("/APP/Profile/getOwnProfile", map[string]string{})
	if err != nil {
		return nil, err
	}
	var userDate = &Profile{}
	if err := json.Unmarshal(resBytes, userDate); err != nil {
		return nil, err
	}
	return userDate, nil
}
func (u *APIProfile) GetPublicProfile(userId int64) (*Profile, error) {
	resBytes, err := u.client.Engine.CallAPI("/APP/Profile/getPublicProfile", map[string]string{
		"user_id": strconv.FormatInt(userId, 10),
	})
	if err != nil {
		return nil, err
	}
	var userDate = &Profile{}
	if err := json.Unmarshal(resBytes, userDate); err != nil {
		return nil, err
	}
	return userDate, nil
}
