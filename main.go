package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "flag"
)

const API_URL string = "http://finance.google.com/finance/info"

type Quote struct {
  TradePrice       float64 `json:"l,string"`
  LastTradeTime    string  `json:"ltt"`
  ChangePrice      float64 `json:"c,string"`
  ChangePercentage float64 `json:"cp,string"`
}

func Decode(resp string, parser *[]Quote) {
  json.Unmarshal([]byte(resp), &parser)
}

func buildUrl(quote string) string {
  url := []string{ API_URL, "?", "q=", quote }
  return strings.Join(url, "")
}

func getQuote(quote string) string {

  resp, err := http.Get(buildUrl(quote))

  if err != nil {
    fmt.Println("Could not get quote")
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  parsedResp := string(body)

  if err != nil {
    fmt.Println("Arghhh")
  }

  fmt.Println(parsedResp)

  return strings.Replace(parsedResp, "//", "", -1)
}

func printQuote(quotes []Quote) {
  for _, quote := range quotes {
    fmt.Printf("Trade Price: %.3f, LastTradeTime: %s, ChangePrice: %.3f, ChangePercentage: %.3f",
      quote.TradePrice,
      quote.LastTradeTime,
      quote.ChangePrice, quote.ChangePercentage)

    fmt.Printf("\n")
  }

  fmt.Printf("\n")
}

func main() {
  ch1 := make(chan []Quote)

  quotes := flag.String("quotes", "GOOGL,TSLA", "stock symbols separate by quotes")
  flag.Parse()

  for {

    go func (msg chan []Quote) {
      test := make([]Quote, 0)
      Decode(getQuote(*quotes), &test)
      ch1 <- test
    }(ch1)

    result := <-ch1
    printQuote(result)
  }
}
