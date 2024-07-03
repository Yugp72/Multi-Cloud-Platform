package gcp_network

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"google.golang.org/api/compute/v1"
)

type SubnetRequest struct {
	ProjectID   int    `json:"projectId"`
	NetworkName string `json:"networkName"`
	SubnetName  string `json:"subnetName"`
	Region      string `json:"region"`
	IPCIDRRange string `json:"ipCidrRange"`
	Token       string `json:"token"`
}

func CreateSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req SubnetRequest
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

	subnet := &compute.Subnetwork{
		Name:        req.SubnetName,
		Region:      req.Region,
		IpCidrRange: req.IPCIDRRange,
		Network:     req.NetworkName,
	}

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)

	_, err = computeService.Subnetworks.Insert(project, req.Region, subnet).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating subnet: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Subnet created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ListSubnetsHandler(w http.ResponseWriter, r *http.Request) {
	var req SubnetRequest
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

	// List subnets
	subnets, err := computeService.Subnetworks.List(project, req.Region).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing subnets: %v", err)
		return
	}

	// Process the list of subnets
	var subnetNames []string
	for _, subnet := range subnets.Items {
		subnetNames = append(subnetNames, subnet.Name)
	}

	resp := struct {
		Subnets []string `json:"subnets"`
	}{
		Subnets: subnetNames,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req SubnetRequest
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

	// Delete the subnet
	_, err = computeService.Subnetworks.Delete(project, req.Region, req.SubnetName).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting subnet: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Subnet deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
