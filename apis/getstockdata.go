package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//GetStocksData fetches stock data from Tiingo stock API
func GetStocksData(ticker string) []StockData {
	resp, err := http.Get("https://api.tiingo.com/tiingo/daily/" + ticker + "/prices?token=c4e60a6a196fa1b481aa7366430393849394d1fa")
	if err != nil {
		log.Fatalf("Error retriving data from stocks api %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Println(string(respBody))

	var data []StockData

	err = json.Unmarshal(respBody, &data)
	data[0].StockTicker = ticker

	if err != nil {
		log.Fatalln(err)
	}

	return data
}
