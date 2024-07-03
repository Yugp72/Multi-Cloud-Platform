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

type SubnetRequest struct {
	VPCID            string `json:"vpcId"`
	SubnetName       string `json:"subnetName"`
	IPv4CIDRBlock    string `json:"ipv4CIDRBlock"`
	IPv6CIDRBlock    string `json:"ipv6CIDRBlock,omitempty"`
	AvailabilityZone string `json:"availabilityZone"`
	Region           string `json:"region"`
	AccountID        int    `json:"accountID"`
}

// CreateSubnetHandler handles POST requests to create a subnet within a VPC
func CreateSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req SubnetRequest
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

	// Create subnet within the VPC
	resp, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
		CidrBlock:        aws.String(req.IPv4CIDRBlock),
		Ipv6CidrBlock:    aws.String(req.IPv6CIDRBlock),
		VpcId:            aws.String(req.VPCID),
		AvailabilityZone: aws.String(req.AvailabilityZone),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating subnet: %v", err)
		return
	}

	// Tag subnet with name
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{resp.Subnet.SubnetId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(req.SubnetName),
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error tagging subnet: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Subnet %s created successfully", *resp.Subnet.SubnetId)
	resp1 := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

// UpdateSubnetRequest represents the JSON request structure for updating a subnet
type UpdateSubnetRequest struct {
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName"`
	IPv4CIDRBlock    string `json:"ipv4CIDRBlock"`
	IPv6CIDRBlock    string `json:"ipv6CIDRBlock,omitempty"`
	AvailabilityZone string `json:"availabilityZone"`
	Region           string `json:"region"`
	AccountID        int    `json:"accountID"`
}

// UpdateSubnetHandler handles POST requests to update a subnet within a VPC
func UpdateSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateSubnetRequest
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

	_, err = svc.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
		SubnetId: aws.String(req.SubnetID),
		MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating subnet: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Subnet %s updated successfully", req.SubnetID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteSubnetRequest represents the JSON request structure for deleting a subnet
type DeleteSubnetRequest struct {
	SubnetID  string `json:"subnetId"`
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

// DeleteSubnetHandler handles POST requests to delete a subnet within a VPC
func DeleteSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteSubnetRequest
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

	// Delete subnet
	_, err = svc.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(req.SubnetID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting subnet: %v", err)
		return
	}

	// Send success response
	respMsg := fmt.Sprintf("Subnet %s deleted successfully", req.SubnetID)
	resp := VPCResponse{Message: respMsg}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListSubnetsRequest represents the JSON request structure for listing subnets
type ListSubnetsRequest struct {
	VPCID     string `json:"vpcId"`
	Region    string `json:"region"`
	AccountID int    `json:"accountID"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ListSubnetsHandler handles POST requests to list subnets within a VPC
func ListSubnetsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListSubnetsRequest
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

	// Describe subnets
	resp, err := svc.DescribeSubnets(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing subnets: %v", err)
		return
	}

	// Send success response
	var subnets []Subnet

	for _, subnet := range resp.Subnets {
		var ipv6CIDRBlock string
		if len(subnet.Ipv6CidrBlockAssociationSet) > 0 {
			ipv6CIDRBlock = *subnet.Ipv6CidrBlockAssociationSet[0].Ipv6CidrBlock
		} else {
			ipv6CIDRBlock = ""
		}

		var subnetTags []Tag
		if subnet.Tags != nil {
			for _, tag := range subnet.Tags {
				subnetTags = append(subnetTags, Tag{
					Key:   *tag.Key,
					Value: *tag.Value,
				})
			}
		}

		var ec2SubnetTags []*ec2.Tag
		for _, tag := range subnetTags {
			ec2SubnetTags = append(ec2SubnetTags, &ec2.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}

		subnets = append(subnets, Subnet{
			ID:                      *subnet.SubnetId,
			IPv4CIDRBlock:           *subnet.CidrBlock,
			IPv6CIDRBlock:           ipv6CIDRBlock,
			AvailabilityZone:        *subnet.AvailabilityZone,
			AvailabilityZoneID:      *subnet.AvailabilityZoneId,
			AvailableIPAddressCount: *subnet.AvailableIpAddressCount,
			State:                   *subnet.State,
			SubnetArn:               *subnet.SubnetArn,
			VpcID:                   *subnet.VpcId,
			Name:                    findTagValue(subnet.Tags, "Name"),
		})
	}

	resp1 := SubnetsResponse{Subnets: subnets}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

// Subnet represents a subnet within a VPC
type Subnet struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	IPv4CIDRBlock    string `json:"ipv4CIDRBlock"`
	IPv6CIDRBlock    string `json:"ipv6CIDRBlock,omitempty"`
	AvailabilityZone string `json:"availabilityZone"`
	// Add more fields here as needed
	AvailabilityZoneID      string `json:"availabilityZoneID"`
	AvailableIPAddressCount int64  `json:"availableIPAddressCount"`
	State                   string `json:"state"`
	SubnetArn               string `json:"subnetArn"`
	VpcID                   string `json:"vpcID"`
}

// SubnetsResponse represents the JSON response structure for listing subnets
type SubnetsResponse struct {
	Subnets []Subnet `json:"subnets"`
}

func findTagValue(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}
