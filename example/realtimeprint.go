package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kinabcd/goplurk"
)

func main() {
	var consumerToken = "..."
	var consumerSecret = "..."
	var accessToken = "..."
	var accessSecret = "..."
	client, err := goplurk.NewClient(consumerToken, consumerSecret, accessToken, accessSecret)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Print my plurks")
	ps, err := client.Timeline.GetPlurks(goplurk.NewGetPlurksOptions().FilterMy())
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, p := range ps.Plurks {
		fmt.Printf("%s %s\n", p.Qualifier, strings.ReplaceAll(p.ContentRaw, "\n", ""))
	}

	fmt.Println("start listen")
	listener := &goplurk.UserChannelListener{}
	listener.AddHandler(goplurk.NewPlurkHandler(func(e *goplurk.NewPlurkEvent) {
		fmt.Printf("(%d) Plurk %s %s\n", e.PlurkId, e.Plurk.Qualifier, strings.ReplaceAll(e.Plurk.ContentRaw, "\n", ""))
	}))
	listener.AddHandler(goplurk.NewResponseHandler(func(e *goplurk.NewResponseEvent) {
		fmt.Printf("(%d) Response %s %s\n", e.PlurkId, e.Response.Qualifier, strings.ReplaceAll(e.Response.ContentRaw, "\n", ""))
	}))
	listener.AddHandler(goplurk.UpdateNotificationHandler(func(e *goplurk.UpdateNotificationEvent) {
		fmt.Printf("Notification Noti=%d, Req=%d\n", e.Counts.Noti, e.Counts.Req)
		// Auto accept friendship request
		if alerts, err := client.Alerts.GetActive(); err == nil {
			for _, e := range alerts {
				if e, ok := e.(*goplurk.AlertsFriendshipRequestEvent); ok {
					client.Alerts.AddAsFriend(e.FromUser.Id)
				}
			}
		}
	}))
	listener.AddHandler(func(eType string, bytes json.RawMessage) error {
		if eType != "new_response" && eType != "new_plurk" && eType != "update_notification" {
			fmt.Printf("Unhandled %s: %s\n", eType, string(bytes))
		}
		return nil
	})
	listener.LogHandler = func(s string, err error) {
		if err != nil {
			fmt.Printf("(EE) %v\n", err)
		} else {
			fmt.Printf("(II) %s\n", s)
		}
	}

	client.Realtime.Listen(context.Background(), listener)
}
