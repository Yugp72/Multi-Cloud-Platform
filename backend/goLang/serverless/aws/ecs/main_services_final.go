package aws_ecs

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ServiceCreateRequest struct {
	AccountID               int                 `json:"projectID"`
	Region                  string              `json:"region"`
	ClusterID               string              `json:"clusterID"`
	ServiceName             string              `json:"serviceName"`
	TaskDefinition          string              `json:"taskDefinition"`
	DesiredCount            int64               `json:"desiredCount"`
	Subnets                 []*string           `json:"subnets"`
	SecurityGroups          []*string           `json:"securityGroups"`
	AssignPublicIP          string              `json:"assignPublicIP"`
	LoadBalancers           []*ecs.LoadBalancer `json:"loadBalancers"`
	DeploymentConfiguration struct {
		MaximumPercent        int64 `json:"maximumPercent"`
		MinimumHealthyPercent int64 `json:"minimumHealthyPercent"`
	} `json:"deploymentConfiguration"`
}

func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	var req ServiceCreateRequest
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

	// Create ECS service
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}
	svc := ecs.New(sess)
	_, err = svc.CreateService(&ecs.CreateServiceInput{
		Cluster:        aws.String(req.ClusterID),
		ServiceName:    aws.String(req.ServiceName),
		TaskDefinition: aws.String(req.TaskDefinition),
		DesiredCount:   aws.Int64(req.DesiredCount),
		// Add other configuration parameters as needed
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        req.Subnets,
				SecurityGroups: req.SecurityGroups,
				AssignPublicIp: aws.String(req.AssignPublicIP),
			},
		},
		DeploymentConfiguration: &ecs.DeploymentConfiguration{
			MaximumPercent:        aws.Int64(req.DeploymentConfiguration.MaximumPercent),
			MinimumHealthyPercent: aws.Int64(req.DeploymentConfiguration.MinimumHealthyPercent),
		},
		LoadBalancers: req.LoadBalancers,
		// Add other deployment-related parameters here
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating service: %v", err)
		return
	}

	// Send success response
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Service created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ClusterListRequest struct {
	AccountID int    `json:"projectID"`
	Region    string `json:"region"`
	ClusterID string `json:"clusterID"`
}

func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	var req ClusterListRequest
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

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}
	svc := ecs.New(sess)

	// List ECS services

	resp, err := svc.ListServices(&ecs.ListServicesInput{
		Cluster: aws.String(req.ClusterID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing services: %v", err)
		return
	}

	// Send response with service names
	var serviceNames []string
	for _, service := range resp.ServiceArns {
		serviceNames = append(serviceNames, *service)
	}
	resp1 := struct {
		Services []string `json:"services"`
	}{
		Services: serviceNames,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

type ServiceUpdateRequest struct {
	ProjectID               int                 `json:"projectID"`
	Region                  string              `json:"region"`
	ClusterID               string              `json:"clusterID"`
	ServiceID               string              `json:"serviceID"`
	TaskDefinition          string              `json:"taskDefinition"`
	DesiredCount            int64               `json:"desiredCount"`
	Subnets                 []*string           `json:"subnets"`
	SecurityGroups          []*string           `json:"securityGroups"`
	AssignPublicIP          string              `json:"assignPublicIP"`
	LoadBalancers           []*ecs.LoadBalancer `json:"loadBalancers"`
	DeploymentConfiguration struct {
		MaximumPercent        int64 `json:"maximumPercent"`
		MinimumHealthyPercent int64 `json:"minimumHealthyPercent"`
	} `json:"deploymentConfiguration"`
}

func UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {
	var req ServiceUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Update ECS service
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}
	svc := ecs.New(sess)
	_, err = svc.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        aws.String(req.ClusterID),
		Service:        aws.String(req.ServiceID),
		TaskDefinition: aws.String(req.TaskDefinition),
		DesiredCount:   aws.Int64(req.DesiredCount),
		// Add other configuration parameters as needed
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        req.Subnets,
				SecurityGroups: req.SecurityGroups,
				AssignPublicIp: aws.String(req.AssignPublicIP),
			},
		},
		DeploymentConfiguration: &ecs.DeploymentConfiguration{
			MaximumPercent:        aws.Int64(req.DeploymentConfiguration.MaximumPercent),
			MinimumHealthyPercent: aws.Int64(req.DeploymentConfiguration.MinimumHealthyPercent),
		},
		LoadBalancers: req.LoadBalancers,
		// Add other deployment-related parameters here
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating service: %v", err)
		return
	}

	// Send success response
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Service updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ServiceDeleteRequest struct {
	AccountID int    `json:"projectID"`
	Region    string `json:"region"`
	ClusterID string `json:"clusterID"`
	ServiceID string `json:"serviceID"`
}

func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	var req ServiceDeleteRequest
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

	// Delete ECS service
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cloudAccount.AccessKey.String, cloudAccount.SecretKey.String, ""),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}
	svc := ecs.New(sess)
	_, err = svc.DeleteService(&ecs.DeleteServiceInput{
		Cluster: aws.String(req.ClusterID),
		Service: aws.String(req.ServiceID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting service: %v", err)
		return
	}

	// Send success response
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Service deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
