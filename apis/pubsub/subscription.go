package pubsub

import (
	"context"
	"fmt"
	"io"
	"sync"

	"cloud.google.com/go/pubsub"
)

//PullMessageAsync pull messages from a subscribed topic
func PullMessageAsync(w io.Writer, projectID, subscriberID string) error {
	ctx := context.Background()
	client, err := GetPubSubClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to create new pubsub client %v", err)
	}

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subscriberID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
