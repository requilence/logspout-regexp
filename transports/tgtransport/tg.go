package tgtransport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var defaultSameMessageMinInterval = time.Minute * 10

type TGTransport struct {
	botToken string
	chat     int64

	throttleSameInterval time.Duration
	throttleSame         map[string]time.Time
	throttleSameMutex    *sync.Mutex
}

func New(options map[string]string) (*TGTransport, error) {
	if options == nil {
		return nil, fmt.Errorf("options need to contain chat&token")
	}

	botToken := options["token"]
	if botToken == "" {
		return nil, fmt.Errorf("token not set")
	}

	chat := options["chat"]
	if chat == "" {
		return nil, fmt.Errorf("chat not set")
	}

	chatId, err := strconv.ParseInt(chat, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("incorrect chat option: should be a number")
	}

	var sameMessageMinInterval time.Duration

	throttleSeconds := options["throttle_seconds"]
	if throttleSeconds != "" {
		throttleSecondsInt, err := strconv.ParseInt(throttleSeconds, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("incorrect throttle_seconds option: should be a number")
		}

		sameMessageMinInterval = time.Second * time.Duration(throttleSecondsInt)
	} else {
		sameMessageMinInterval = defaultSameMessageMinInterval
	}

	return &TGTransport{
		botToken:             botToken,
		chat:                 chatId,
		throttleSameInterval: sameMessageMinInterval,
		throttleSame:         make(map[string]time.Time),
		throttleSameMutex:    &sync.Mutex{},
	}, nil
}

type Payload struct {
	ChatID int64  `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
	Text   string `json:"text"`
}

func (tgn *TGTransport) Name() string {
	return "tg"
}

func (tgn *TGTransport) Write(containerId, containerName, matchedString, re string) error {
	groupEvent := containerId + "/" + re

	tgn.throttleSameMutex.Lock()
	if t, exists := tgn.throttleSame[groupEvent]; exists {
		if time.Now().Sub(t) < tgn.throttleSameInterval {
			tgn.throttleSameMutex.Unlock()
			return nil
		}
	}

	tgn.throttleSame[groupEvent] = time.Now()
	tgn.throttleSameMutex.Unlock()

	matchedString = strings.Replace(matchedString, "&", "&amp;", -1)
	matchedString = strings.Replace(matchedString, "<", "&lt;", -1)
	matchedString = strings.Replace(matchedString, ">", "&gt;", -1)

	data := Payload{
		ChatID: tgn.chat,
		ParseMode: "html",
		Text:   fmt.Sprintf("Found log match for <strong>%s(%s)</strong> container\nregexp: <code>%s</code>\n<pre>%s</pre>", strings.TrimLeft(containerName,"/"), containerId, re, matchedString),
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", os.ExpandEnv("https://api.telegram.org/bot"+tgn.botToken+"/sendMessage"), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
