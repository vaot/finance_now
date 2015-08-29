package google_api

import (
  "strings"
  "net/http"
  "io/ioutil"
  "log"
)

const GOOGLE_API_URL string = "http://finance.google.com/finance/info"

func BuildUrl(quote string) string {
  url := []string{ GOOGLE_API_URL, "?", "q=", quote }
  return strings.Join(url, "")
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

  return strings.Replace(parsedResp, "//", "", -1)
}
