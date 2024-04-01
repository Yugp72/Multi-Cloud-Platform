package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Items struct {
	Attributes map[string]string `json:"attributes"`
}

func createTableItemHandler(w http.ResponseWriter, r *http.Request) {
	var items Items
	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	av := make(map[string]*dynamodb.AttributeValue)
	for key, value := range items.Attributes {
		av[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Movies"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create item: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Item created successfully")
}

// package main

// import (
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"

// 	"fmt"
// 	"log"
// )

// func mainCreateTable() {

// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	// Create DynamoDB client
// 	svc := dynamodb.New(sess)

// 	tableName := "Movies"

// 	input := &dynamodb.CreateTableInput{
// 		AttributeDefinitions: []*dynamodb.AttributeDefinition{
// 			{
// 				AttributeName: aws.String("Year"),
// 				AttributeType: aws.String("N"),
// 			},
// 			{
// 				AttributeName: aws.String("Title"),
// 				AttributeType: aws.String("S"),
// 			},
// 		},
// 		KeySchema: []*dynamodb.KeySchemaElement{
// 			{
// 				AttributeName: aws.String("Year"),
// 				KeyType:       aws.String("HASH"),
// 			},
// 			{
// 				AttributeName: aws.String("Title"),
// 				KeyType:       aws.String("RANGE"),
// 			},
// 		},
// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
// 			ReadCapacityUnits:  aws.Int64(10),
// 			WriteCapacityUnits: aws.Int64(10),
// 		},
// 		TableName: aws.String(tableName),
// 	}

// 	_, err := svc.CreateTable(input)
// 	if err != nil {
// 		log.Fatalf("Got error calling CreateTable: %s", err)
// 	}

// 	fmt.Println("Created the table", tableName)

// }
