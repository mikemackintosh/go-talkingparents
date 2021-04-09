package main

import (
	"flag"
	"log"
	"time"

  "github.com/mikemackintosh/talkingparents"
)

var (
	// flag variables
	flagUsername  string
	flagPassword  string
	flagTimeframe time.Duration
)


// Set flags
func init() {
	flag.StringVar(&flagUsername, "u", "", "Username")
	flag.StringVar(&flagPassword, "p", "", "Password")
	flag.DurationVar(&flagTimeframe, "w", 24*time.Hour, "Time window")
}

// main
func main() {
	flag.Parse()

  // Create a new client.
  client, err := tp.NewClient()
  if err != nil {
    log.Println(err)
  }

  // Authenticate to the client.
  if err := client.Authenticate(flagUsername, flagPassword); err != nil {
    log.Println(err)
  }

  conversations, err := client.ListConversations()
  if err != nil {
    log.Println(err)
  }

  for _, c := range conversations.Issues {
    thread, err := client.GetThread(c.ThreadId)
    if err != nil {
      log.Println(err)
    }

    for _, m := range thread.GetUntimelyMessages(flagTimeframe) {
      log.Println(m)
    }
  }
}
