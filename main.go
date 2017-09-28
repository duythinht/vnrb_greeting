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
	via     string
}

func main() {
	rand.Seed(time.Now().Unix())
	cli := new(opt)
	flag.StringVar(&cli.channel, "channel", "C0GCPHQNM", "Channel to say greeting")
	flag.StringVar(&cli.token, "token", "Unknown", "slack user token to say greeting")
	flag.StringVar(&cli.text, "text", "", "Message to say, if empty, pick random message each greeting")
	flag.StringVar(&cli.via, "via", "", "Name of bot")
	flag.Parse()

	now := time.Now()
	h, m, s := now.Clock()
	timestampAtDay := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	delta := ((33 * time.Hour) - timestampAtDay) % (24 * time.Hour)

	//send(cli.token, cli.channel, formatText(cli.text, cli.via))

	fmt.Printf("Schedule to next greeting in %s\n", delta.String())

	timer := time.NewTimer(delta)

	for {
		t := <-timer.C
		formatedText := formatText(cli.text, cli.via)
		fmt.Printf("Send a message at %s\n", t.String())
		send(cli.token, cli.channel, formatedText)
		timer.Reset(24 * time.Hour)
	}
}

func send(token string, channel string, text string) {
	requestURL := fmt.Sprintf("%s?token=%s&channel=%s&text=%s", VNRB_SLACK_URL, token, channel, text)
	fmt.Println("Request to:", requestURL)
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

//format message
func formatText(text string, via string) string {
	var message string
	if len(text) == 0 {
		message = greetings[rand.Intn(len(greetings))]
	} else {
		message = text
	}

	if len(via) == 0 {
		return url.QueryEscape(fmt.Sprintf("%s", message))
	}
	return url.QueryEscape(fmt.Sprintf("%s\n_- by %s -_", message, via))
}
