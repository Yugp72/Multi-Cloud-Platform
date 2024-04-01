package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"cloud.google.com/go/storage"
)

// BucketRequest represents the JSON request structure
type BucketRequest struct {
	BucketName string `json:"bucketName"`
}

// ObjectRequest represents the JSON request structure for object operations
type ObjectRequest struct {
	BucketName string `json:"bucketName"`
	ObjectName string `json:"objectName"`
}

// BucketResponse represents the JSON response structure
type BucketResponse struct {
	Message string `json:"message"`
}

// ObjectResponse represents the JSON response structure for object operations
type ObjectResponse struct {
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
}

// Initialize a Google Cloud Storage client
var client *storage.Client

func main() {
	// Create a Google Cloud Storage client
	ctx := context.Background()
	var err error
	client, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating GCS client: %v", err)
	}
	defer client.Close()

	// Define routes
	http.HandleFunc("/createBucket", CreateBucketHandler)
	http.HandleFunc("/uploadObject", UploadObjectHandler)
	http.HandleFunc("/getObject", GetObjectHandler)
	http.HandleFunc("/deleteObject", DeleteObjectHandler)

	// Start the HTTP server
	port := "8000"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// CreateBucketHandler handles POST requests to create a GCS bucket
func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req BucketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Create a new GCS bucket
	err = createGCSBucket(req.BucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCS bucket: %v", err)
		return
	}

	// Send success response
	resp := BucketResponse{Message: "GCS bucket created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UploadObjectHandler handles POST requests to upload an object to a GCS bucket
func UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req ObjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error retrieving file from request")
		return
	}
	fmt.Print(handler.Filename)
	defer file.Close()

	// Upload the object to the GCS bucket
	objectURL, err := uploadGCSObject(req.BucketName, req.ObjectName, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error uploading object to GCS: %v", err)
		return
	}

	// Send success response with the object URL
	resp := ObjectResponse{Message: "Object uploaded successfully", URL: objectURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetObjectHandler handles GET requests to retrieve an object from a GCS bucket
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req ObjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Download the object from the GCS bucket
	objectContent, err := downloadGCSObject(req.BucketName, req.ObjectName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error downloading object from GCS: %v", err)
		return
	}

	// Send the object content as the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+req.ObjectName)
	w.Write(objectContent)
}

// DeleteObjectHandler handles DELETE requests to delete an object from a GCS bucket
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req ObjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Delete the object from the GCS bucket
	err = deleteGCSObject(req.BucketName, req.ObjectName)
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

func uploadGCSObject(bucketName, objectName string, file multipart.File) (string, error) {
	ctx := context.Background()
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	// Object successfully uploaded, return its URL
	return "https://storage.googleapis.com/" + bucketName + "/" + objectName, nil
}

func createGCSBucket(bucketName string) error {
	ctx := context.Background()
	return client.Bucket(bucketName).Create(ctx, "divine-treat-413716", nil)
}

func downloadGCSObject(bucketName, objectName string) ([]byte, error) {
	ctx := context.Background()
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

func deleteGCSObject(bucketName, objectName string) error {
	ctx := context.Background()
	return client.Bucket(bucketName).Object(objectName).Delete(ctx)
}
