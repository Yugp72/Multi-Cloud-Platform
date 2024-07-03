package gcp_kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

// ClusterRequest represents the JSON request structure for GCP cluster operations
type ClusterRequest struct {
	AccountID         int    `json:"accountID"`
	Zone              string `json:"zone"`
	ClusterName       string `json:"clusterName"`
	KubernetesVersion string `json:"kubernetesVersion"`
	NodeCount         int64  `json:"nodeCount"`
	InstanceType      string `json:"instanceType"`
	Token             string `json:"token"`
}

type DeleteClusterRequest struct {
	AccountID   int    `json:"accountID"`
	Zone        string `json:"zone"`
	ClusterName string `json:"clusterName"`
	Token       string `json:"token"`
}

// ClusterResponse represents the JSON response structure for GCP cluster operations
type ClusterResponse struct {
	Message   string `json:"message"`
	ClusterID string `json:"clusterID,omitempty"`
}

func initContainerService(token string) (*container.Service, error) {
	ctx := context.Background()
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	containerService, err := container.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating container service: %v", err)
	}

	return containerService, nil
}

// CreateClusterHandler handles POST requests to create a GCP Kubernetes cluster
func CreateClusterHandler(w http.ResponseWriter, r *http.Request) {
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
	containerService, err := initContainerService(req.Token)

	// Generate a unique ID for the cluster
	clusterID := uuid.New().String()

	// Create the GCP Kubernetes cluster
	ctx := context.Background()
	_, err = containerService.Projects.Zones.Clusters.Create(cloudAccount.ProjectID.String, req.Zone, &container.CreateClusterRequest{
		Cluster: &container.Cluster{
			Name:             req.ClusterName,
			InitialNodeCount: req.NodeCount,
			NodeConfig: &container.NodeConfig{
				MachineType: req.InstanceType,
				DiskSizeGb:  100, // Example disk size
			},
			MasterAuth: &container.MasterAuth{
				Username: "admin", // Example admin username
				Password: "admin", // Example admin password
			},
			Locations: []string{req.Zone},
		},
	}).Context(ctx).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCP Kubernetes cluster: %v", err)
		return
	}

	resp := ClusterResponse{Message: fmt.Sprintf("GCP Kubernetes cluster created with ID: %s", clusterID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteClusterHandler handles POST requests to delete a GCP Kubernetes cluster
func DeleteClusterHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteClusterRequest
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
	containerService, err := initContainerService(req.Token)

	// Delete the GCP Kubernetes cluster
	_, err = containerService.Projects.Zones.Clusters.Delete(cloudAccount.ProjectID.String, req.Zone, req.ClusterName).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting GCP Kubernetes cluster: %v", err)
		return
	}

	resp := ClusterResponse{Message: fmt.Sprintf("GCP Kubernetes cluster deleted: %s", req.ClusterName)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListClusterRequest struct {
	AccountID int    `json:"accountID"`
	Zone      string `json:"zone"`
	Token     string `json:"token"`
}

// ListClustersHandler handles POST requests to list GCP Kubernetes clusters
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
	containerService, err := initContainerService(req.Token)

	// List the GCP Kubernetes clusters
	ctx := context.Background()
	clusters, err := containerService.Projects.Zones.Clusters.List(cloudAccount.ProjectID.String, req.Zone).Context(ctx).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing GCP Kubernetes clusters: %v", err)
		return
	}

	var clusterNames []string
	for _, cluster := range clusters.Clusters {
		clusterNames = append(clusterNames, cluster.Name)
	}

	resp := ClusterResponse{Message: "Successfully retrieved GCP Kubernetes clusters", ClusterID: fmt.Sprintf("%v", clusterNames)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
