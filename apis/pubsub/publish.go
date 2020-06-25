package pubsub

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"github.com/theshabh58/stock-reporter/rest-apis/apis"
)

//Publish a message on to a topic
func Publish(wr io.Writer, projectID, topicID string, msg []apis.StockData) error {

	ctx := context.Background()
	pubsubclient, err := GetPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	topic := pubsubclient.Topic(topicID)
	defer topic.Stop()

	var results []*pubsub.PublishResult

	r := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(fmt.Sprintf("%v", msg)),
	})

	results = append(results, r)

	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			return fmt.Errorf("Error getting id %v", err)
		}
		fmt.Printf("Published a message to topic %s, with messageID: %v\n", topicID, id)
	}

	return nil

}
