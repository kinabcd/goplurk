package goplurk

import (
	"strconv"
)

type APIProfile struct {
	client *Client
}

func (u *APIProfile) GetOwnProfile() (*Profile, error) {
	res := &Profile{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Profile/getOwnProfile", map[string]string{}, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
func (u *APIProfile) GetPublicProfile(userId int64) (*Profile, error) {
	body := map[string]string{
		"user_id": strconv.FormatInt(userId, 10),
	}
	res := &Profile{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Profile/getPublicProfile", body, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
