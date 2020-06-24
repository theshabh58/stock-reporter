package apis

import "time"

//StockReport model
type StockReport struct {
	ID     string  `json:"id" firestore:"-"`
	User   User    `json:"user" firestore:"user"`
	Stocks []Stock `json:"stocks" firestore:"stocks"`
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
	StockTicker  string    `json:"stockName" firestore:"stockName"`
	LowestPrice  float64   `json:"low" firestore:"lowestPrice"`
	HighestPrice float64   `json:"high" firestore:"highestPrice"`
	OpeningPrice float64   `json:"open" firestore:"openingPrice"`
	ClosingPrice float64   `json:"close" firestore:"closingPrice"`
	Date         time.Time `json:"date" firestore:"date"`
}
