package goplurk

import (
	"encoding/json"
	"time"
)

type User struct {
	Id                int64    `json:"id"`
	NickName          string   `json:"nick_name"`
	DisplayName       string   `json:"display_name"`
	FullName          string   `json:"full_name"`
	NameColor         string   `json:"name_color"`
	Premium           bool     `json:"premium"`
	HasProfileImage   int64    `json:"has_profile_image"`
	Avatar            int64    `json:"avatar"`
	ShowLocation      int64    `json:"show_location"`
	Location          string   `json:"location"`
	Timezone          string   `json:"timezone"`
	DefaultLang       string   `json:"default_lang"`
	DateFormat        int64    `json:"dateformat"`
	DateOfBirth       string   `json:"date_of_birth"`
	Birthday          Birthday `json:"birthday"`
	BdayPrivacy       int64    `json:"bday_privacy"`
	Gender            int64    `json:"gender"`
	Karma             float64  `json:"karma"`
	Recruited         int64    `json:"recruited"`
	Relationship      string   `json:"relationship"`
	Status            string   `json:"status"`
	TimelinePrivacy   int64    `json:"timeline_privacy"`
	VerifiedAccount   bool     `json:"verified_account"`
	FriendListPrivacy string   `json:"friend_list_privacy"`
	EmailConfirmed    bool     `json:"email_confirmed"`
	PhoneVerified     *int64   `json:"phone_verified"`
	PinnedPlurkId     *int64   `json:"pinned_plurk_id"`
	BackgroundId      int64    `json:"background_id"`
	ShowAds           bool     `json:"show_ads"`
	Version           string   `json:"_version"`
}

type Birthday struct {
	Year  int64 `json:"year"`
	Month int64 `json:"month"`
	Day   int64 `json:"day"`
}
type KarmaTrendNode struct {
	Time  time.Time
	Karma float64
}
type KarmaStates struct {
	CurrentKarma    float64  `json:"current_karma"`
	KarmaFallReason string   `json:"karma_fall_reason"`
	KarmaGraph      string   `json:"karma_graph"`
	KarmaTrendRaw   []string `json:"karma_trend"`
	KarmaTrend      []KarmaTrendNode
}

type Plurk struct {
	PlurkId             int64   `json:"plurk_id"`
	Qualifier           string  `json:"qualifier"`
	QualifierTranslated string  `json:"qualifier_translated"`
	IsUnread            int64   `json:"is_unread"`
	PlurkType           int64   `json:"plurk_type"`
	UserId              int64   `json:"user_id"`
	OwnerId             int64   `json:"owner_id"`
	Posted              string  `json:"posted"`
	NoComments          int64   `json:"no_comments"`
	Content             string  `json:"content"`
	ContentRaw          string  `json:"content_raw"`
	ResponseCount       int64   `json:"response_count"`
	ResponseSeen        int64   `json:"responses_seen"`
	LimitedTo           string  `json:"limited_to"`
	Favorite            bool    `json:"favorite"`
	FavoriteCount       int64   `json:"favorite_count"`
	Favorers            []int64 `json:"favorers"`
	Replurkable         bool    `json:"replurkable"`
	Replurked           bool    `json:"replurked"`
	ReplurkerId         int64   `json:"replurker_id"`
	ReplurkerCount      int64   `json:"replurkers_count"`
	Replurkers          []int64 `json:"replurkers"`
}

type Plurks struct {
	Plurks     []Plurk         `json:"plurks"`
	PlurkUsers map[string]User `json:"plurk_users"`
}

type Response struct {
	Id                  int64   `json:"id"`
	UserId              int64   `json:"user_id"`
	PlurkId             int64   `json:"plurk_id"`
	Content             string  `json:"content"`
	ContentRaw          string  `json:"content_raw"`
	Qualifier           string  `json:"qualifier"`
	QualifierTranslated *string `json:"qualifier_translated"`
	Posted              string  `json:"posted"`
	Lang                string  `json:"lang"`
	LastEdited          *string `json:"last_edited"`
	Coins               *string `json:"coins"`
	Editability         int64   `json:"editability"`
}

type Responses struct {
	Responses     []Response      `json:"responses"`
	ResponsesSeen int64           `json:"responses_seen"`
	ResponseCount int64           `json:"response_count"`
	Friends       map[string]User `json:"friends"`
}

type Profile struct {
	FriendsCount int64 `json:"friends_count"`
	FansCount    int64 `json:"fans_count"`

	UserInfo User   `json:"user_info"`
	Privacy  string `json:"privacy"`

	Plurks []Plurk `json:"plurks"`

	// OwnProfileOnly
	UnreadCount *int64          `json:"unread_count"`
	PlurksUsers map[string]User `json:"plurks_users"`

	// PublicProfileOnly
	AreFriends        *bool `json:"are_friends"`
	IsFan             *bool `json:"is_fan"`
	IsFollowing       *bool `json:"is_following"`
	HasReadPermission *bool `json:"has_read_permission"`
}

type UnreadCount struct {
	All       int64 `json:"all"`
	My        int64 `json:"my"`
	Private   int64 `json:"private"`
	Responsed int64 `json:"responded"`
}

type PlurkCountsInfo struct {
	ResponseCount  int64 `json:"response_count"`
	FavoriteCount  int64 `json:"favorite_count"`
	ReplurkerCount int64 `json:"replurkers_count"`
}

type UserChannel struct {
	ChannelName string `json:"channel_name"`
	CometServer string `json:"comet_server"`
}

// Realtime event "new_plurk"
type NewPlurkEvent struct {
	Plurk
}

// Realtime event "new_response"
type NewResponseEvent struct {
	PlurkId       int64    `json:"plurk_id"`
	Plurk         Plurk    `json:"plurk"`
	ResponseCount int64    `json:"response_count"`
	Response      Response `json:"response"`
	User          map[string]User
}

// Realtime event "update_notification"
// {"counts": {"noti": 1, "req": 0}, "type": "update_notification"}
type UpdateNotificationEvent struct {
	Counts struct {
		Noti int64 `json:"noti"`
		Req  int64 `json:"req"`
	} `json:"counts"`
}

type AlertsEvent struct {
	Type   string `json:"type"`
	Posted string `json:"posted"`
}

type AlertsUnhandledEvent struct {
	AlertsEvent
	RawMessage json.RawMessage
}

// requires action from the user
type AlertsFriendshipRequestEvent struct {
	AlertsEvent
	FromUser User `json:"from_user"`
}

// requires action from the user
type AlertsFriendshipPendingEvent struct {
	AlertsEvent
	FromUser User `json:"to_user"`
}

type AlertsNewFanEvent struct {
	AlertsEvent
	NewFan User `json:"new_fan"`
}

// Friendship accepted notification
type AlertsFriendshipAcceptedEvent struct {
	AlertsEvent
	FriendInfo User `json:"friend_info"`
}

// New friend notification
type AlertsNewFriendEvent struct {
	AlertsEvent
	NewFriend User `json:"new_friend"`
}

// New private plurk
type AlertsPrivatePlurkEvent struct {
	AlertsEvent
	Owner   User  `json:"owner"`
	PlurkId int64 `json:"plurk_id"`
}

// User's plurk got liked
type AlertsPlurkLikedEvent struct {
	AlertsEvent
	FromUser  User  `json:"from_user"`
	PlurkId   int64 `json:"plurk_id"`
	NumOthers int64 `json:"num_others"`
}

// User's plurk got replurked
type AlertsPlurkReplurkedEvent struct {
	AlertsEvent
	FromUser  User  `json:"from_user"`
	PlurkId   int64 `json:"plurk_id"`
	NumOthers int64 `json:"num_others"`
}

// User got mentioned in a plurk
type AlertsMentionedEvent struct {
	AlertsEvent
	FromUser  User  `json:"from_user"`
	PlurkId   int64 `json:"plurk_id"`
	NumOthers int64 `json:"num_others"`
	// response_id may be null if user was mentioned in the plurk and not in a response.
	ResponseId *int64 `json:"response_id"`
}

// User's own plurk got responded
type AlertsMyRespondedEvent struct {
	AlertsEvent
	FromUser   User  `json:"from_user"`
	PlurkId    int64 `json:"plurk_id"`
	NumOthers  int64 `json:"num_others"`
	ResponseId int64 `json:"response_id"`
}
