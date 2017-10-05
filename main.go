package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
		"Ohayo!!!",
	}
)

type opt struct {
	token   string
	channel string
	text    string
	by      string
	at      string
}

func main() {
	rand.Seed(time.Now().Unix())
	cli := new(opt)
	flag.StringVar(&cli.channel, "channel", "C0GCPHQNM", "Channel to say greeting")
	flag.StringVar(&cli.token, "token", "Unknown", "slack user token to say greeting")
	flag.StringVar(&cli.text, "text", "", "Message to say, if empty, pick random message each greeting")
	flag.StringVar(&cli.by, "by", "", "Name of bot")
	flag.StringVar(&cli.at, "at", "09:00", "related time around to send greeting, format: hh:mm")
	flag.Parse()

	now := time.Now()
	h, m, s := now.Clock()
	timestampAtDay := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	delta := ((33 * time.Hour) - timestampAtDay) % (24 * time.Hour)

	//send(cli.token, cli.channel, formatText(cli.text, cli.via))

	delta, err := randomTimeDelta(cli.at)

	fatalIfErr(err)
	fmt.Printf("Schedule to next greeting in %s, at %s\n", delta.String(), time.Now().Add(delta).String())
	timer := time.NewTimer(delta)

	for {
		t := <-timer.C
		fmt.Printf("Send a message at %s\n", t.String())
		for _, ch := range strings.Split(cli.channel, ",") {
			formatedText := formatText(cli.text, cli.by)
			send(cli.token, ch, formatedText)
		}
		delta, _ := randomTimeDelta(cli.at)
		fmt.Printf("Schedule to next greeting in %s, at %s\n", delta.String(), time.Now().Add(delta).String())
		timer.Reset(delta)
	}
}

func randomTimeDelta(at string) (time.Duration, error) {
	if ok, _ := regexp.MatchString("\\d{2}:\\d{2}", at); ok {
		atTimeArr := strings.Split(at, ":")
		atHour, _ := strconv.Atoi(atTimeArr[0])
		if atHour > 23 {
			return 0, errors.New("Invalid time, hour shoud < 24")
		}
		atMinute, _ := strconv.Atoi(atTimeArr[1])
		if atMinute > 59 {
			return 0, errors.New("Invalid time, min should < 60")
		}

		h, m, s := time.Now().Clock()
		tsnow := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
		delta := (time.Duration(24+atHour)*time.Hour + time.Duration(atMinute)*time.Minute - tsnow) % (24 * time.Hour)
		extra := time.Duration(rand.Intn(900)) * time.Second
		return delta + extra, nil
	}
	return 0, errors.New("Invalid time input, should be format at hh:mm")
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
