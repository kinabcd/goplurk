package main

import (
	"context"
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
	client.Realtime.Listen(context.Background(), func(event interface{}) {
		if e, ok := event.(*goplurk.NewResponseEvent); ok {
			fmt.Printf("(%d) Response %s %s\n", e.PlurkId, e.Response.Qualifier, strings.ReplaceAll(e.Response.ContentRaw, "\n", ""))
		} else if e, ok := event.(*goplurk.NewPlurkEvent); ok {
			fmt.Printf("(%d) Plurk %s %s\n", e.PlurkId, e.Plurk.Qualifier, strings.ReplaceAll(e.Plurk.ContentRaw, "\n", ""))
		} else if e, ok := event.(*goplurk.UpdateNotificationEvent); ok {
			fmt.Printf("Notification Noti=%d, Req=%d\n", e.Counts.Noti, e.Counts.Req)
		} else if e, ok := event.(*goplurk.RealtimeLogEvent); ok {
			if e.Err != nil {
				fmt.Printf("(EE) %v\n", e.Err)
			} else {
				fmt.Printf("(II) %s\n", *e.Log)
			}
		}
	})
}
