package goplurk

import (
	"encoding/json"
	"strconv"
)

type APIAlerts struct {
	client *Client
}

func (a *APIAlerts) GetActive() ([]any, error) {
	datas := []json.RawMessage{}
	if err := a.client.Engine.CallAPIUnmarshal("/APP/Alerts/getActive", map[string]string{}, &datas); err != nil {
		return nil, err
	}
	res := make([]any, len(datas))
	for i, data := range datas {
		res[i] = resolveAlertsEvent(data)
	}
	return res, nil
}

func (a *APIAlerts) GetHistory() ([]any, error) {
	datas := []json.RawMessage{}
	if err := a.client.Engine.CallAPIUnmarshal("/APP/Alerts/getHistory", map[string]string{}, &datas); err != nil {
		return nil, err
	}
	res := make([]any, len(datas))
	for i, data := range datas {
		res[i] = resolveAlertsEvent(data)
	}
	return res, nil
}
func (a *APIAlerts) AddAsFan(userId int64) error {
	_, err := a.client.Engine.CallAPI("/APP/Alerts/addAsFan", map[string]string{
		"user_id": strconv.FormatInt(userId, 10),
	})
	return err
}
func (a *APIAlerts) AddAllAsFan() error {
	_, err := a.client.Engine.CallAPI("/APP/Alerts/addAllAsFriends", map[string]string{})
	return err
}

func (a *APIAlerts) AddAsFriend(userId int64) error {
	_, err := a.client.Engine.CallAPI("/APP/Alerts/addAsFriend", map[string]string{
		"user_id": strconv.FormatInt(userId, 10),
	})
	return err
}
func (a *APIAlerts) DenyAsFriend(userId int64) error {
	_, err := a.client.Engine.CallAPI("/APP/Alerts/denyFriendship", map[string]string{
		"user_id": strconv.FormatInt(userId, 10),
	})
	return err
}
func (a *APIAlerts) RemoveNotification(userId int64) error {
	_, err := a.client.Engine.CallAPI("/APP/Alerts/removeNotification", map[string]string{
		"user_id": strconv.FormatInt(userId, 10),
	})
	return err
}
func resolveAlertsEvent(bytes json.RawMessage) any {
	pass1 := &AlertsEvent{}
	if err := json.Unmarshal(bytes, pass1); err != nil {
		return &AlertsUnhandledEvent{
			AlertsEvent: *pass1,
			RawMessage:  bytes,
		}
	}

	var pass2 any = nil
	switch pass1.Type {
	case "friendship_request":
		pass2 = &AlertsFriendshipRequestEvent{}
	case "friendship_pending":
		pass2 = &AlertsFriendshipPendingEvent{}
	case "new_fan":
		pass2 = &AlertsNewFanEvent{}
	case "friendship_accepted":
		pass2 = &AlertsFriendshipAcceptedEvent{}
	case "new_friend":
		pass2 = &AlertsNewFriendEvent{}
	case "private_plurk":
		pass2 = &AlertsPrivatePlurkEvent{}
	case "plurk_liked":
		pass2 = &AlertsPlurkLikedEvent{}
	case "plurk_replurked":
		pass2 = &AlertsPlurkReplurkedEvent{}
	case "mentioned":
		pass2 = &AlertsMentionedEvent{}
	case "my_responded":
		pass2 = &AlertsMyRespondedEvent{}
	default:
		return &AlertsUnhandledEvent{
			AlertsEvent: *pass1,
			RawMessage:  bytes,
		}
	}
	if err := json.Unmarshal(bytes, pass2); err != nil {
		return &AlertsUnhandledEvent{
			AlertsEvent: *pass1,
			RawMessage:  bytes,
		}
	}
	return pass2
}
