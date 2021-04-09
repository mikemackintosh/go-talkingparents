package tp

import (
	"fmt"
	"time"
	"math"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/hako/durafmt"
)
// Thread is a thread view, which contains a slice of messages.
type Thread struct {
	Userid             int       `json:"UserID"`
	Subject            string    `json:"Subject"`
	Timezoneidentifier string    `json:"TimeZoneIdentifier"`
	Messages           []Message `json:"ThreadItems"`
}

// Message is a message within a thread.
type Message struct {
	Itemtype        int         `json:"ItemType"`
	Itemid          int         `json:"ItemID"`
	Message         string      `json:"Message"`
	EntryDate       string      `json:"EntryDate"`
	EntryDateUTC    string      `json:"EntryDateUtc"`
	AltentryDateUTC interface{} `json:"AltEntryDateUtc"`
	Caseid          int         `json:"CaseID"`
	Threadid        int         `json:"ThreadID"`
	Userid          int         `json:"UserID"`
	Name            string      `json:"Name"`
	Filename        interface{} `json:"FileName"`
	Initials        string      `json:"Initials"`
	Duration        int         `json:"Duration"`
	Attachments     interface{} `json:"Attachments"`
	IsBlobDeleted   bool        `json:"IsBlobDeleted"`
}

// GetThread will get the thread for the corresponding conversation ID.
func (c *Client) GetThread(id int) (*Thread, error) {
	var err error

	// make the request
	req, err := http.NewRequest("GET", ENDPOINT_THREADS, nil)

	// set the query string for the thread id
	q := req.URL.Query()
	q.Set("id", strconv.Itoa(id))
	req.URL.RawQuery = q.Encode()

	// set standard http headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-US,en")
	req.Header.Add("Origin", "https://talkingparents.com")
	req.Header.Add("User-Agent", USER_AGENT)

	// Set security headers
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Dest", "empty")

	// send it
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("thread %d was not found", id)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for thread %d", id)
	}

	var thread Thread
	if err := json.Unmarshal(body, &thread); err != nil {
		return nil, fmt.Errorf("error parsing thread %d", id)
	}

	return &thread, nil
}


// GetUntimelyMessages will list messages where they were not viewed within a specific timeframe.
// Ex: thread.GetUntimelyMessages(24 * time.Hour)
func (t *Thread) GetUntimelyMessages(d time.Duration) ([]string) {
	var lastMessage = map[string]*time.Time{}
	var output = []string{}
	for _, m := range t.Messages {
		people.AddPersonIfNotExists(m.Name)

		ts, _ := time.Parse(DATE_FORMAT, m.EntryDate)

		switch m.Itemtype {
		case SENT_MESSAGE:
			lastMessage[m.Name] = &ts
		case VIEW_MESSAGE:
			sender, _ := people.Not(m.Name)
			if t, _ := lastMessage[sender]; t != nil {
				delta := ts.Sub(*t)
				if delta > d {
					output = append(output, fmt.Sprintf("%s sent at %s, %s viewed at %s, %v later", sender, t, m.Name, ts, durafmt.Parse(delta).String()))
				}
			}
		default:
			continue
		}
	}

	return output
}

// RoundTime will round time up or down.
func RoundTime(input float64) int {
	var result float64
	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}
	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)
	return int(i)
}
