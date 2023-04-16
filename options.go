package goplurk

import (
	"strconv"
	"strings"
	"time"
)

type Options interface {
	Get() map[string]string
}

type baseOptions struct {
	options map[string]string
}

func (o *baseOptions) Get() map[string]string {
	if o.options == nil {
		o.options = map[string]string{}
	}
	return o.options
}
func (o *baseOptions) set(key, value string) {
	o.Get()[key] = value
}

type PlurkAddOptions struct {
	baseOptions
}

func NewPlurkAddOptions() *PlurkAddOptions {
	return &PlurkAddOptions{}
}

func (o *PlurkAddOptions) LimitedTo(ids ...int64) *PlurkAddOptions {
	limitedToStrs := []string{}
	for _, limited := range ids {
		limitedToStrs = append(limitedToStrs, strconv.FormatInt(limited, 10))
	}
	o.set("limited_to", "["+strings.Join(limitedToStrs, ",")+"]")
	return o
}
func (o *PlurkAddOptions) NoComments() *PlurkAddOptions {
	o.set("no_comments", "1")
	return o
}
func (o *PlurkAddOptions) FriendsOnlyComments() *PlurkAddOptions {
	o.set("no_comments", "2")
	return o
}
func (o *PlurkAddOptions) Lang(lang string) *PlurkAddOptions {
	o.set("lang", lang)
	return o
}

type GetPlurksOptions struct {
	baseOptions
}

func NewGetPlurksOptions() *GetPlurksOptions {
	return &GetPlurksOptions{}
}

func (o *GetPlurksOptions) Offset(offset time.Time) *GetPlurksOptions {
	o.set("offset", offset.UTC().Format(time.RFC3339))
	return o
}

// Limit the number of plurks returned
func (o *GetPlurksOptions) Limit(limit int64) *GetPlurksOptions {
	o.set("limit", strconv.FormatInt(limit, 10))
	return o
}

// Limit the my plurks returned
func (o *GetPlurksOptions) FilterMy() *GetPlurksOptions {
	o.set("filter", "my")
	return o
}

// Limit the favorite plurks returned
func (o *GetPlurksOptions) FilterFavorite() *GetPlurksOptions {
	o.set("filter", "favorite")
	return o
}

// Limit the replurked plurks returned
func (o *GetPlurksOptions) FilterReplurked() *GetPlurksOptions {
	o.set("filter", "replurked")
	return o
}

// Limit the private plurks returned
func (o *GetPlurksOptions) FilterPrivate() *GetPlurksOptions {
	o.set("filter", "private")
	return o
}

// Limit the responded plurks returned
func (o *GetPlurksOptions) FilterResponded() *GetPlurksOptions {
	o.set("filter", "responded")
	return o
}

// Limit the mentioned plurks returned
func (o *GetPlurksOptions) FilterMentioned() *GetPlurksOptions {
	o.set("filter", "mentioned")
	return o
}
