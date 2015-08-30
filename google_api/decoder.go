package google_api

import(
  "strconv"
  "encoding/json"
)

type Quote struct {
  TradePrice       string `json:"l"`
  LastTradeTime    string `json:"ltt"`
  ChangePrice      string `json:"c"`
  ChangePercentage string `json:"cp"`
  Symbol           string `json:"t"`
}

type Quotes []Quote

func (quote *Quote) GetTradePrice() float64 {
  f,_ := strconv.ParseFloat(quote.TradePrice, 32)
  return f
}

func (quote *Quote) GetChangePrice() float64 {
  f,_ := strconv.ParseFloat(quote.ChangePrice, 32)
  return f
}

func (quote *Quote) GetChangePricePercentage() float64 {
  f,_ := strconv.ParseFloat(quote.ChangePercentage, 32)
  return f
}


func Decode(resp string, parser *Quotes) {
  json.Unmarshal([]byte(resp), &parser)
}
