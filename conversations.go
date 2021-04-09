package tp

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

// Conversations is a list of available conversations.
type Conversations struct {
	Issues []Conversation `json:"Issues"`
	ThreadCount int         `json:"ThreadCount"`
	SearchTerm  interface{} `json:"SearchTerm"`
}

// Conversation is a conversation summary.
type Conversation struct {
  ThreadId           int           `json:"ThreadID"`
  Subject            string        `json:"Subject"`
  LastReplyDate      string        `json:"LastReplyDate"`
  Hascall            bool          `json:"HasCall"`
  Isnew              bool          `json:"IsNew"`
  UserId             int           `json:"UserID"`
  ThreadItems        []interface{} `json:"ThreadItems"`
  TimezoneIdentifier interface{}   `json:"TimeZoneIdentifier"`
  LastReplyName      string        `json:"LastReplyName"`
  CreateDate         string        `json:"CreateDate"`
  CreatedBy          string        `json:"CreatedBy"`
  FromReplyId        int           `json:"FromReplyID"`
}


// ListConversations will list all conversations.
func (c *Client) ListConversations() (*Conversations, error) {
	var err error

	// make the request
	req, err := http.NewRequest("GET", ENDPOINT_CONVERSATIONS, nil)

	// set the query string for the thread id
	q := req.URL.Query()
	q.Set("p", "1")
	req.URL.RawQuery = q.Encode()

	// set standard http headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-US,en")
	req.Header.Add("Origin", "https://talkingparents.com")
	req.Header.Add("User-Agent", USER_AGENT)

	// send it
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("conversations not found")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for conversations")
	}

	var conversations Conversations
	if err := json.Unmarshal(body, &conversations); err != nil {
		return nil, fmt.Errorf("error parsing conversations")
	}

	return &conversations, nil
}
