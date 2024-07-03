package gcp_gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	db "btep.project/databaseConnection"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

type BucketRequest struct {
	BucketName string `json:"bucketName"`
	AccountID  int    `json:"accountID"`
	Token      string `json:"token"`
}
type BucketResponse struct {
	Message string `json:"message"`
}

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var req BucketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}
	token := req.Token
	ctx := context.Background()
	creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS client: %v", err)
		return
	}
	defer client.Close()
	if err := client.Bucket(req.BucketName).Create(ctx, "divine-treat-413716", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS bucket: %v", err)
		return
	}
	resp := BucketResponse{Message: "GCS bucket created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ObjectRequest struct {
	BucketName string                `json:"bucketName"`
	ObjectName string                `json:"objectName"`
	File       *multipart.FileHeader `json:"file"`
	AccountID  int                   `json:"accountID"`
	Token      string                `json:"token"`
}

// ObjectResponse represents the JSON response structure for object operations
type ObjectResponse struct {
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
}

// UploadObjectHandler handles POST requests to upload an object to a GCS bucket
func UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Max size 10MB, adjust as needed
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	// Parse the form fields
	var req ObjectRequest
	err = json.NewDecoder(r.Body).Decode(&req)
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
	fmt.Println(cloudAccount)

	// Extract the token from the request
	token := req.Token

	// Create a GCS client with token authentication
	ctx := context.Background()
	creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS client: %v", err)
		return
	}
	defer client.Close()

	// Upload the object to the GCS bucket
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error retrieving file from form: %v", err)
		return
	}
	defer file.Close()

	wc := client.Bucket(req.BucketName).Object(req.ObjectName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error uploading object to GCS: %v", err)
		return
	}
	if err := wc.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error closing object writer: %v", err)
		return
	}

	// Send success response with the object URL
	objectURL := "https://storage.googleapis.com/" + req.BucketName + "/" + req.ObjectName
	resp := ObjectResponse{Message: "Object uploaded successfully", URL: objectURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type GetObjectRequest struct {
	BucketName string `json:"bucketName"`
	ObjectName string `json:"objectName"`
	AccountID  int    `json:"accountID"`
	Token      string `json:"token"`
}

// GetObjectHandler handles GET requests to retrieve an object from a GCS bucket
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req GetObjectRequest
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
	fmt.Println(cloudAccount)

	// Extract the token from the request
	token := req.Token

	// Create a GCS client with token authentication
	ctx := context.Background()
	creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS client: %v", err)
		return
	}
	defer client.Close()

	// Download the object from the GCS bucket
	rc, err := client.Bucket(req.BucketName).Object(req.ObjectName).NewReader(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error downloading object from GCS: %v", err)
		return
	}
	defer rc.Close()

	// Read the object content
	objectContent, err := ioutil.ReadAll(rc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading object content: %v", err)
		return
	}

	// Send success response with the object content
	w.Header().Set("Content-Disposition", "attachment; filename="+req.ObjectName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(objectContent)
}

type DeleteObjectRequest struct {
	BucketName string `json:"bucketName"`
	ObjectName string `json:"objectName"`
	AccountID  int    `json:"accountID"`
	Token      string `json:"token"`
}

// DeleteObjectHandler handles DELETE requests to delete an object from a GCS bucket
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req DeleteObjectRequest
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
	fmt.Println(cloudAccount)

	// Extract the token from the request
	token := req.Token

	// Create a GCS client with token authentication
	ctx := context.Background()
	creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS client: %v", err)
		return
	}
	defer client.Close()

	// Delete the object from the GCS bucket
	err = client.Bucket(req.BucketName).Object(req.ObjectName).Delete(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting object from GCS: %v", err)
		return
	}

	// Send success response
	resp := ObjectResponse{Message: "Object deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req BucketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	token := req.Token

	// Create a GCS client with token authentication
	ctx := context.Background()
	creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS client: %v", err)
		return
	}
	defer client.Close()

	// Delete the GCS bucket
	client.Bucket(req.BucketName).Delete(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting GCS bucket: %v", err)
		return
	}

	// Send success response
	resp := BucketResponse{Message: "GCS bucket deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListBucketRequest struct {
	Token     string `json:"token"`
	AccountID int    `json:"accountID"`
}

func ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req ListBucketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Extract the token from the request
	token := req.Token
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)

	// Create a GCS client with token authentication
	ctx := context.Background()
	creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS client: %v", err)
		return
	}
	defer client.Close()

	// List all GCS buckets
	buckets := client.Buckets(ctx, cloudAccount.ProjectID.String)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing GCS buckets: %v", err)
		return
	}

	// Send success response with the list of buckets
	type BucketDetails struct {
		Name    string `json:"name"`
		Created string `json:"created"`
	}

	var bucketDetails []BucketDetails

	for {
		bucketAttrs, err := buckets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error listing GCS buckets: %v", err)
			return
		}

		bucketDetail := BucketDetails{
			Name:    bucketAttrs.Name,
			Created: bucketAttrs.Created.Format(time.RFC3339),
		}
		bucketDetails = append(bucketDetails, bucketDetail)
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	// Encode bucketDetails slice as JSON and write it to the response
	if err := json.NewEncoder(w).Encode(bucketDetails); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding response: %v", err)
		return
	}
}
