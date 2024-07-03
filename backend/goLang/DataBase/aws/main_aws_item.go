package aws_dynamodb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db "btep.project/databaseConnection"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// ItemRequest represents the JSON request structure for DynamoDB item operations
type ItemRequest struct {
	TableName string                 `json:"tableName"`
	Key       map[string]interface{} `json:"key"`
	Region    string                 `json:"region"`
	AccountID int                    `json:"accountID"`
}

// ItemResponse represents the JSON response structure for DynamoDB item operations
type ItemResponse struct {
	Message string `json:"message"`
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var req ItemRequest
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

	item := make(map[string]*dynamodb.AttributeValue)
	for k, v := range req.Key {
		// Check if the value is a number
		var attrValue *dynamodb.AttributeValue
		switch v := v.(type) {
		case int, int8, int16, int32, int64:
			attrValue = &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%v", v)),
			}
		case float32, float64:
			attrValue = &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%f", v)),
			}
		default:
			attrValue = &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%v", v)),
			}
		}
		item[k] = attrValue
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(req.TableName),
		Item:      item,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating item in DynamoDB: %v", err)
		return
	}

	resp := ItemResponse{Message: "Item created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ItemQuery struct {
	TableName string                 `json:"tableName"`
	QueryKey  map[string]interface{} `json:"queryKey"`
	Region    string                 `json:"region"`
	AccountID int                    `json:"accountID"`
}

func ReadItemHandler(w http.ResponseWriter, r *http.Request) {
	var req ItemQuery
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

	// Build the filter expression and expression attribute values dynamically
	var conditions []string
	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)
	for key, value := range req.QueryKey {
		conditions = append(conditions, fmt.Sprintf("%s = :%s", key, key))

		var attrValue *dynamodb.AttributeValue

		switch v := value.(type) {
		case int, int8, int16, int32, int64:
			attrValue = &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%v", v)),
			}
		case float32, float64:
			attrValue = &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%f", v)),
			}
		default:
			attrValue = &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%v", v)),
			}
		}

		expressionAttributeValues[":"+key] = attrValue
	}

	filterExpression := strings.Join(conditions, " AND ")

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(req.TableName),
		FilterExpression:          aws.String(filterExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}
	fmt.Println("input: ", input)

	result, err := svc.Scan(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error querying items from DynamoDB: %v", err)
		return
	}
	fmt.Println("result: ", result)
	// Convert the result to JSON and send it in the response
	resultJSON, err := json.Marshal(result.Items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding JSON response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}

type updateRequest struct {
	TableName string                 `json:"tableName"`
	QueryKey  map[string]interface{} `json:"queryKey"`
	UpdateMap map[string]interface{} `json:"updateMap"`
	Region    string                 `json:"region"`
	AccountID int                    `json:"accountID"`
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	var req updateRequest
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

	// Step 1: Scan for items based on the query key
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(req.TableName),
	}
	scanInput.FilterExpression, scanInput.ExpressionAttributeValues = buildFilterExpression(req.QueryKey)

	scanResult, err := svc.Scan(scanInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error scanning items from DynamoDB: %v", err)
		return
	}

	// Step 2: Update the matching items with the update map
	for _, item := range scanResult.Items {
		updateInput := &dynamodb.UpdateItemInput{
			TableName:                 aws.String(req.TableName),
			Key:                       item, // Use the item as the key for the update
			UpdateExpression:          aws.String(buildUpdateExpression(req.UpdateMap)),
			ExpressionAttributeValues: buildExpressionAttributeValues(req.UpdateMap),
		}

		_, err = svc.UpdateItem(updateInput)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error updating item in DynamoDB: %v", err)
			return
		}
	}

	resp := ItemResponse{Message: "Items updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Helper function to build the filter expression and expression attribute values based on the query key
// Helper function to build the filter expression and expression attribute values based on the query key
func buildFilterExpression(queryKey map[string]interface{}) (*string, map[string]*dynamodb.AttributeValue) {
	var expressions []string
	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)

	for key, value := range queryKey {
		expressions = append(expressions, fmt.Sprintf("%s = :%s", key, key))
		expressionAttributeValues[":"+key] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%v", value)),
		}
	}

	filterExpression := strings.Join(expressions, " AND ")
	return aws.String(filterExpression), expressionAttributeValues
}

// Helper function to build the update expression based on the update map
func buildUpdateExpression(updateMap map[string]interface{}) string {
	var expressions []string
	for key := range updateMap {
		expressions = append(expressions, fmt.Sprintf("%s = :%s", key, key))
	}
	return "SET " + strings.Join(expressions, ", ")
}

// Helper function to build the expression attribute values based on the update map
func buildExpressionAttributeValues(updateMap map[string]interface{}) map[string]*dynamodb.AttributeValue {
	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)
	for key, value := range updateMap {
		expressionAttributeValues[":"+key] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%v", value)), // Assuming all values are strings for simplicity
		}
	}
	return expressionAttributeValues
}

// DeleteItemHandler handles POST requests to delete an item from DynamoDB
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	var req ItemRequest
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

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(req.TableName),
		Key:       buildKey(req.Key),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting item from DynamoDB: %v", err)
		return
	}

	resp := ItemResponse{Message: "Item deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Helper function to build DynamoDB Key from input
func buildKey(input map[string]interface{}) map[string]*dynamodb.AttributeValue {
	key := make(map[string]*dynamodb.AttributeValue)
	for k, v := range input {
		key[k] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%v", v)),
		}
	}
	return key
}

type ListItemsQuery struct {
	TableName string `json:"tableName"`
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

func ListItemsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListItemsQuery
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting cloud account details: %v", err), http.StatusInternalServerError)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		http.Error(w, "Error initializing AWS session", http.StatusInternalServerError)
		return
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String(req.TableName),
	}
	fmt.Println("input: ", input)

	result, err := svc.Scan(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying items from DynamoDB: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the result to JSON and send it in the response
	resultJSON, err := json.Marshal(result.Items)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}
