package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "log"
  "flag"
)

const API_URL string = "http://finance.google.com/finance/info"

func Decode(resp string, parser *Quotes) {
  json.Unmarshal([]byte(resp), &parser)
}

func buildUrl(quote string) string {
  url := []string{ API_URL, "?", "q=", quote }
  return strings.Join(url, "")
}

func getQuote(quote string) string {

  resp, err := http.Get(buildUrl(quote))

  if err != nil {
    log.Fatal("Could not get quote")
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  parsedResp := string(body)

  if err != nil {
    log.Fatal("Cannot read body")
  }

  return strings.Replace(parsedResp, "//", "", -1)
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

func main() {

  ch1 := make(chan Quotes)

  query := flag.String("quotes", "GOOGL,TSLA", "stock symbols separate by quotes")
  flag.Parse()

  for {

    go func (msg chan Quotes) {
      var quotes Quotes
      Decode(getQuote(*query), &quotes)
      ch1 <- quotes
    }(ch1)

    printQuote(<-ch1)
  }
}
