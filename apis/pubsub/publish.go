package pubsub

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/pubsub"
)

//Publish a message on to a topic
func publish(wr io.Writer, projectID, topicID string, msg []string) error {

	ctx := context.Background()
	pubsubclient, err := GetPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	topic := pubsubclient.Topic(topicID)
	defer topic.Stop()

	var results []*pubsub.PublishResult

	r := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("Test"),
	})

	results = append(results, r)

	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			log.Fatalf("Error getting id %v", err)
		}
		fmt.Printf("Published a message to topic %s, with messageID: %v", topicID, id)
	}

	return nil

}
