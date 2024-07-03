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

type ClusterRequest struct {
	AccountID      int    `json:"accountId"`
	Region         string `json:"region"`
	ClusterID      string `json:"clusterId"`
	Infrastructure string `json:"infrastructure"`
}

type ClusterDeleteRequest struct {
	AccountID int    `json:"accountId"`
	ClusterID string `json:"clusterId"`
	Region    string `json:"region"`
}

type ClusterResponse struct {
	Message string `json:"message"`
}

func CreateECSHandler(w http.ResponseWriter, r *http.Request) {
	var req ClusterRequest
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create ECS service client
	svc := ecs.New(sess)

	// Choose the infrastructure provider
	var capacityProviders []*ecs.CapacityProviderStrategyItem
	var capacityProviders1 []*string
	switch req.Infrastructure {
	case "Fargate":
		capacityProviders = []*ecs.CapacityProviderStrategyItem{
			{
				CapacityProvider: aws.String("FARGATE"),
				Weight:           aws.Int64(1),
			},
			{
				CapacityProvider: aws.String("FARGATE_SPOT"),
				Weight:           aws.Int64(1),
			},
		}
	case "EC2":
		capacityProviders = []*ecs.CapacityProviderStrategyItem{
			{
				CapacityProvider: aws.String("EC2"),
				Weight:           aws.Int64(1),
			},
		}
		capacityProviders1 = []*string{
			aws.String("EC2"),
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid infrastructure choice")
		return
	}

	// Create the ECS cluster
	_, err = svc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName:                     aws.String(req.ClusterID),
		CapacityProviders:               capacityProviders1,
		DefaultCapacityProviderStrategy: capacityProviders, // Choose default capacity provider strategy
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating ECS cluster: %v", err)
		return
	}

	// Send success response
	resp := ClusterResponse{Message: "ECS cluster created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteECSHandler(w http.ResponseWriter, r *http.Request) {
	var req ClusterDeleteRequest
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create ECS service client
	svc := ecs.New(sess)

	// Perform the deletion operation
	_, err = svc.DeleteCluster(&ecs.DeleteClusterInput{
		Cluster: aws.String(req.ClusterID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting ECS cluster: %v", err)
		return
	}

	// Send success response
	resp := ClusterResponse{Message: "ECS cluster deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListClusterRequest struct {
	AccountID int    `json:"accountId"`
	Region    string `json:"region"`
}

func ListClustersHandler(w http.ResponseWriter, r *http.Request) {
	var req ListClusterRequest
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create ECS service client
	svc := ecs.New(sess)

	// List all ECS clusters
	resp, err := svc.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing ECS clusters: %v", err)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
