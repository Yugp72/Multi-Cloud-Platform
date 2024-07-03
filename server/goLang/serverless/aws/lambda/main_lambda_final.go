package aws_lambda

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaFunctionRequest struct {
	AccountID        int               `json:"accountID"`
	Region           string            `json:"region"`
	FunctionName     string            `json:"functionName"`
	Handler          string            `json:"handler"`
	Role             string            `json:"role"`
	Runtime          string            `json:"runtime"`
	S3Bucket         string            `json:"s3Bucket"`
	S3Key            string            `json:"s3Key"`
	Environment      map[string]string `json:"environment"`
	MemorySize       int64             `json:"memorySize"`
	Timeout          int64             `json:"timeout"`
	SecurityGroupIds []string          `json:"securityGroupIds"`
	SubnetIds        []string          `json:"subnetIds"`
}

type LambdaFunctionResponse struct {
	Message string `json:"message"`
}

func initLambdaService(accessKey, secretKey, region string) (*lambda.Lambda, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %v", err)
	}

	svc := lambda.New(sess)

	return svc, nil
}

func CreateLambdaFunctionHandler(w http.ResponseWriter, r *http.Request) {
	var req LambdaFunctionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	// Initialize AWS Lambda service
	lambdaSvc, err := initLambdaService(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, req.Region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error initializing AWS Lambda service: %v", err), http.StatusInternalServerError)
		return
	}

	params := &lambda.CreateFunctionInput{
		FunctionName: aws.String(req.FunctionName),
		Handler:      aws.String(req.Handler),
		Role:         aws.String(req.Role),
		Runtime:      aws.String(req.Runtime),
		Code: &lambda.FunctionCode{
			S3Bucket: aws.String(req.S3Bucket),
			S3Key:    aws.String(req.S3Key),
		},
		Environment: &lambda.Environment{
			Variables: aws.StringMap(req.Environment),
		},
		MemorySize: aws.Int64(req.MemorySize),
		Timeout:    aws.Int64(req.Timeout),
		VpcConfig: &lambda.VpcConfig{
			SecurityGroupIds: aws.StringSlice(req.SecurityGroupIds),
			SubnetIds:        aws.StringSlice(req.SubnetIds),
		},
	}

	_, err = lambdaSvc.CreateFunction(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Lambda function: %v", err), http.StatusInternalServerError)
		return
	}

	resp := LambdaFunctionResponse{Message: "Lambda function created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListLambdaFunctionsRequest struct {
	AccountID int `json:"accountID"`
}

type DeleteLambdaFunctionRequest struct {
	AccountID    int    `json:"accountID"`
	Region       string `json:"region"`
	FunctionName string `json:"functionName"`
}

// for all region list
func ListLambdaFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListLambdaFunctionsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)

	// Initialize AWS Lambda service
	lambdaSvc, err := initLambdaService(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, cloudAccount.Region.String)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error initializing AWS Lambda service: %v", err), http.StatusInternalServerError)
		return
	}

	// List Lambda functions
	result, err := lambdaSvc.ListFunctions(&lambda.ListFunctionsInput{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing Lambda functions: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract function names from the response
	var functions []string
	for _, function := range result.Functions {
		functions = append(functions, *function.FunctionName)
	}

	// Send success response
	resp := struct {
		Functions []string `json:"functions"`
	}{
		Functions: functions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteLambdaFunctionHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteLambdaFunctionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	// Initialize AWS Lambda service
	lambdaSvc, err := initLambdaService(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, req.Region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error initializing AWS Lambda service: %v", err), http.StatusInternalServerError)
		return
	}

	// Delete Lambda function
	_, err = lambdaSvc.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String(req.FunctionName),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting Lambda function: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	resp := LambdaFunctionResponse{Message: "Lambda function deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
