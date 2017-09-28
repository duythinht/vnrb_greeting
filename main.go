package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	VNRB_SLACK_URL = "https://vietnamrb.slack.com/api/chat.postMessage"
)

var (
	greetings = []string{
		"Good morning!!!",
		"Morning",
		"Morning :)",
		"Good morning all",
		"Gooood morningggggg!!!",
		"Morning :smile:",
		"Morning everyone :coffee:",
	}
)

type opt struct {
	token   string
	channel string
	text    string
}

func main() {
	rand.Seed(time.Now().Unix())
	cli := new(opt)
	flag.StringVar(&cli.channel, "channel", "C0GCPHQNM", "Channel to say greeting")
	flag.StringVar(&cli.token, "token", "Unknown", "slack user token to say greeting")
	flag.StringVar(&cli.text, "text", pickRandomGreeting(), "Message to say")
	flag.Parse()

	now := time.Now()

	currentMin := now.Minute()
	hourDelta := (33 - now.Hour()) % 24

	delta := time.Duration(hourDelta)*time.Hour - time.Duration(currentMin)*time.Minute

	//send(cli.token, cli.channel, cli.text)

	fmt.Printf("Schedule to next greeting in %s to channel: %s, text: '%s'\n", delta.String(), cli.channel, cli.text)

	timer := time.NewTimer(delta)

	for {
		t := <-timer.C
		fmt.Printf("Send message '%s' at %s\n", cli.text, t.String())
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

func pickRandomGreeting() string {
	return greetings[rand.Intn(len(greetings))]
}
