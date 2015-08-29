package yahoo_api

import (
  "strings"
  "net/http"
  "io/ioutil"
  "log"
  "net/url"
)

const(
  YAHOO_API_URL = "https://query.yahooapis.com/v1/public/yql"
)

func BuildUrl(quote string) string {

  var queryUrl *url.URL
  queryUrl, err := url.Parse(YAHOO_API_URL)

  if err != nil {
    panic("Cannot parse url")
  }

  parameters := url.Values{}
  parameters.Add("q", strings.Join([]string{"select * from yahoo.finance.quotes where symbol in (", quote, ")"}, ""))
  parameters.Add("format", "json")
  parameters.Add("env", "http://datatables.org/alltables.env")
  queryUrl.RawQuery = parameters.Encode()

  return queryUrl.String()
}

func GetQuote(quote string) string {

  resp, err := http.Get(BuildUrl(quote))

  if err != nil {
    log.Fatal("Could not get quote")
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  parsedResp := string(body)

  if err != nil {
    log.Fatal("Cannot read body")
  }

  return parsedResp
}
