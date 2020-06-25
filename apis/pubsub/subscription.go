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

	//consume 10 messages
	var mutex sync.Mutex
	recievedMessages := 10
	subscribe := client.Subscription(subscriberID)

	cctx, cancel := context.WithCancel(ctx)
	err = subscribe.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mutex.Lock()
		defer mutex.Unlock()
		fmt.Fprintf(w, "Got message: %q \n", string(msg.Data))
		msg.Ack()
		recievedMessages++
		if recievedMessages == 10 {
			cancel()
		}
	})

	if err != nil {
		return fmt.Errorf("Error recieving messages %v", err)
	}

	return nil
}
