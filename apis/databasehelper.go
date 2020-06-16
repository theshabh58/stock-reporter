package apis

import (
	"context"
	"fmt"

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

//GetStockReportByEmail find stock report for the given email
func (db *Database) GetStockReportByEmail(ctx context.Context, email, collectionName string) (string, error) {
	//Search for user by given email
	it := db.Client.Collection(collectionName).Where("user.email", "==", email).Documents(ctx)
	defer it.Stop()

	doc, err := it.Next()
	if err != nil {
		return "", fmt.Errorf("No matching stock report found: %v", err)
	}

	return doc.Ref.ID, nil
}

//AddNewStock adds a new stock to the stock report with the matching id
func (db *Database) AddNewStock(ctx context.Context, stock, id, collectionName string) error {
	//Get the stock report
	dsnap, err := db.Client.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		return err
	}
	var doc StockReport
	dsnap.DataTo(&doc)

	//Update the stocks in the stock report by adding the new stock
	_, err = db.Client.Collection(collectionName).Doc(id).Set(ctx, map[string]interface{}{
		"stocks": append(doc.Stocks, Stock{Ticker: stock}),
	}, firestore.MergeAll)

	if err != nil {
		return err
	}

	return nil
}
