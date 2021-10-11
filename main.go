package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Update struct {
	Kind     string `json:"kind"`
	Text     string `json:"text"`
	Image    string `json:"image"`
	Priority int    `json:"priority"`
	URL      string `json:"url"`
}

var URL = os.Getenv("BINGO_HOOK_URL")
var COLORS = [5]int{15158332, 15105570, 16776960, 3447003, 9807270}

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var update Update
		err := json.Unmarshal(scanner.Bytes(), &update)
		if err != nil {
			fmt.Println("Error parsing json")
		} else {
			message := MakeMessage(&update)
			SendMessage(client, message)
		}
	}

}

func SendMessage(client *http.Client, message interface{}) {
	messageJsonB, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error making json")
	}

	reqContent := bytes.NewReader(messageJsonB)
	req, err := http.NewRequest("POST", URL, reqContent)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	client.Do(req)
}

func MakeMessage(u *Update) map[string]interface{} {
	if u.Priority == 0 {
		u.Priority = 5
	} else if u.Priority < 1 {
		u.Priority = 1
	} else if u.Priority > 5 {
		u.Priority = 5
	}
	embed := map[string]interface{}{
		"title": u.Text,
		"color": COLORS[u.Priority-1],
	}
	if strings.TrimSpace(u.Image) != "" {
		embed["thumbnail"] = map[string]string{
			"url": strings.TrimSpace(u.Image),
		}
	}

	if strings.TrimSpace(u.URL) != "" {
		embed["url"] = strings.TrimSpace(u.URL)
	}

	message := map[string]interface{}{
		"username": u.Kind,
		"content":  fmt.Sprintf("PT%d", u.Priority),
		"embeds":   []map[string]interface{}{embed},
	}

	return message

}
