package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"golang/db" // Replace "your-package-path" with the actual import path for the db package
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// BucketRequest represents the JSON request structure
type BucketRequest struct {
	BucketName string `json:"bucketName"`
	ProjectID  string `json:"projectID"`

	AccountID int `json:"accountID"`
}

// ObjectRequest represents the JSON request structure for object operations
type ObjectRequest struct {
	BucketName string `json:"bucketName"`
	ObjectKey  string `json:"objectKey"`
	Content    string `json:"content"`
	AccountID  int    `json:"accountID"`
}

// BucketResponse represents the JSON response structure
type BucketResponse struct {
	Message string `json:"message"`
}

// ObjectResponse represents the JSON response structure for object operations
type ObjectResponse struct {
	Message string `json:"message"`
}

// UploadObjectHandler handles POST requests to upload an object to a GCP bucket
func UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	var req ObjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get CloudAccount details
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	ctx := context.Background()

	// Authenticate and create a storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(*cloudAccount.CredentialsJSON)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCP storage client: %v", err)
		return
	}
	defer client.Close()

	bucketName := req.BucketName
	objectKey := req.ObjectKey
	content := req.Content

	// Convert content string to byte array
	contentBytes := []byte(content)

	wc := client.Bucket(bucketName).Object(objectKey).NewWriter(ctx)
	defer wc.Close()

	// Write content to the object
	if _, err := wc.Write(contentBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error uploading object to GCP: %v", err)
		return
	}

	// Send success response
	resp := ObjectResponse{Message: "Object uploaded successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetObjectHandler handles GET requests to retrieve an object from a GCP bucket
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	var req ObjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get CloudAccount details
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	ctx := context.Background()

	// Authenticate and create a storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(*cloudAccount.CredentialsJSON)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCP storage client: %v", err)
		return
	}
	defer client.Close()

	bucketName := req.BucketName
	objectKey := req.ObjectKey

	rc, err := client.Bucket(bucketName).Object(objectKey).NewReader(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error downloading object from GCP: %v", err)
		return
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading object content: %v", err)
		return
	}

	// Send the object content in the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+objectKey)
	w.Write(data)
}

// DeleteObjectHandler handles POST requests to delete an object from a GCP bucket
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	var req ObjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get CloudAccount details
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	ctx := context.Background()

	// Authenticate and create a storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(*cloudAccount.CredentialsJSON)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCP storage client: %v", err)
		return
	}
	defer client.Close()

	bucketName := req.BucketName
	objectKey := req.ObjectKey

	err = client.Bucket(bucketName).Object(objectKey).Delete(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting object from GCP: %v", err)
		return
	}

	// Send success response
	resp := BucketResponse{Message: "Object deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
