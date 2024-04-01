package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	tableName := params["tableName"]
	key := make(map[string]*dynamodb.AttributeValue)
	for attr, val := range params {
		if attr != "tableName" {
			key[attr] = &dynamodb.AttributeValue{
				S: aws.String(val),
			}
		}
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	input := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}
	_, err = svc.DeleteItem(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete item: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted item from table %s", tableName)
}

// package main

// // snippet-start:[dynamodb.go.delete_item.imports]
// import (
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"

// 	"fmt"
// 	"log"
// )

// // snippet-end:[dynamodb.go.delete_item.imports]

// func mainDeleteItem() {

// 	// and region from the shared configuration file ~/.aws/config.
// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	// Create DynamoDB client
// 	svc := dynamodb.New(sess)

// 	// snippet-start:[dynamodb.go.delete_item.call]
// 	tableName := "Movies"
// 	movieName := "The Big New Movie"
// 	movieYear := "2015"

// 	input := &dynamodb.DeleteItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"Year": {
// 				N: aws.String(movieYear),
// 			},
// 			"Title": {
// 				S: aws.String(movieName),
// 			},
// 		},
// 		TableName: aws.String(tableName),
// 	}

// 	_, err := svc.DeleteItem(input)
// 	if err != nil {
// 		log.Fatalf("Got error calling DeleteItem: %s", err)
// 	}

// 	fmt.Println("Deleted '" + movieName + "' (" + movieYear + ") from table " + tableName)
// 	// snippet-end:[dynamodb.go.delete_item.call]
// }
