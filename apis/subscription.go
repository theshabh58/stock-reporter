package apis

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

//PullStockPriceMessages pull messages from a subscribed topic
func PullStockPriceMessages(ctx context.Context) error {

	client, err := GetPubSubClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return fmt.Errorf("Failed to create new pubsub client %v", err)
	}

	sub := client.Subscription(os.Getenv("SUBSCRIPTION_ID"))

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
