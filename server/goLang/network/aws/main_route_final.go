package aws_vpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type RouteTableRequest struct {
	VPCID     string   `json:"vpcId"`
	SubnetID  string   `json:"subnetId"`
	Route     []string `json:"route"`
	Region    string   `json:"region"`
	AccountID int      `json:"accountID"`
}

type DeleteRouteTableRequest struct {
	RouteTableID string `json:"routeTableId"`
	Region       string `json:"region"`
	AccountID    int    `json:"accountID"`
}

// CreateRouteTableHandler handles POST requests to create a route table within a VPC
func CreateRouteTableHandler(w http.ResponseWriter, r *http.Request) {
	var req RouteTableRequest
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

	// Create route table within the VPC
	resp, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(req.VPCID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating route table: %v", err)
		return
	}

	// Associate route table with subnet
	_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(*resp.RouteTable.RouteTableId),
		SubnetId:     aws.String(req.SubnetID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error associating route table with subnet: %v", err)
		return
	}

	// Add routes to the route table
	for _, route := range req.Route {
		_, err = svc.CreateRoute(&ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String(route),
			RouteTableId:         aws.String(*resp.RouteTable.RouteTableId),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error adding route to route table: %v", err)
			return
		}
	}

	// Send success response
	respMsg := fmt.Sprintf("Route table %s created and associated with subnet %s successfully", *resp.RouteTable.RouteTableId, req.SubnetID)
	resp1 := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

// DeleteRouteTableRequest represents the JSON request structure for deleting a route table

// DeleteRouteTableHandler handles POST requests to delete a route table
func DeleteRouteTableHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteRouteTableRequest
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

	// Delete the route table
	_, err = svc.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(req.RouteTableID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting route table: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Route table %s deleted successfully", req.RouteTableID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ListRouteTableHandler(w http.ResponseWriter, r *http.Request) {
	// Get cloud account details
	cloudAccount, err := GetCloudAccountDetails(1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create EC2 service client
	svc := ec2.New(sess)

	// List route tables
	resp, err := svc.DescribeRouteTables(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing route tables: %v", err)
		return
	}

	// Send success response
	resp1 := RouteTableResponse{
		Message:     "Route tables listed successfully",
		RouteTables: resp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

type RouteTableResponse struct {
	Message     string                         `json:"message"`
	RouteTables *ec2.DescribeRouteTablesOutput `json:"routeTables"`
}
