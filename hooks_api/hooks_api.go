package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/fzzy/radix/redis"
  "strings"
)


func main() {
  m := martini.Classic()
  client, _ := redis.Dial("tcp", "localhost:6379")

  m.Get("/alerts", func(r *http.Request) string {
    qs := r.URL.Query()

    err := client.Cmd("HSET", "alerts", strings.ToUpper(qs.Get("quote")), qs.Get("action")).Err

    if err != nil {
      return "We failed to stop your alert"
    }

    return ("We just stopped your " + qs.Get("quote") + " alert.")
  })

  m.RunOnAddr(":6000")
}
