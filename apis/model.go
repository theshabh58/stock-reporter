package apis

//StockReport model
type StockReport struct {
	ID     string
	User   User    `json:"user" firestore:"user"`
	Stocks []Stock `json"stocks" firestore"stocks"`
}

//Stock model
type Stock struct {
	Ticker string `json:"ticker" firestore:"ticker"`
}

//User model
type User struct {
	FirstName string `json:"firstName" firestore:"firstName"`
	LastName  string `json:"lastName" firestore:"lastName"`
	Email     string `json:"email" firestore"email"`
}
