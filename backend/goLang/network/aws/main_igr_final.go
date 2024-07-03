package aws_vpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
) // InternetGatewayRequest represents the JSON request structure for creating internet gateways
type InternetGatewayRequest struct {
	VPCID     string `json:"vpcId"`
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

type ListInternetGatewaysResponse struct {
	InternetGateways []string `json:"internetGateways"`
}

// CreateInternetGatewayHandler handles POST requests to create an internet gateway within a VPC
func CreateInternetGatewayHandler(w http.ResponseWriter, r *http.Request) {
	var req InternetGatewayRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get cloud account details
	cloudAccount, err := GetCloudAccountDetails(req.AccountID)
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

	// Create EC2 service client
	svc := ec2.New(sess)

	// Create internet gateway
	resp, err := svc.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating internet gateway: %v", err)
		return
	}

	// Attach internet gateway to VPC
	_, err = svc.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(*resp.InternetGateway.InternetGatewayId),
		VpcId:             aws.String(req.VPCID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error attaching internet gateway to VPC: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Internet gateway %s created and attached to VPC %s successfully", *resp.InternetGateway.InternetGatewayId, req.VPCID)
	resp1 := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

// AttachInternetGatewayRequest represents the JSON request structure for attaching an internet gateway
type AttachInternetGatewayRequest struct {
	VPCID             string `json:"vpcId"`
	InternetGatewayID string `json:"internetGatewayId"`
	Region            string `json:"region"`
	AccountID         int    `json:"accountID"`
}

// AttachInternetGatewayHandler handles POST requests to attach an internet gateway to a VPC
func AttachInternetGatewayHandler(w http.ResponseWriter, r *http.Request) {
	var req AttachInternetGatewayRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get cloud account details
	cloudAccount, err := GetCloudAccountDetails(req.AccountID)
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

	// Create EC2 service client
	svc := ec2.New(sess)

	// Attach internet gateway to VPC
	_, err = svc.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(req.InternetGatewayID),
		VpcId:             aws.String(req.VPCID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error attaching internet gateway to VPC: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Internet gateway %s attached to VPC %s successfully", req.InternetGatewayID, req.VPCID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DetachInternetGatewayRequest represents the JSON request structure for detaching an internet gateway
type DetachInternetGatewayRequest struct {
	VPCID             string `json:"vpcId"`
	InternetGatewayID string `json:"internetGatewayId"`
	Region            string `json:"region"`
	AccountID         int    `json:"accountID"`
}

// DetachInternetGatewayHandler handles POST requests to detach an internet gateway from a VPC
func DetachInternetGatewayHandler(w http.ResponseWriter, r *http.Request) {
	var req DetachInternetGatewayRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get cloud account details
	cloudAccount, err := GetCloudAccountDetails(req.AccountID)
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

	// Create EC2 service client
	svc := ec2.New(sess)

	// Detach internet gateway from VPC
	_, err = svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(req.InternetGatewayID),
		VpcId:             aws.String(req.VPCID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error detaching internet gateway from VPC: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Internet gateway %s detached from VPC %s successfully", req.InternetGatewayID, req.VPCID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteInternetGatewayRequest represents the JSON request structure for deleting an internet gateway
type DeleteInternetGatewayRequest struct {
	InternetGatewayID string `json:"internetGatewayId"`
	Region            string `json:"region"`
	AccountID         int    `json:"accountID"`
}

// DeleteInternetGatewayHandler handles POST requests to delete an internet gateway
func DeleteInternetGatewayHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteInternetGatewayRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get cloud account details
	cloudAccount, err := GetCloudAccountDetails(req.AccountID)
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

	// Create EC2 service client
	svc := ec2.New(sess)

	// Delete internet gateway
	_, err = svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(req.InternetGatewayID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting internet gateway: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Internet gateway %s deleted successfully", req.InternetGatewayID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListInternetGatewaysRequest struct {
	AccountID int    `json:"accountID"`
	Region    string `json:"region"`
}

// ListAllInternetGatewaysHandler handles POST requests to list all internet gateways
func ListAllInternetGatewaysHandler(w http.ResponseWriter, r *http.Request) {
	// Get cloud account details
	var req ListInternetGatewaysRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	cloudAccount, err := GetCloudAccountDetails(req.AccountID)
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

	// Create EC2 service client
	svc := ec2.New(sess)

	// Describe internet gateways across all regions
	resp, err := svc.DescribeInternetGateways(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error describing internet gateways: %v", err)
		return
	}

	// Extract internet gateway details from the response
	var internetGateways []InternetGatewayDetails
	for _, ig := range resp.InternetGateways {
		var attachments []AttachmentDetails
		for _, attachment := range ig.Attachments {
			attachments = append(attachments, AttachmentDetails{
				State: *attachment.State,
				VpcID: *attachment.VpcId,
			})
		}

		internetGateways = append(internetGateways, InternetGatewayDetails{
			InternetGatewayID: *ig.InternetGatewayId,
			OwnerID:           *ig.OwnerId,
			Attachments:       attachments,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(internetGateways)
}

// Structs for request and response

type InternetGatewayDetails struct {
	InternetGatewayID string              `json:"internetGatewayId"`
	OwnerID           string              `json:"ownerId"`
	Attachments       []AttachmentDetails `json:"attachments"`
}

type AttachmentDetails struct {
	State string `json:"state"`
	VpcID string `json:"vpcId"`
}
