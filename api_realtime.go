package goplurk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type comet struct {
	NewOffset int64             `json:"new_offset"`
	Data      []json.RawMessage `json:"data"`
}

type RealtimeLogEvent struct {
	Log *string
	Err error
}

func newRealtimeLogEvent(s string) *RealtimeLogEvent {
	return &RealtimeLogEvent{Log: &s}
}

type APIRealtime struct {
	client *Client
}

func (u *APIRealtime) GetUserChannel() (*UserChannel, error) {
	res := &UserChannel{}
	// {"comet_server": "https://comet03.plurk.com/comet/1235515351741/?channel=generic-4-f733d8522327edf87b4d1651e6395a6cca0807a0", "channel_name": "generic-4-f733d8522327edf87b4d1651e6395a6cca0807a0"}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Realtime/getUserChannel", map[string]string{}, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
func (u *APIRealtime) Listen(ctx context.Context, callback func(interface{})) {
	var channel *UserChannel = nil
	var offset int64 = 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if channel == nil {
				newChannel, err := u.GetUserChannel()
				if err != nil {
					callback(&RealtimeLogEvent{Err: err})
					time.Sleep(5 * time.Second)
					continue
				}
				channel = newChannel
				offset = 0
				callback(newRealtimeLogEvent(fmt.Sprintf("newChannel %s", newChannel.CometServer)))

			}
			serverUrl, err := url.Parse(channel.CometServer)
			if err != nil {
				callback(&RealtimeLogEvent{Err: err})
				time.Sleep(5 * time.Second)
				continue
			}
			q := serverUrl.Query()
			q.Set("offset", strconv.FormatInt(offset, 10))
			serverUrl.RawQuery = q.Encode()
			body, err := getUrl(serverUrl)
			if err != nil {
				callback(&RealtimeLogEvent{Err: err})
				time.Sleep(5 * time.Second)
				continue
			}

			newComet, err := resolveComet(body)
			if err != nil {
				callback(&RealtimeLogEvent{Err: err})
				time.Sleep(5 * time.Second)
				continue
			}
			if offset != newComet.NewOffset {
				callback(newRealtimeLogEvent(fmt.Sprintf("offset %d", newComet.NewOffset)))
			}
			offset = newComet.NewOffset
			if offset == -3 {
				// Your offset is wrong and you need to resync your data
				channel = nil
				continue
			}
			for _, rawEvent := range newComet.Data {
				if event, err := resolveEvent(rawEvent); err != nil {
					callback(&RealtimeLogEvent{Err: err})
				} else {
					callback(event)
				}
			}
		}
	}
}

func getUrl(_url *url.URL) ([]byte, error) {
	res, err := http.Get(_url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get: %s, %v", _url.String(), err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	return body, nil
}

func resolveComet(bytes []byte) (*comet, error) {
	bytes = []byte(strings.TrimSuffix(strings.TrimPrefix(string(bytes), "CometChannel.scriptCallback("), ");"))
	res := &comet{}
	if err := json.Unmarshal(bytes, res); err != nil {
		return nil, err
	}
	return res, nil
}

func resolveEvent(bytes json.RawMessage) (interface{}, error) {
	pass1 := &struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, pass1); err != nil {
		return nil, err
	}
	switch pass1.Type {
	case "new_response":
		res := &NewResponseEvent{}
		if err := json.Unmarshal(bytes, res); err != nil {
			return nil, err
		}
		return res, nil

	case "new_plurk":
		res := &NewPlurkEvent{}
		if err := json.Unmarshal(bytes, res); err != nil {
			return nil, err
		}
		return res, nil
	case "update_notification":
		res := &UpdateNotificationEvent{}
		if err := json.Unmarshal(bytes, res); err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("not handled event: %s", string(bytes))
	}
}
