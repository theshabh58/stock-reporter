package apis

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//GetStockData fetches stock data from Tiingo stock API
//https://cloud.iexapis.com/stable/stock
func GetStockData(ticker string) (StockData, error) {
	var data StockData
	iexAPIURL := "https://cloud.iexapis.com"
	apiRoute := "/stable/stock/"
	reqType := "/quote?"
	apiToken := os.Getenv("IEX_API_TOKEN")
	queryParams := "&displayPercent=true"
	url := iexAPIURL + apiRoute + ticker + reqType + apiToken + queryParams

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error retriving data from iexCloud %v", err)
		return data, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return data, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Fatalln(err)
		return data, err
	}

	return data, nil
}
