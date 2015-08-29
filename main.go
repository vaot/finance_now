package main

import (
  "fmt"
  "strings"
  "encoding/json"
  "flag"
  "os"
  "bytes"
  "net/http"
  "strconv"
  // "github.com/vaot/finance_now/yahoo_api"
  "github.com/vaot/finance_now/google_api"
  "github.com/fzzy/radix/redis"
  "time"
)

var client, _ = redis.Dial("tcp", "localhost:6379")
const MAX_ALERTS int = 10

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

func SlackHandler(symbol string, price float64) {
  slackPayload := &SlackRequest{}

  priceStr := strconv.FormatFloat(price, 'f', 3, 64)

  slackPayload.Channel = "#stocks"
  slackPayload.Username = "Finance Now"
  slackPayload.IconEmoji = ":moneybag:"
  slackPayload.Text = strings.Join([]string{ symbol, " stocks just reached the limit set: ", priceStr }, "")

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

func ShouldRunHandler(quote string) bool {
  status,_ := client.Cmd("HGET", "alerts", quote).Str()
  client.Cmd("HINCRBY", "alerts:times", quote, -1)
  times,_ := client.Cmd("HGET", "alerts:times", quote).Int()
  return (status == "running" && times > 0)
}

func Watcher(quotes Quotes, mapping map[string]string) {
  for _, quote := range quotes {

    if val, ok := mapping[quote.Symbol]; ok {

      fVal,_ := strconv.ParseFloat(val, 32)

      if fVal >= quote.getTradePrice() && ShouldRunHandler(quote.Symbol) {
        time.Sleep(2000 * time.Millisecond)
        go SlackHandler(quote.Symbol, quote.getTradePrice())
      }

    }

  }
}

func main() {
  ch1 := make(chan Quotes)

  query := flag.String("quotes", "GOOGL,TSLA", "stock symbols separate by quotes")
  limits := flag.String("limits", "234.4,280.3", "stock symbols separate by quotes")

  flag.Parse()

  formattedQuery := strings.ToUpper(*query)
  query = &formattedQuery

  mapping := MapQuotesToLimits(query, limits)

  for key,_ := range mapping {
    client.Cmd("HSET", "alerts", key, "running")
    client.Cmd("HSET", "alerts:times", key, MAX_ALERTS)
  }

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
