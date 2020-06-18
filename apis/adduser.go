package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

//AddUser adds a new user for stock reports
func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get request body
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		log.Println("Invalid request body provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Verify all properties in the request body are provided
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		log.Println("Missing data in the request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Get required environment variables/constants
	const collectionName = "stock-reports"
	projectID := os.Getenv("PROJECT_ID")

	//Connect to the database
	ctx := r.Context()
	db, err := ConnectToDatabase(ctx, projectID)
	defer db.Client.Close()

	if err != nil {
		log.Println("Error connecting to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Build sock report to insert
	stockReport := StockReport{
		User:   user,
		Stocks: []Stock{},
	}

	//Insert into the database
	id, err := db.InsertStockReport(ctx, stockReport, collectionName)

	if err != nil {
		log.Println("Error inserting post into the database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Proceed with returning response if no errors encountered
	stockReport.ID = id
	log.Println("Successfully inserted new stock report:", id)

	json.NewEncoder(w).Encode(stockReport)
	w.WriteHeader(http.StatusCreated)
}
