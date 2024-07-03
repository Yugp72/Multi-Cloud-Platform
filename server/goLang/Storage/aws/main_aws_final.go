package aws_s3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	db "btep.project/databaseConnection"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
	BucketName string                `json:"bucketName"`
	ObjectKey  string                `json:"objectKey"`
	Content    *multipart.FileHeader `json:"content"`
	Region     string                `json:"region"`
	AccountID  int                   `json:"accountID"`
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
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	fmt.Println(cloudAccount)
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
	svc := s3.New(sess)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(req.BucketName),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating bucket: %v", err)
		return
	}
	resp := BucketResponse{Message: "Bucket created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // Set maxMemory to 10 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	// Extract fields from form data
	bucketName := r.FormValue("bucketName")
	region := r.FormValue("region")
	accountIDStr := r.FormValue("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid account ID: %v", err)
		return
	}

	// Get CloudAccount details
	cloudAccount, err := db.GetCloudAccountDetails(accountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(*&cloudAccount.AccessKey.String, *&cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session: %v", err)
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Retrieve the uploaded file
	file, handler, err := r.FormFile("content")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error retrieving file from form data: %v", err)
		return
	}
	defer file.Close()

	// Use the file name as the object key
	objectKey := handler.Filename

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
	fmt.Println(cloudAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
		// Credentials: credentials.NewStaticCredentials(*cloudAccount.accessKey, *cloudAccount.secretKey, ""),

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
	fmt.Println(cloudAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
		// Credentials: credentials.NewStaticCredentials(*cloudAccount.accessKey, *cloudAccount.secretKey, ""),

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

type DeleteBucketRequest struct {
	BucketName string `json:"bucketName"`
	Region     string `json:"region"`
	AccountID  int    `json:"accountID"`
}

// DeleteBucketHandler handles POST requests to delete an S3 bucket
func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var req BucketRequest
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
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create S3 service client
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(req.BucketName),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing objects in the bucket: %v", err)
		return
	}

	for _, obj := range resp.Contents {
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(req.BucketName),
			Key:    obj.Key,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error deleting object %s from bucket: %v", *obj.Key, err)
			return
		}
	}

	// Delete the bucket from S3
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(req.BucketName),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting bucket from S3: %v", err)
		return
	}

	// Send success response
	resp1 := BucketResponse{Message: "Bucket deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

type ListBucketRequest struct {
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}
type BucketDetails struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type ListBucketResponse struct {
	Buckets []*BucketDetails `json:"buckets"`
}

func ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListBucketRequest
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
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	// all buckets with its detials
	bucketsOutput, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing buckets: %v", err)
		//here return null responce

	}

	var bucketDetails []*BucketDetails
	for _, bucket := range bucketsOutput.Buckets {
		bucketDetails = append(bucketDetails, &BucketDetails{Name: *bucket.Name, Date: bucket.CreationDate.String()})
	}

	// Send success response
	resp := ListBucketResponse{Buckets: bucketDetails}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
