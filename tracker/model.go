package tracker

type StockData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Change float64 `json:"change"`
	Time   string  `json:"time"`
}
