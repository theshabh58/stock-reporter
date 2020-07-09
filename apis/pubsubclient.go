package apis

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

//GetPubSubClient creates a pubsub client
func GetPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	return client, nil
}
