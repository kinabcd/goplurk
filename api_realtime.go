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

type RealtimeEventHandler func(eventType string, btyes json.RawMessage) error

type UserChannelListener struct {
	LogHandler    func(log string, err error)
	EventHandlers []RealtimeEventHandler
}

func (l *UserChannelListener) Log(s string) {
	if l.LogHandler != nil {
		l.LogHandler(s, nil)
	}
}

func (l *UserChannelListener) Err(s error) {
	if l.LogHandler != nil {
		l.LogHandler("", s)
	}
}
func (l *UserChannelListener) AddHandler(h RealtimeEventHandler) {
	l.EventHandlers = append(l.EventHandlers, h)
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
func (u *APIRealtime) Listen(ctx context.Context, listener *UserChannelListener) {
	if listener == nil {
		return
	}
	client := &http.Client{
		Timeout: 90 * time.Second,
	}
	var channel *UserChannel = nil
	var offset int64 = 0
	var waitForError = false
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if waitForError {
				waitForError = false
				c, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				<-c.Done()
				continue
			}
			if channel == nil {
				newChannel, err := u.GetUserChannel()
				if err != nil {
					listener.Err(err)
					waitForError = true
					continue
				}
				channel = newChannel
				offset = 0
				listener.Log(fmt.Sprintf("newChannel %s", newChannel.CometServer))
			}
			serverUrl, err := url.Parse(channel.CometServer)
			if err != nil {
				listener.Err(err)
				waitForError = true
				continue
			}
			q := serverUrl.Query()
			q.Set("offset", strconv.FormatInt(offset, 10))
			serverUrl.RawQuery = q.Encode()
			body, err := getUrl(ctx, client, serverUrl)
			if err != nil {
				listener.Err(err)
				waitForError = true
				continue
			}

			newComet, err := resolveComet(body)
			if err != nil {
				listener.Err(err)
				waitForError = true
				continue
			}
			if offset != newComet.NewOffset {
				listener.Log(fmt.Sprintf("offset %d", newComet.NewOffset))
			} else {
				// workaround.
				// Channel sometimes closes without any notification.
				// Call this API to activate it.
				u.GetUserChannel()
			}
			offset = newComet.NewOffset
			if offset == -3 {
				// Your offset is wrong and you need to resync your data
				channel = nil
				continue
			}
			for _, rawEvent := range newComet.Data {
				resolveEvent(rawEvent, listener)
			}
		}
	}
}

func getUrl(ctx context.Context, client *http.Client, _url *url.URL) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, _url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s, %v", _url.String(), err)
	}
	res, err := client.Do(req)
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

func resolveEvent(bytes json.RawMessage, listener *UserChannelListener) {
	if listener == nil {
		return
	}
	pass1 := &struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, pass1); err != nil {
		listener.Err(err)
	}
	if listener.EventHandlers != nil {
		for _, handler := range listener.EventHandlers {
			if err := handler(pass1.Type, bytes); err != nil {
				listener.Err(err)
			}
		}
	}
}

func NewResponseHandler(handler func(*NewResponseEvent)) RealtimeEventHandler {
	return standardRealtimeEventHandler("new_response", handler)
}
func NewPlurkHandler(handler func(*NewPlurkEvent)) RealtimeEventHandler {
	return standardRealtimeEventHandler("new_plurk", handler)
}
func UpdateNotificationHandler(handler func(*UpdateNotificationEvent)) RealtimeEventHandler {
	return standardRealtimeEventHandler("update_notification", handler)
}

func standardRealtimeEventHandler[T any](eventType string, handler func(*T)) RealtimeEventHandler {
	return func(inEventType string, bytes json.RawMessage) error {
		if inEventType == eventType {
			res := new(T)
			if err := json.Unmarshal(bytes, res); err != nil {
				return err
			}
			handler(res)
		}
		return nil
	}
}
