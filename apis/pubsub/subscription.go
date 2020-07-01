package pubsub

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

//PullMessageAsync pull messages from a subscribed topic
func PullMessageAsync(w io.Writer, projectID, subscriberID string) error {
	ctx := context.Background()
	client, err := GetPubSubClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to create new pubsub client %v", err)
	}

	sub := client.Subscription(subscriberID)

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Printf("Message Recieved: %s\n", m.Data)
		m.Ack()
	})

	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
