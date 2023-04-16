package goplurk

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type APITimeline struct {
	client *Client
}

func (u *APITimeline) GetPlurk(plurkId int64) (*Plurk, error) {
	res := &struct {
		Plurk *Plurk `json:"plurk"`
	}{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Timeline/getPlurk", map[string]string{
		"plurk_id": strconv.FormatInt(plurkId, 10),
	}, res); err != nil {
		return nil, err
	}
	return res.Plurk, nil
}

func (u *APITimeline) GetPlurkCountsInfo(plurkId int64) (*PlurkCountsInfo, error) {
	res := &PlurkCountsInfo{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Timeline/getPlurkCountsInfo", map[string]string{
		"plurk_id": strconv.FormatInt(plurkId, 10),
	}, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Return plurks older than offset
// see GetPlurksOptions
func (u *APITimeline) GetPlurks(optionSets ...Options) (*Plurks, error) {
	return u.getPlurks("/APP/Timeline/getPlurks", optionSets)
}

// Return unread plurks older than offset
// see GetPlurksOptions
func (u *APITimeline) GetUnreadPlurks(optionSets ...Options) (*Plurks, error) {
	return u.getPlurks("/APP/Timeline/getUnreadPlurks", optionSets)
}

// Return public plurks older than offset
// see GetPlurksOptions
func (u *APITimeline) GetPublicPlurks(optionSets ...Options) (*Plurks, error) {
	return u.getPlurks("/APP/Timeline/getPublicPlurks", optionSets)
}

func (u *APITimeline) getPlurks(_api string, optionSets []Options) (*Plurks, error) {
	body := map[string]string{}
	for _, optionSet := range optionSets {
		maps.Copy(body, optionSet.Get())
	}
	plurks := &Plurks{}
	if err := u.client.Engine.CallAPIUnmarshal(_api, body, plurks); err != nil {
		return nil, err
	}
	return plurks, nil
}

func (u *APITimeline) PlurkAdd(qualifier string, content string, optionSets ...Options) (*Plurk, error) {
	if qualifier == "" {
		qualifier = ":"
	}
	if content == "" {
		return nil, fmt.Errorf("content can not be empty")
	}
	var body = map[string]string{}
	body["qualifier"] = qualifier
	body["content"] = content
	for _, optionSet := range optionSets {
		maps.Copy(body, optionSet.Get())
	}
	res, err := u.client.Engine.CallAPI("/APP/Timeline/plurkAdd", body)
	if err != nil {
		return nil, err
	}
	plurk := Plurk{}
	if err := json.Unmarshal(res, &plurk); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v, %s", err, string(res))
	}
	return &plurk, nil
}

func (u *APITimeline) PlurkDelete(plurkId int64) error {
	_, err := u.client.Engine.CallAPI("/APP/Timeline/plurkDelete", map[string]string{
		"plurk_id": strconv.FormatInt(plurkId, 10),
	})
	if err != nil {
		return fmt.Errorf("failed to delete plurk: %v", err)
	}
	return nil
}

func (u *APITimeline) MutePlurks(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/mutePlurks", plurkIds)
}

func (u *APITimeline) UnmutePlurks(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/unmutePlurks", plurkIds)
}

func (u *APITimeline) FavoritePlurks(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/favoritePlurks", plurkIds)
}

func (u *APITimeline) UnfavoritePlurks(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/unfavoritePlurks", plurkIds)
}

func (u *APITimeline) Replurk(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/replurk", plurkIds)
}

func (u *APITimeline) Unreplurk(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/unreplurk", plurkIds)
}

func (u *APITimeline) MarkAsRead(plurkIds []int64) error {
	return u.opPlurk("/APP/Timeline/markAsRead", plurkIds)
}

func (u *APITimeline) opPlurk(_url string, plurkIds []int64) error {
	if len(plurkIds) != 0 {
		plurkIdStrs := []string{}
		for _, limited := range plurkIds {
			plurkIdStrs = append(plurkIdStrs, strconv.FormatInt(limited, 10))
		}
		_, err := u.client.Engine.CallAPI(_url, map[string]string{
			"ids": "[" + strings.Join(plurkIdStrs, ",") + "]",
		})
		return err
	} else {
		return fmt.Errorf("plurkIds can not be empty")
	}

}
