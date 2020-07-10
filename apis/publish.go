package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

//PubSubMessage payload for the event
type PubSubMessage struct {
	Data []byte `json:"data"`
}

//PublishStockData used to publish stock data to topic
func PublishStockData(ctx context.Context, msg PubSubMessage) error {

	var report StockReport

	err := json.Unmarshal(msg.Data, &report)
	if err != nil {
		log.Fatalf("error unmarshaling payload from topic: %v", err)
		return err
	}

	//Fetch stock data from api
	for _, v := range report.Stocks {
		data, err := GetStockData(v.Ticker)
		if err != nil {
			log.Fatalf("error fetching stock data from tiingo api: %v", err)
			return err
		}
		report.StockReport = append(report.StockReport, data...)
	}

	//publish stock data to topic
	err = publishStockDataToTopic(ctx, report)
	if err != nil {
		return err
	}
	return nil
}

func publishStockDataToTopic(ctx context.Context, report StockReport) error {
	//publish a message to topic
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return err
	}

	topic := client.Topic(os.Getenv("TOPIC_ID"))
	defer topic.Stop()

	var results []*pubsub.PublishResult

	message, err := json.Marshal(report)
	if err != nil {
		return err
	}

	r := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})

	results = append(results, r)

	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			log.Fatalf("Error getting id %v", err)
			return err
		}
		fmt.Printf("Published a message to topic with messageID: %v\n", id)
	}
	return nil
}
