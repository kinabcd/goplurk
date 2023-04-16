package goplurk

import (
	"golang.org/x/exp/maps"
)

type APIPolling struct {
	client *Client
}

// Return plurks newer than offset
// You must set offset in GetPlurksOptions or you will get nothing.
func (u *APIPolling) GetPlurks(optionSets ...Options) (*Plurks, error) {
	body := map[string]string{}
	for _, optionSet := range optionSets {
		maps.Copy(body, optionSet.Get())
	}
	plurks := &Plurks{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Polling/getPlurks", body, plurks); err != nil {
		return nil, err
	}
	return plurks, nil
}

func (u *APIPolling) GetUnreadCount() (*UnreadCount, error) {
	res := &UnreadCount{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Polling/getUnreadCount", map[string]string{}, res); err != nil {
		return nil, err
	}
	return res, nil
}
