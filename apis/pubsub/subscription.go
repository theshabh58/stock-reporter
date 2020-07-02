package pubsub

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

//PullStockPriceMessages pull messages from a subscribed topic
func PullStockPriceMessages(w io.Writer, projectID, subscriberID string) error {
	ctx := context.Background()
	client, err := GetPubSubClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to create new pubsub client %v", err)
	}

	sub := client.Subscription(subscriberID)

	cctx, cancel := context.WithTimeout(ctx, 20000*time.Millisecond)
	defer cancel()
	err = sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Printf("Message Recieved: %s\n", m.Data)
		m.Ack()
	})

	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
