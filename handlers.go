package main

import (
  "strconv"
  "strings"
  "net/http"
  "os"
  "encoding/json"
  "bytes"
)

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
