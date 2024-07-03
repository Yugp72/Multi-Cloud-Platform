package aws_dynamodb

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type TableRequest struct {
	TableName             string              `json:"tableName"`
	Region                string              `json:"region"`
	AccountID             int                 `json:"accountID"`
	AttributeDefinitions  []map[string]string `json:"attributeDefinitions"`
	KeySchema             []map[string]string `json:"keySchema"`
	ProvisionedThroughput map[string]int64    `json:"provisionedThroughput"`
}

type TableResponse struct {
	Message string `json:"message"`
}

// CreateTableHandler handles POST requests to create a table in DynamoDB
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	var req TableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	svc := dynamodb.New(sess)
	fmt.Println(req.AttributeDefinitions)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: buildAttributeDefinitions(req.AttributeDefinitions),
		KeySchema:            buildKeySchema(req.KeySchema),
		TableName:            aws.String(req.TableName),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(req.ProvisionedThroughput["ReadCapacityUnits"]),
			WriteCapacityUnits: aws.Int64(req.ProvisionedThroughput["WriteCapacityUnits"]),
		},
	}

	_, err = svc.CreateTable(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating table in DynamoDB: %v", err)
		return
	}

	resp := TableResponse{Message: "Table created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type DeleteTableRequest struct {
	TableName string `json:"tableName"`
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

// DeleteTableHandler handles POST requests to delete a table from DynamoDB
func DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteTableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(req.TableName),
	}

	_, err = svc.DeleteTable(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting table from DynamoDB: %v", err)
		return
	}

	resp := TableResponse{Message: "Table deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type UpdateTableRequest struct {
	TableName             string              `json:"tableName"`
	Region                string              `json:"region"`
	AccountID             int                 `json:"accountID"`
	AttributeDefinitions  []map[string]string `json:"attributeDefinitions,omitempty"`
	ProvisionedThroughput map[string]int64    `json:"provisionedThroughput,omitempty"`
}

// UpdateTableHandler handles POST requests to update a table in DynamoDB
func UpdateTableHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateTableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateTableInput{
		TableName:            aws.String(req.TableName),
		AttributeDefinitions: buildAttributeDefinitions(req.AttributeDefinitions),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(req.ProvisionedThroughput["readCapacityUnits"]),
			WriteCapacityUnits: aws.Int64(req.ProvisionedThroughput["writeCapacityUnits"]),
		},
	}

	_, err = svc.UpdateTable(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating table in DynamoDB: %v", err)
		return
	}

	resp := TableResponse{Message: "Table updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Helper function to build DynamoDB AttributeDefinitions from input
func buildAttributeDefinitions(input []map[string]string) []*dynamodb.AttributeDefinition {
	var attributeDefinitions []*dynamodb.AttributeDefinition
	for _, attr := range input {
		attributeDefinitions = append(attributeDefinitions, &dynamodb.AttributeDefinition{
			AttributeName: aws.String(attr["attributeName"]),
			AttributeType: aws.String(attr["attributeType"]),
		})
	}
	return attributeDefinitions
}

// Helper function to build DynamoDB KeySchema from input
func buildKeySchema(input []map[string]string) []*dynamodb.KeySchemaElement {
	var keySchema []*dynamodb.KeySchemaElement
	for _, key := range input {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(key["attributeName"]),
			KeyType:       aws.String(key["keyType"]),
		})
	}
	return keySchema
}

type ListTableRequest struct {
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

func ListTablesHandler(w http.ResponseWriter, r *http.Request) {
	var req ListTableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get cloud account details
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create DynamoDB service client
	svc := dynamodb.New(sess)

	// Input for ListTables operation
	input := &dynamodb.ListTablesInput{}

	// ListTables operation to get a list of table names
	result, err := svc.ListTables(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing tables in DynamoDB: %v", err)
		return
	}

	// Initialize a slice to store table descriptions
	var tableDetails []*dynamodb.TableDescription

	// Iterate over the list of table names
	for _, tableName := range result.TableNames {
		// DescribeTable operation to get details of each table
		describeInput := &dynamodb.DescribeTableInput{
			TableName: tableName,
		}
		tableOutput, err := svc.DescribeTable(describeInput)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error describing table %s: %v", *tableName, err)
			return
		}

		// Append the table description to the slice
		tableDetails = append(tableDetails, tableOutput.Table)
	}

	// Send success response with the list of table details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tableDetails)
}
