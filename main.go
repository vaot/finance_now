package main

import (
  "fmt"
  "strings"
  "flag"
  "strconv"
  // "github.com/vaot/finance_now/yahoo_api"
  "github.com/vaot/finance_now/google_api"
  "github.com/fzzy/radix/redis"
  "time"
  "os"
  "net/url"
)

var redisUrl ,_ = url.Parse(os.Getenv("REDIS_URL"))
var client, _ = redis.Dial("tcp", redisUrl.Host)

const MAX_ALERTS int = 10

func printQuote(quotes google_api.Quotes) {
  for _, quote := range quotes {
    fmt.Printf("%s :::::: Trade Price: %.2f, LastTradeTime: %s, ChangePrice: %.2f, ChangePercentage: %.2f",
      quote.Symbol,
      quote.GetTradePrice(),
      quote.LastTradeTime,
      quote.GetChangePrice(),
      quote.GetChangePricePercentage())

    fmt.Printf("\n")
  }

  fmt.Printf("\n")
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
  fmt.Println(times)
  fmt.Println(status)
  return (status == "running" && times > 0)
}

func Watcher(quotes google_api.Quotes, mapping map[string]string) {
  for _, quote := range quotes {

    if val, ok := mapping[quote.Symbol]; ok {

      fVal,_ := strconv.ParseFloat(val, 32)
      fmt.Println("Watching for limit: " + quote.Symbol + " " + val)
      fmt.Println(ShouldRunHandler(quote.Symbol))
      if quote.GetTradePrice() >= fVal && ShouldRunHandler(quote.Symbol) {
        time.Sleep(2000 * time.Millisecond)
        go SlackHandler(quote.Symbol, quote.GetTradePrice())
      }

    }

  }
}

func main() {
  ch1 := make(chan google_api.Quotes)

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

    go func (msg chan google_api.Quotes) {
      var quotes google_api.Quotes
      google_api.Decode(google_api.GetQuote(*query), &quotes)
      ch1 <- quotes
    }(ch1)

    quote := <-ch1

    Watcher(quote, mapping)
    printQuote(quote)
  }
}
