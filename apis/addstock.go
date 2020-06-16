package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

//AddStockRequest request body for AddStock endpoint
type AddStockRequest struct {
	UserEmail   string `json:"userEmail"`
	StockTicker string `json:"stockTicker"`
}

//AddStock adds a new stock to the given user
func AddStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get request body
	decoder := json.NewDecoder(r.Body)
	var input AddStockRequest
	err := decoder.Decode(&input)
	if err != nil {
		log.Println("Invalid request body provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Verify all properties in the request body are provided
	if input.UserEmail == "" || input.StockTicker == "" {
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

	//Find the matching stock report using the given email
	id, err := db.GetStockReportByEmail(ctx, input.UserEmail, collectionName)

	if err != nil {
		log.Println("Error finding stock report with the provided email:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update that stock report by adding the new stock
	err = db.AddNewStock(ctx, input.StockTicker, id, collectionName)

	if err != nil {
		log.Println("Error adding the provided stock:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
