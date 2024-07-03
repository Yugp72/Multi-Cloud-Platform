package aws_ec2

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

// InstanceRequest represents the JSON request structure for EC2 instance operations
type InstanceRequest struct {
	InstanceType     string   `json:"instanceType"`
	AmiID            string   `json:"amiID"`
	KeyName          string   `json:"keyName"`
	SecurityGroupIDs []string `json:"securityGroupIDs"`
	SubnetID         string   `json:"subnetID"`
	Region           string   `json:"region"`
	Name             string   `json:"name"`
	AccountID        int      `json:"accountID"`
}

// InstanceResponse represents the JSON response structure for EC2 instance operations
type InstanceResponse struct {
	Message     string   `json:"message"`
	InstanceIDs []string `json:"instanceIDs,omitempty"`
}

type InstanceListRequest struct {
	AccountID int    `json:"accountID"`
	Region    string `json:"region"`
}

type TerminalRequest struct {
	AccountID  int    `json:"accountID"`
	InstanceID string `json:"instanceID"`
}

// Initialize AWS session

// CreateInstanceHandler handles POST requests to create an EC2 instance
// CreateInstanceHandler handles POST requests to create an EC2 instance
func CreateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var req InstanceRequest
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

	// Initialize AWS session with the provided region
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session: %v", err)
		return
	}
	svc := ec2.New(sess)

	// Create an EC2 instance
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(req.AmiID),
		InstanceType:     aws.String(req.InstanceType),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          aws.String(req.KeyName),
		SecurityGroupIds: aws.StringSlice(req.SecurityGroupIDs),
		SubnetId:         aws.String(req.SubnetID),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(req.Name), // Set your desired instance name
					},
				},
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating EC2 instance: %v", err)
		return
	}

	instanceID := *runResult.Instances[0].InstanceId
	resp := InstanceResponse{Message: fmt.Sprintf("EC2 instance created with ID: %s", instanceID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type InstanceDetails struct {
	InstanceID       string                 `json:"instance_id"`
	InstanceType     string                 `json:"instance_type"`
	LaunchTime       string                 `json:"launch_time"`
	PrivateIPAddress string                 `json:"private_ip_address"`
	PublicIPAddress  string                 `json:"public_ip_address"`
	AvailabilityZone string                 `json:"availability_zone"`
	StateCode        string                 `json:"state_code"`
	SecurityGroups   []*ec2.GroupIdentifier `json:"security_groups"`
	PlatformDetails  string                 `json:"platform_details"`
	Tags             []*ec2.Tag             `json:"tags"`
}

type ListInstanceResponse struct {
	Instances []*InstanceDetails `json:"instances"`
}

// ListInstancesHandler handles GET requests to list EC2 instances
func ListInstancesHandler(w http.ResponseWriter, r *http.Request) {
	var req InstanceListRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	fmt.Println(cloudAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Initialize AWS session with the provided region
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cloudAccount.Region.String),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})

	fmt.Println("sess: ", sess)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session: %v", err)
		return
	}
	svc := ec2.New(sess)

	fmt.Println("svc: ", svc)

	// List EC2 instances
	result, err := svc.DescribeInstances(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing EC2 instances: %v", err)
		return
	}

	fmt.Println("result: ", result)

	var instances []InstanceDetails
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceDetails := InstanceDetails{
				InstanceID:       "",
				InstanceType:     "",
				LaunchTime:       "",
				PrivateIPAddress: "",
				PublicIPAddress:  "",
				AvailabilityZone: "",
				StateCode:        "",
				SecurityGroups:   nil,
				PlatformDetails:  "",
				Tags:             []*ec2.Tag{},
			}
			if instance.InstanceId != nil {
				instanceDetails.InstanceID = *instance.InstanceId
			}
			if instance.InstanceType != nil {
				instanceDetails.InstanceType = *instance.InstanceType
			}
			if instance.Placement != nil && instance.Placement.AvailabilityZone != nil {
				instanceDetails.AvailabilityZone = *instance.Placement.AvailabilityZone
			}
			if instance.LaunchTime != nil {
				instanceDetails.LaunchTime = instance.LaunchTime.String()
			}
			if instance.PrivateIpAddress != nil {
				instanceDetails.PrivateIPAddress = *instance.PrivateIpAddress
			}
			if instance.PublicIpAddress != nil {
				instanceDetails.PublicIPAddress = *instance.PublicIpAddress
			}
			if instance.State != nil && instance.State.Code != nil {
				instanceDetails.StateCode = *instance.State.Name
			}
			if instance.SecurityGroups != nil {
				instanceDetails.SecurityGroups = instance.SecurityGroups
			}
			if instance.Platform != nil {
				instanceDetails.PlatformDetails = *instance.Platform
			}
			if instance.Tags != nil {
				instanceDetails.Tags = instance.Tags
			}

			instances = append(instances, instanceDetails)
		}
	}

	fmt.Println("instances: ", instances)

	var instancePointers []*InstanceDetails
	for i := range instances {
		instancePointers = append(instancePointers, &instances[i])
	}
	resp := ListInstanceResponse{Instances: instancePointers}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// TerminateInstanceHandler handles POST requests to terminate an EC2 instance
func TerminateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var req TerminalRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	if req.InstanceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Instance ID is required")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	fmt.Println(cloudAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Initialize AWS session with the provided region
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cloudAccount.Region.String),
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session: %v", err)
		return
	}
	svc := ec2.New(sess)

	_, err = svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(req.InstanceID)},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error terminating EC2 instance: %v", err)
		return
	}

	resp := InstanceResponse{Message: fmt.Sprintf("EC2 instance terminated: %s", req.InstanceID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
