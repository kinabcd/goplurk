package main

import (
	"github.com/kinabcd/goplurk"
)

func main() {
	var consumerToken = "..."
	var consumerSecret = "..."
	var accessToken = "..."
	var accessSecret = "..."
	client, _ := goplurk.NewClient(consumerToken, consumerSecret, accessToken, accessSecret)
	client.Timeline.PlurkAdd("says", "somecontent")
}
