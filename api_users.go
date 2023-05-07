package goplurk

import (
	"strconv"
	"strings"
	"time"
)

type APIUsers struct {
	client *Client
}

func (u *APIUsers) Me() (*User, error) {
	res := &User{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Users/me", map[string]string{}, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (u *APIUsers) GetKarmaStats() (*KarmaStates, error) {
	res := &KarmaStates{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Users/getKarmaStats", map[string]string{}, res); err != nil {
		return nil, err
	}
	res.KarmaTrend = make([]KarmaTrendNode, len(res.KarmaTrendRaw))
	for i, n := range res.KarmaTrendRaw {
		splits := strings.Split(n, "-")
		timeInt, err := strconv.ParseInt(splits[0], 10, 64)
		if err != nil {
			return nil, err
		}
		karma, err := strconv.ParseFloat(splits[1], 64)
		if err != nil {
			return nil, err
		}

		res.KarmaTrend[i].Time = time.Unix(timeInt, 0)
		res.KarmaTrend[i].Karma = karma
	}

	return res, nil
}
