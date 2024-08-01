package goplurk

import (
	"strconv"
)

type APIProfile struct {
	client *Client
}

// Returns data that's private for the current user. This can be used to construct a profile and render a timeline of the latest plurks.
func (u *APIProfile) GetOwnProfile() (*Profile, error) {
	body := map[string]string{
		"include_plurks": "false",
	}
	res := &Profile{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Profile/getOwnProfile", body, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Fetches public information such as a user's public plurks and basic information. Fetches also if the current user is following the user, are friends with or is a fan.
// userIdOrNickName must be int64 for user_id or string for nick_name
func (u *APIProfile) GetPublicProfile(userIdOrNickName any) (*Profile, error) {
	body := map[string]string{
		"include_plurks": "false",
	}
	userId := func(num any) (int64, bool) {
		switch v := num.(type) {
		case int64:
			return v, true
		case int32:
			return int64(v), true
		case int16:
			return int64(v), true
		case int8:
			return int64(v), true
		case int:
			return int64(v), true
		case uint64:
			return int64(v), true
		case uint32:
			return int64(v), true
		case uint16:
			return int64(v), true
		case uint8:
			return int64(v), true
		case uint:
			return int64(v), true
		}
		return 0, false
	}
	if userId, ok := userId(userIdOrNickName); ok {
		body["user_id"] = strconv.FormatInt(userId, 10)
	} else if nickName, ok := userIdOrNickName.(string); ok {
		body["nick_name"] = nickName
	} else {
		return nil, ErrInvalidParameters
	}

	res := &Profile{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Profile/getPublicProfile", body, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
