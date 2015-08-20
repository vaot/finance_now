package main

import "strconv"

type Quote struct {
  TradePrice       string `json:"l"`
  LastTradeTime    string `json:"ltt"`
  ChangePrice      string `json:"c"`
  ChangePercentage string `json:"cp"`
  Symbol           string `json:"t"`
}

type Quotes []Quote

func (quote *Quote) getTradePrice() float64 {
  f,_ := strconv.ParseFloat(quote.TradePrice, 32)
  return f
}

func (quote *Quote) getChangePrice() float64 {
  f,_ := strconv.ParseFloat(quote.ChangePrice, 32)
  return f
}

func (quote *Quote) getChangePricePercentage() float64 {
  f,_ := strconv.ParseFloat(quote.ChangePercentage, 32)
  return f
}
