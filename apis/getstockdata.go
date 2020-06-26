package apis

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//GetStockData fetches stock data from Tiingo stock API
func GetStockData(ticker string) []interface{} {
	tiingoAPIURL := "https://api.tiingo.com/tiingo/daily/"
	apiRoute := "/prices?"

	resp, err := http.Get(tiingoAPIURL + ticker + apiRoute + os.Getenv("TIINGO_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error retriving data from stocks api %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var data []StockData

	err = json.Unmarshal(respBody, &data)
	data[0].StockTicker = ticker

	if err != nil {
		log.Fatalln(err)
	}

	stockData := []interface{}{data}

	return stockData
}
