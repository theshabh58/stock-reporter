package apis

import (
	"context"

	"cloud.google.com/go/firestore"
)

//Database struct for database operationss
type Database struct {
	Client *firestore.Client
}

//ConnectToDatabase attempts to the database returns the database client
func ConnectToDatabase(ctx context.Context, projectID string) (*Database, error) {
	dbClient, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		return nil, err
	}

	return &Database{Client: dbClient}, nil
}

//InsertStockReport inserts a new stock report into the database
func (db *Database) InsertStockReport(ctx context.Context, report StockReport, collectionName string) (string, error) {
	//Add post to the database
	result, _, err := db.Client.Collection(collectionName).Add(ctx, report)

	if err != nil {
		return "", err
	}

	return result.ID, nil
}
