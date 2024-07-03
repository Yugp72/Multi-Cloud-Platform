package gcp_network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"golang.org/x/oauth2"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func initComputeService(token string) (*compute.Service, error) {
	ctx := context.Background()
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	computeService, err := compute.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating compute service: %v", err)
	}

	return computeService, nil
}

type NetworkRequest struct {
	ProjectID   int    `json:"projectId"`
	NetworkName string `json:"networkName"`
	Token       string `json:"token"`
	AccountID   int    `json:"account"`
}

func CreateNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var req NetworkRequest
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

	network := &compute.Network{
		Name:                  req.NetworkName,
		AutoCreateSubnetworks: true, // Use subnet mode
		Description:           "Created via API",
	}

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing compute service: %v", err)
		return
	}

	_, err = computeService.Networks.Insert(project, network).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating network: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Network created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ListNetworksHandler(w http.ResponseWriter, r *http.Request) {
	var req NetworkRequest
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

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing compute service: %v", err)
		return
	}

	// List networks
	networks, err := computeService.Networks.List(project).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing networks: %v", err)
		return
	}

	// Process the list of networks
	var networkNames []string
	for _, network := range networks.Items {
		networkNames = append(networkNames, network.Name)
	}

	resp := struct {
		Networks []string `json:"networks"`
	}{
		Networks: networkNames,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var req NetworkRequest
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

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing compute service: %v", err)
		return
	}

	// Delete the network
	_, err = computeService.Networks.Delete(project, req.NetworkName).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting network: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Network deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
