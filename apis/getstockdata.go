package apis

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//GetStockData fetches stock data from Tiingo stock API
func GetStockData(ticker string) ([]StockData, error) {
	tiingoAPIURL := "https://api.tiingo.com/tiingo/daily/"
	apiRoute := "/prices?"
	apiToken := os.Getenv("TIINGO_API_TOKEN")

	resp, err := http.Get(tiingoAPIURL + ticker + apiRoute + apiToken)
	if err != nil {
		log.Fatalf("Error retriving data from stocks api %v", err)
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	var data []StockData

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	data[0].Ticker = ticker

	return data, nil
}
