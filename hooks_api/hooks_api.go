package main

import (
	"github.com/fzzy/radix/redis"
	"github.com/go-martini/martini"
	"net/http"
	"strings"
	// "encoding/json"
	"fmt"
	"github.com/vaot/finance_now/google_api"
	"os"
  "net/url"
)

func alertsHandler(r *http.Request) string {
  redisUrl, _ := url.Parse(os.Getenv("REDIS_URL"))
	client, _ := redis.Dial("tcp", redisUrl.Host)

	qs := r.URL.Query()
	fmt.Println(qs.Get("text"))
	err := client.Cmd("HSET", "alerts", strings.ToUpper(qs.Get("text")), qs.Get("action")).Err

	if err != nil {
		return "We failed to stop your alert"
	}

	return ("We just stopped your " + qs.Get("text") + " alert.")
}

func quotesHandler(r *http.Request) string {
	qs := r.URL.Query()
	quote := qs.Get("text")

	var quotes google_api.Quotes
	google_api.Decode(google_api.GetQuote(quote), &quotes)

	report := []string{}

	for _, quote := range quotes {
		message := quote.Symbol + " " + quote.TradePrice
		report = append(report, message)
	}

	return strings.Join(report, ", ")
}

func main() {
	m := martini.Classic()

	m.Get("/alerts", alertsHandler)
	m.Get("/quotes", quotesHandler)

	m.RunOnAddr(":" + os.Getenv("PORT"))
}
