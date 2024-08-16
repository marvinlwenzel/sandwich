package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const VERSION = "0.0.1"
const DEFAULT_WAIT_SECONDS = 15

type discordUrl string
type webhookUrl string
type discordUserId int64

// func (id discordUserId) String() string {
// 	return string(int64(id))
// }

type WebhookMessage struct {
	Content         string          `json:"content"`
	Username        string          `json:"username"`
	Flags           int32           `json:"flags"`
	AllowedMentions DiscordMentions `json:"allowed_mentions"`
}

func (wm WebhookMessage) MarshalJSON() ([]byte, error) {
	if len(wm.AllowedMentions.Users) == 0 && wm.Flags == 0 {
		return json.Marshal(struct {
			Content  string `json:"content"`
			Username string `json:"username"`
		}{
			Content:  wm.Content,
			Username: wm.Username,
		})
	} else if len(wm.AllowedMentions.Users) == 0 && wm.Flags > 0 {
		return json.Marshal(struct {
			Content  string `json:"content"`
			Username string `json:"username"`
			Flags    int32  `json:"flags"`
		}{
			Content:  wm.Content,
			Username: wm.Username,
			Flags:    wm.Flags,
		})
	} else if len(wm.AllowedMentions.Users) > 0 && wm.Flags == 0 {
		return json.Marshal(struct {
			Content         string          `json:"content"`
			Username        string          `json:"username"`
			AllowedMentions DiscordMentions `json:"flags"`
		}{
			Content:         wm.Content,
			Username:        wm.Username,
			AllowedMentions: wm.AllowedMentions,
		})
	} else {
		return json.Marshal(struct {
			Content         string          `json:"content"`
			Username        string          `json:"username"`
			Flags           int32           `json:"flags"`
			AllowedMentions DiscordMentions `json:"allowed_mentions"`
		}{
			Content:         wm.Content,
			Username:        wm.Username,
			Flags:           wm.Flags,
			AllowedMentions: wm.AllowedMentions,
		})
	}
}

type DiscordMentions struct {
	Users []discordUserId `json:"users"`
}

func (dms DiscordMentions) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(dms.Users))
	for i := range dms.Users {
		ss[i] = strconv.FormatInt(int64(dms.Users[i]), 10)
	}
	return json.Marshal(ss)
}
func main() {
	rawTargetUrl, hasTargetUrl := os.LookupEnv("SANDWICH_CHECK_URL")
	if !hasTargetUrl {
		fmt.Fprintln(os.Stderr, "No URL from `SANDWICH_CHECK_URL`")
		os.Exit(101)
	}
	targetUrl := discordUrl(rawTargetUrl)

	rawWebhook, hasWebhook := os.LookupEnv("SANDWICH_WEBHOOK")
	if !hasWebhook {
		fmt.Fprintln(os.Stderr, "No URL from `SANDWICH_WEBHOOK`")
		os.Exit(102)
	}
	webhook := webhookUrl(rawWebhook)

	rawWaitSeconds, hasWaitSeconds := os.LookupEnv("SANDWICH_INTERVAL_SECONDS")
	var (
		waitSeconds uint64
	)
	if !hasWaitSeconds {
		fmt.Printf("No SANDWICH_INTERVAL_SECONDS set. Using default of %s\n", DEFAULT_WAIT_SECONDS)
		waitSeconds = DEFAULT_WAIT_SECONDS
	} else {
		num, err := strconv.ParseUint(rawWaitSeconds, 10, 64)
		if err != nil {
			fmt.Printf("SANDWICH_INTERVAL_SECONDS '%s' can not be parsed. Panic\n", DEFAULT_WAIT_SECONDS)
			panic("Can not parse int from env")
		}
		waitSeconds = num
	}

	targetDownUserIds := discordUserIdsFromEnvVar("SANDWICH_TARGET_DOWN_IDS")

	fmt.Printf("S.A.N.D.W.I.C.H. Version %s\n", VERSION)

	fmt.Printf("Targeting %s\n", targetUrl)
	fmt.Printf("Reporting to Webhook %s\n", webhook)

	fmt.Printf("Pinging target down to %s\n", targetDownUserIds)
	fmt.Printf("Pinging in interval of %s seconds\n", waitSeconds)

	helloData := WebhookMessage{
		Content:  fmt.Sprintf("I am back. As version %s", VERSION),
		Username: "S.A.N.D.W.I.C.H.",
		Flags:    4096,
		AllowedMentions: DiscordMentions{
			Users: []discordUserId{},
		},
	}

	hello, _ := json.Marshal(helloData)
	resp, err := http.Post(string(webhook), "application/json", bytes.NewBuffer(hello))

	if err != nil || int(resp.StatusCode/100) != 2 {
		fmt.Printf("Err: %s\n", err)
		fmt.Printf("resp: %s\n", resp)
		panic("Could not send Welcome Message to Webhook. Config or Impl might be screwed.")
	}

	go checkFoundry(targetUrl, webhook, targetDownUserIds, waitSeconds)
	select {}
}

type FoundryStatus int64

const (
	Startup FoundryStatus = iota
	Dead
	InternalError
	UpAndNoWorld
	UpAndWorld
)

func (status FoundryStatus) String() string {
	return []string{"StartUp", "Dead", "InternalError", "UpAndNoWorld", "UpAndWorld"}[status]
}

func checkFoundry(targetUrl discordUrl, webhook webhookUrl, targetDownUserIds []discordUserId, waitSeconds uint64) {
	waitTime := time.Duration(waitSeconds) * time.Second
	status := Startup
	for {
		newStatus, worldName := foundryStatusOfServer(string(targetUrl))
		if newStatus == status {
			// nothing per default
		} else if newStatus == Dead || newStatus == InternalError {
			pingString := ""
			for i := range targetDownUserIds {
				pingString = fmt.Sprintf("%s <@%s>", pingString, strconv.FormatInt(int64(targetDownUserIds[i]), 10))
			}

			data := WebhookMessage{
				Content:  fmt.Sprintf("%s\n# %s has gone %s", pingString, targetUrl, newStatus),
				Username: "S.A.N.D.W.I.C.H.",
				Flags:    0,
				AllowedMentions: DiscordMentions{
					Users: []discordUserId{},
				},
			}
			body, _ := json.Marshal(data)
			_, _ = http.Post(string(webhook), "application/json", bytes.NewBuffer(body))
			fmt.Printf("New Status %s", newStatus)
		} else if newStatus == UpAndNoWorld {
			data := WebhookMessage{
				Content:  fmt.Sprintf("%s is online with no world open.", targetUrl),
				Username: "S.A.N.D.W.I.C.H.",
				Flags:    4096,
				AllowedMentions: DiscordMentions{
					Users: []discordUserId{},
				},
			}
			body, _ := json.Marshal(data)
			_, _ = http.Post(string(webhook), "application/json", bytes.NewBuffer(body))
			fmt.Printf("New Status %s", newStatus)
		} else if newStatus == UpAndWorld {
			data := WebhookMessage{
				Content:  fmt.Sprintf("%s is hosting %s", targetUrl, worldName),
				Username: "S.A.N.D.W.I.C.H.",
				Flags:    4096,
				AllowedMentions: DiscordMentions{
					Users: []discordUserId{},
				},
			}
			body, _ := json.Marshal(data)
			_, _ = http.Post(string(webhook), "application/json", bytes.NewBuffer(body))
			fmt.Printf("New Status %s", newStatus)
		}
		status = newStatus
		time.Sleep(waitTime)
	}
}

func foundryStatusOfServer(s string) (FoundryStatus, string) {
	resp, err := http.Get(s)
	if err != nil {
		return Dead, ""
	}

	if resp.StatusCode >= 500 {
		return InternalError, ""
	}

	title, err := getPageTitle(resp.Body)

	if err != nil {
		return InternalError, ""
	}

	if title == "Foundry Virtual Tabletop" {
		return UpAndNoWorld, ""
	}

	return UpAndWorld, title
}

func getPageTitle(respBody io.Reader) (string, error) {
	// Parse the HTML document
	doc, err := html.Parse(respBody)
	if err != nil {
		return "", err
	}

	// Recursively search for the <title> element
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			return n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			title := f(c)
			if title != "" {
				return title
			}
		}
		return ""
	}

	return f(doc), nil
}

func discordUserIdsFromEnvVar(envName string) []discordUserId {
	rawTargetDownIds, hasTargetDownIds := os.LookupEnv(envName)
	var (
		targetDownUserIdStrings []string
	)
	if hasTargetDownIds {
		targetDownUserIdStrings = strings.Split(rawTargetDownIds, ",")
	} else {
		targetDownUserIdStrings = []string{}
	}

	targetDownUserIds := make([]discordUserId, len(targetDownUserIdStrings))
	for i, s := range targetDownUserIdStrings {
		num, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			fmt.Printf("Raw: %s\n", rawTargetDownIds)
			fmt.Printf("Err: %s\n", err)
			panic("Can not parse userid from string")
		}
		targetDownUserIds[i] = discordUserId(num)
	}
	return targetDownUserIds
}