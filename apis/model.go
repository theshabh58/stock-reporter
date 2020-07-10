package apis

import "time"

//StockReport model
type StockReport struct {
	ID          string      `json:"id" firestore:"-"`
	User        User        `json:"user" firestore:"user"`
	Stocks      []Stock     `json:"stocks" firestore:"stocks"`
	StockReport []StockData `json:"stockReport" firestore:"-"`
}

//Stock model
type Stock struct {
	Ticker string `json:"ticker" firestore:"ticker"`
}

//User model
type User struct {
	FirstName string `json:"firstName" firestore:"firstName"`
	LastName  string `json:"lastName" firestore:"lastName"`
	Email     string `json:"email" firestore:"email"`
}

//StockData model
type StockData struct {
	Ticker       string    `json:"ticker"`
	LowestPrice  float64   `json:"low"`
	HighestPrice float64   `json:"high"`
	OpeningPrice float64   `json:"open"`
	ClosingPrice float64   `json:"close"`
	Date         time.Time `json:"date"`
}
