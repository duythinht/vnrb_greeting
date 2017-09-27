package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	VNRB_SLACK_URL = "https://vietnamrb.slack.com/api/chat.postMessage"
)

type opt struct {
	token   string
	channel string
	text    string
}

func main() {
	cli := new(opt)
	flag.StringVar(&cli.channel, "channel", "C0GCPHQNM", "Channel to say greeting")
	flag.StringVar(&cli.token, "token", "Unknown", "slack user token to say greeting")
	flag.StringVar(&cli.text, "text", "Good morning", "Message to say")
	flag.Parse()

	now := time.Now()

	currentMin := now.Minute()
	hourDelta := (33 - now.Hour()) % 24

	delta := time.Duration(hourDelta)*time.Hour - time.Duration(currentMin)*time.Minute

	fmt.Printf("Schedule to next greeting in %s\n", delta.String())

	timer := time.NewTimer(delta)

	for {
		t := <-timer.C
		fmt.Printf("Send message '%s' at %s", cli.text, t.String())
		send(cli.token, cli.channel, cli.text)
		timer.Reset(24 * time.Hour)
	}
}

func send(token string, channel string, text string) {
	requestURL := fmt.Sprintf("%s?token=%s&channel=%s&text=%s", VNRB_SLACK_URL, token, channel, url.QueryEscape(text))
	resp, err := http.Get(requestURL)
	fatalIfErr(err)
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		fatalIfErr(errors.New(fmt.Sprintf("Error when request to slack api, status: %d", resp.StatusCode)))
	}
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
