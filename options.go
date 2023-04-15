package goplurk

import (
	"strconv"
	"strings"
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
	o.Get()["limited_to"] = "[" + strings.Join(limitedToStrs, ",") + "]"
	return o
}
func (o *PlurkAddOptions) NoComments() *PlurkAddOptions {
	o.Get()["no_comments"] = "1"
	return o
}
func (o *PlurkAddOptions) FriendsOnlyComments() *PlurkAddOptions {
	o.Get()["no_comments"] = "2"
	return o
}
func (o *PlurkAddOptions) Lang(lang string) *PlurkAddOptions {
	o.Get()["lang"] = lang
	return o
}
