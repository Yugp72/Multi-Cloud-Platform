package aws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"btep.project/backend/golang/db"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketRequest represents the JSON request structure
type BucketRequest struct {
	BucketName string `json:"bucketName"`
	Region     string `json:"region"`
	AccountID  int    `json:"accountID"`
}

// ObjectRequest represents the JSON request structure for object operations
type ObjectRequest struct {
	BucketName string `json:"bucketName"`
	ObjectKey  string `json:"objectKey"`
	Content    string `json:"content"`
	Region     string `json:"region"`
	AccountID  int    `json:"accountID"`
}

type BucketResponse struct {
	Message string `json:"message"`
}

// CreateBucketHandler handles POST requests to create an S3 bucket
func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var req BucketRequest
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
		Credentials: aws.NewStaticCredentials(*cloudAccount.AccessKey, *cloudAccount.SecretKey, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Create the S3 bucket
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(req.BucketName),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating bucket: %v", err)
		return
	}

	// Send success response
	resp := BucketResponse{Message: "Bucket created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UploadObjectHandler handles POST requests to upload an object to an S3 bucket
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

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*cloudAccount.Region),
		Credentials: aws.NewStaticCredentials(*cloudAccount.AccessKey, *cloudAccount.SecretKey, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Parse the form data
	err = r.ParseMultipartForm(10 << 20) // Set maxMemory to 10 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data")
		return
	}

	bucketName := req.BucketName
	file, handler, err := r.FormFile("Content")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error retrieving file from form data")
		return
	}
	defer file.Close()

	objectKey := handler.Filename // Use filename as object key

	// Upload the object to the S3 bucket
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error uploading object to S3: %v", err)
		return
	}

	// Send success response
	resp := BucketResponse{Message: "Object uploaded successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetObjectHandler handles GET requests to retrieve an object from an S3 bucket
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*cloudAccount.Region),
		Credentials: aws.NewStaticCredentials(*cloudAccount.AccessKey, *cloudAccount.SecretKey, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Download the object from the S3 bucket
	objOutput, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(req.BucketName),
		Key:    aws.String(req.ObjectKey),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error downloading object from S3: %v", err)
		return
	}

	// Read the object content
	objContent, err := ioutil.ReadAll(objOutput.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading object content: %v", err)
		return
	}

	// Send the object content in the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+req.ObjectKey)
	w.Write(objContent)
}

// DeleteObjectHandler handles POST requests to delete an object from an S3 bucket
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*cloudAccount.Region),
		Credentials: aws.NewStaticCredentials(*cloudAccount.AccessKey, *cloudAccount.SecretKey, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Delete the object from the S3 bucket
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(req.BucketName),
		Key:    aws.String(req.ObjectKey),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting object from S3: %v", err)
		return
	}

	// Send success response
	resp := BucketResponse{Message: "Object deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
