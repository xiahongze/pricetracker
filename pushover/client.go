package pushover

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type (
	// Client is the Pushover client
	Client struct {
		AppToken string
		User     string
	}

	// Message is the request body for Pushover
	// Not all fields have been implemented
	// Please refer to https://pushover.net/api
	Message struct {
		Msg      string
		Device   string
		Title    string
		URL      string
		Priority int
	}
)

var msgAPI = "https://api.pushover.net/1/messages.json"

// Send sends a message to the queue
func (c *Client) Send(msg *Message) {
	if msg.Msg == "" {
		log.Printf("can't send an empty message")
		return
	}
	form := url.Values{
		"token":   []string{c.AppToken},
		"user":    []string{c.User},
		"message": []string{msg.Msg},
	}
	if msg.Device != "" {
		form.Add("device", msg.Device)
	}
	if msg.Title != "" {
		form.Add("title", msg.Title)
	}
	if msg.URL != "" {
		form.Add("url", msg.URL)
	}
	if msg.Priority != 0 {
		form.Add("priority", string(msg.Priority))
	}
	resp, err := http.PostForm(msgAPI, form)
	if err != nil {
		log.Printf("failed sending message %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed reading resp body %v", err)
			return
		}
		log.Printf("resp: %s", body)
	}
}
