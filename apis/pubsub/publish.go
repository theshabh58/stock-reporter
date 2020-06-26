package pubsub

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

//PublishStockData a message on to a topic
func PublishStockData(wr io.Writer, projectID, topicID string, msg string) error {

	if msg != "" {
		ctx := context.Background()
		pubsubclient, err := GetPubSubClient(ctx, projectID)
		if err != nil {
			return err
		}

		topic := pubsubclient.Topic(topicID)
		defer topic.Stop()

		var results []*pubsub.PublishResult

		r := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(msg),
		})

		results = append(results, r)

		for _, r := range results {
			id, err := r.Get(ctx)
			if err != nil {
				return fmt.Errorf("Error getting id %v", err)
			}
			fmt.Printf("Published a message to topic %s, with messageID: %v\n", topicID, id)
		}
	}

	return nil

}
