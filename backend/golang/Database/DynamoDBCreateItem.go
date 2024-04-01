package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Year   int
	Title  string
	Plot   string
	Rating float64
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tableName := "XYZ"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	year := strconv.Itoa(item.Year)

	message := fmt.Sprintf("Successfully added '%s' (%s) to table %s", item.Title, year, tableName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, message)
}

func main() {
	http.HandleFunc("/create-item", createItemHandler)
	http.HandleFunc("/create-table", createTableItemHandler)
	http.HandleFunc("/delete-item", deleteItemHandler)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// package main

// import (
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

// 	"fmt"
// 	"log"
// 	"strconv"
// )

// // snippet-end:[dynamodb.go.create_item.struct]
// type Item struct {
// 	Year   int
// 	Title  string
// 	Plot   string
// 	Rating float64
// }

// func main() {

// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	// Create DynamoDB client
// 	svc := dynamodb.New(sess)

// 	// snippet-start:[dynamodb.go.create_item.assign_struct]
// 	item := Item{
// 		Year:   2015,
// 		Title:  "The Big New Movie",
// 		Plot:   "Nothing happens at all.",
// 		Rating: 0.0,
// 	}

// 	av, err := dynamodbattribute.MarshalMap(item)
// 	if err != nil {
// 		log.Fatalf("Got error marshalling new movie item: %s", err)
// 	}

// 	tableName := "Movies"

// 	input := &dynamodb.PutItemInput{
// 		Item:      av,
// 		TableName: aws.String(tableName),
// 	}

// 	_, err = svc.PutItem(input)
// 	if err != nil {
// 		log.Fatalf("Got error calling PutItem: %s", err)
// 	}

// 	year := strconv.Itoa(item.Year)

// 	fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)
// 	// snippet-end:[dynamodb.go.create_item.call]
// }
