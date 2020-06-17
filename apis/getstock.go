package apis

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

//GetStockReports gets the stock reports from the database and publishes them to a topic
func GetStockReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get required environment variables/constants
	const collectionName = "stock-reports"
	const topicID = "fetch-stock-data"
	projectID := os.Getenv("PROJECT_ID")

	//Connect to the database
	ctx := context.Background()
	db, err := ConnectToDatabase(ctx, projectID)
	defer db.Client.Close()

	if err != nil {
		log.Println("Error connecting to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Get all the stock reports
	stockReports, err := db.GetStockReports(ctx, collectionName)

	if err != nil {
		log.Println("Error retrieving stock reports:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Create client to connect to pubsub topic
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println("Error creating pubsub client:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := client.Topic(topicID)

	//Publish message for eacb stock report
	for _, report := range stockReports {
		//Convert message into bytes
		msg := new(bytes.Buffer)
		json.NewEncoder(msg).Encode(report)

		//Publish message
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte(msg.Bytes()),
		})

		//Verify successful publishing of message
		id, err := result.Get(ctx)
		if err != nil {
			log.Println("Error publishing message:", err)
		}

		log.Println("Successfully published message. ID:", id)
	}

	w.WriteHeader(http.StatusCreated)
}
