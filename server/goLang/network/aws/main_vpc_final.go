package aws_vpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type VPCRequest struct {
	VPCName       string `json:"vpcName"`
	IPv4CIDRBlock string `json:"ipv4CIDRBlock"`
	IPv6CIDRBlock string `json:"ipv6CIDRBlock,omitempty"`
	Tenancy       string `json:"tenancy,omitempty"`
	Region        string `json:"region"`
	AccountID     int    `json:"accountID"`
}

type VPCResponse struct {
	Message string `json:"message"`
}

// GetCloudAccountDetails retrieves cloud account details from the database
func GetCloudAccountDetails(accountID int) (*db.CloudAccount, error) {
	return db.GetCloudAccountDetails(accountID)
}

// CreateVPCHandler handles POST requests to create a VPC
func CreateVPCHandler(w http.ResponseWriter, r *http.Request) {
	var req VPCRequest
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

	// Create VPC
	resp, err := svc.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock:       aws.String(req.IPv4CIDRBlock),
		InstanceTenancy: aws.String(req.Tenancy),
		Ipv6CidrBlock:   aws.String(req.IPv6CIDRBlock),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating VPC: %v", err)
		return
	}

	// Tag VPC with name
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{resp.Vpc.VpcId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(req.VPCName),
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error tagging VPC: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("VPC %s created successfully", *resp.Vpc.VpcId)
	resp1 := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

// DeleteVPCRequest represents the JSON request structure for deleting a VPC
type DeleteVPCRequest struct {
	VPCID     string `json:"vpcId"`
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

// DeleteVPCHandler handles POST requests to delete a VPC
func DeleteVPCHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteVPCRequest
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

	// Delete VPC
	_, err = svc.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: aws.String(req.VPCID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting VPC: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("VPC %s deleted successfully", req.VPCID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListVPCsRequest struct {
	AccountID int    `json:"accountID"`
	Region    string `json:"region"`
}

// ListVPCsHandler handles POST requests to list all VPCs
type ListVPCsResponse struct {
	Vpcs []Vpc `json:"Vpcs"`
}

type Vpc struct {
	CidrBlock       string `json:"CidrBlock"`
	DhcpOptionsId   string `json:"DhcpOptionsId"`
	InstanceTenancy string `json:"InstanceTenancy"`
	IsDefault       bool   `json:"IsDefault"`
	OwnerId         string `json:"OwnerId"`
	State           string `json:"State"`
	VpcId           string `json:"VpcId"`
	// Add other necessary fields here
}

func ListVPCsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListVPCsRequest
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

	// Describe VPCs
	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error describing VPCs: %v", err)
		return
	}

	// Extract necessary VPC information
	var vpcs []Vpc
	for _, v := range resp.Vpcs {
		vpc := Vpc{
			CidrBlock:       *v.CidrBlock,
			DhcpOptionsId:   *v.DhcpOptionsId,
			InstanceTenancy: *v.InstanceTenancy,
			IsDefault:       *v.IsDefault,
			OwnerId:         *v.OwnerId,
			State:           *v.State,
			VpcId:           *v.VpcId,
			// Add other necessary fields here
		}
		vpcs = append(vpcs, vpc)
	}

	// Send success response
	response := ListVPCsResponse{Vpcs: vpcs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
