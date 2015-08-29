package main

import (
  "fmt"
  "strings"
  "encoding/json"
  "flag"
  "os"
  "bytes"
  "net/http"
  // "github.com/vaot/finance_now/yahoo_api"
  "github.com/vaot/finance_now/google_api"
)

func Decode(resp string, parser *Quotes) {
  json.Unmarshal([]byte(resp), &parser)
}

func printQuote(quotes Quotes) {
  for _, quote := range quotes {
    fmt.Printf("%s :::::: Trade Price: %.2f, LastTradeTime: %s, ChangePrice: %.2f, ChangePercentage: %.2f",
      quote.Symbol,
      quote.getTradePrice(),
      quote.LastTradeTime,
      quote.getChangePrice(),
      quote.getChangePricePercentage())

    fmt.Printf("\n")
  }

  fmt.Printf("\n")
}

type SlackRequest struct {
  Channel string `json:"channel"`
  Username string `json:"username"`
  Text string `json:"text"`
  IconEmoji string `json:"icon_emoji"`
}

func SlackHandler(price string) {
  slackPayload := &SlackRequest{}

  slackPayload.Channel = "#stocks"
  slackPayload.Username = "Finance Now"
  slackPayload.IconEmoji = ":moneybag:"
  slackPayload.Text = strings.Join([]string{"TSLA stocks just reached the limit set: ", price}, "")

  jsonString, _ := json.Marshal(slackPayload)
  var SLACK_WEBHOOK_URL string = os.Getenv("SLACK_WEBHOOK_URL")
  http.Post(SLACK_WEBHOOK_URL, "application/json", bytes.NewReader(jsonString))
}

func MapQuotesToLimits(quotes *string, limits *string) map[string]string {
  limitsMap := make(map[string]string)
  quotesArray := strings.Split(*quotes, ",")
  limitsArray := strings.Split(*limits, ",")

  for i := 0; i < len(quotesArray); i++ {
    limitsMap[quotesArray[i]] = limitsArray[i]
  }

  return limitsMap
}

func Watcher(quotes Quotes, mapping *map[string]string) {
  // To do
}

func main() {
  ch1 := make(chan Quotes)

  query := flag.String("quotes", "GOOGL,TSLA", "stock symbols separate by quotes")
  limits := flag.String("watch", "234.4,280.3", "stock symbols separate by quotes")

  mapping := MapQuotesToLimits(query, limits)

  flag.Parse()

  for {

    go func (msg chan Quotes) {
      var quotes Quotes
      Decode(google_api.GetQuote(*query), &quotes)
      ch1 <- quotes
    }(ch1)

    quote := <-ch1

    Watcher(quote, mapping)
    printQuote(quote)
  }
}
