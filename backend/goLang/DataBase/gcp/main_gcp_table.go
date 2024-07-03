// ListTablesHandler handles GET requests to list collections (tables) in Firestore
package gcp_firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"cloud.google.com/go/firestore"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type TableRequest struct {
	TableName  string              `json:"tableName"`
	ProjectID  string              `json:"projectID"`
	AccountID  int                 `json:"accountID"`
	Token      string              `json:"token"`
	Attributes []map[string]string `json:"attributes,omitempty"`
}

type TableResponse struct {
	Message string `json:"message"`
}

// CreateFirestoreClient creates a Firestore client with OAuth 2.0 token authentication
func CreateFirestoreClient(ctx context.Context, token, projectID string) (*firestore.Client, error) {
	// Create Firestore client with OAuth 2.0 token authentication
	client, err := firestore.NewClient(ctx, projectID, option.WithTokenSource(oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateTableHandler handles POST requests to create a collection (table) in Firestore
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	var req TableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}
	fmt.Println("decoded request", req)

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	fmt.Println("cloudAccount", cloudAccount)
	// Create Firestore client with OAuth 2.0 token authentication
	ctx := context.Background()
	client, err := CreateFirestoreClient(ctx, req.Token, cloudAccount.ProjectID.String)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Firestore client: %v", err)
		return
	}
	defer client.Close()

	fmt.Println("client firebase", client)

	// Create the collection (table) with the given name
	_, err = client.Collection(req.TableName).Doc("placeholder").Set(ctx, map[string]interface{}{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating collection in Firestore: %v", err)
		return
	}

	fmt.Println("collection created")

	resp := TableResponse{Message: "Collection created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateTableHandler handles POST requests to update documents in a collection (table) in Firestore
func UpdateTableHandler(w http.ResponseWriter, r *http.Request) {
	var req TableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create Firestore client with OAuth 2.0 token authentication
	ctx := context.Background()
	client, err := CreateFirestoreClient(ctx, req.Token, cloudAccount.ProjectID.String)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Firestore client: %v", err)
		return
	}
	defer client.Close()

	// Construct a query to filter documents in the collection based on certain criteria
	query := client.Collection(req.TableName).Where("fieldName", "==", "fieldValue")

	// Iterate over the documents returned by the query and update them
	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		fmt.Println(doc.Data())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error iterating over documents: %v", err)
			return
		}

		// Update the document fields based on the update request
		// Example: doc.Ref.Set(ctx, map[string]interface{}{"fieldName": "updatedValue"})
	}

	resp := TableResponse{Message: "Documents updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteTableHandler handles POST requests to delete a collection (table) from Firestore
func DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	var req TableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create Firestore client with OAuth 2.0 token authentication
	ctx := context.Background()
	client, err := CreateFirestoreClient(ctx, req.Token, cloudAccount.ProjectID.String)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Firestore client: %v", err)
		return
	}
	defer client.Close()

	// Fetch all documents in the collection (table)
	iter := client.Collection(req.TableName).Documents(ctx)
	batchSize := 100 // Set the batch size to avoid memory issues for large collections
	batch := client.Batch()
	deletedCount := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error iterating over documents: %v", err)
			return
		}

		// Delete each document in the collection
		batch.Delete(doc.Ref)
		// Execute the batch delete operation in batches
		deletedCount++
		if deletedCount%batchSize == 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error committing batch delete: %v", err)
				return
			}
			batch = client.Batch()
		}
	}

	// Commit any remaining documents in the batch
	_, err = batch.Commit(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error committing batch delete: %v", err)
		return
	}

	resp := TableResponse{Message: "Collection deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ListTablesHandler(w http.ResponseWriter, r *http.Request) {
	var req TableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create Firestore client with OAuth 2.0 token authentication
	ctx := context.Background()
	client, err := CreateFirestoreClient(ctx, req.Token, cloudAccount.ProjectID.String)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Firestore client: %v", err)
		return
	}
	defer client.Close()

	// Retrieve a list of all collections (tables) in the Firestore database
	collections, err := client.Collections(ctx).GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving collections from Firestore: %v", err)
		return
	}

	// Extract collection names from the collection references
	var tableNames []string
	for _, colRef := range collections {
		tableNames = append(tableNames, colRef.ID)
	}

	// Send the list of table names in the response
	resp := TableResponse{Message: "Collections listed successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// func ListTablesHandler(w http.ResponseWriter, r *http.Request) {
// 	var req TableRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, "Invalid request body")
// 		return
// 	}

// 	// Fetch cloud account details from the database
// 	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
// 		return
// 	}

// 	// Create Firestore client with OAuth 2.0 token authentication
// 	ctx := context.Background()
// 	client, err := CreateFirestoreClient(ctx, req.Token, cloudAccount.ProjectID.String)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "Error creating Firestore client: %v", err)
// 		return
// 	}
// 	defer client.Close()

// 	// Retrieve a list of all collections (tables) in the Firestore database
// 	collections, err := client.Collections(ctx).GetAll()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "Error retrieving collections from Firestore: %v", err)
// 		return
// 	}

// 	// Extract collection names from the collection references
// 	var tableNames []string
// 	for _, colRef := range collections {
// 		tableNames = append(tableNames, colRef.ID)
// 	}

// 	// Send the list of table names in the response
// 	resp := TableResponse{Message: "Collections listed successfully"}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(resp)
// }
