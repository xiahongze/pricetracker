package pushover

import (
	"fmt"
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
		User     string
		Msg      string
		Device   string
		Title    string
		URL      string
		Priority int
	}
)

var msgAPI = "https://api.pushover.net/1/messages.json"

// Send sends a message to the queue
func (c *Client) Send(msg *Message) error {
	if msg.Msg == "" {
		return fmt.Errorf("can't send an empty message")
	}
	user := msg.User
	if user == "" {
		user = c.User
	}
	form := url.Values{
		"token":   []string{c.AppToken},
		"user":    []string{user},
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
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Printf("resp: %s", body)
	}
	return nil
}
