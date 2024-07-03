package gcp_network

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"google.golang.org/api/compute/v1"
)

type CloudRouterRequest struct {
	ProjectID       int    `json:"projectId"`
	RouterName      string `json:"routerName"`
	Region          string `json:"region"`
	Network         string `json:"network"`
	ASN             int64  `json:"asn"`
	Description     string `json:"description"`
	BgpPeeringName  string `json:"bgpPeeringName"`
	InterfaceName   string `json:"interfaceName"`
	IPAddress       string `json:"ipAddress"`
	PeerIPAddress   string `json:"peerIpAddress"`
	PeerASN         int64  `json:"peerAsn"`
	AdvertisedRoute string `json:"advertisedRoute"`
	Token           string `json:"token"`
}

func CreateCloudRouterHandler(w http.ResponseWriter, r *http.Request) {
	var req CloudRouterRequest
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

	router := &compute.Router{
		Name:    req.RouterName,
		Region:  req.Region,
		Network: req.Network,
		Bgp: &compute.RouterBgp{
			Asn:                req.ASN,
			AdvertiseMode:      "DEFAULT",
			AdvertisedGroups:   []string{"ALL_SUBNETS"},
			AdvertisedIpRanges: []*compute.RouterAdvertisedIpRange{},
		},
		Description: req.Description,
	}

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)

	_, err = computeService.Routers.Insert(project, req.Region, router).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating cloud router: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Cloud router created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func ListCloudRoutersHandler(w http.ResponseWriter, r *http.Request) {
	var req CloudRouterRequest
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

	// List cloud routers
	routersList, err := computeService.Routers.List(project, req.Region).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing cloud routers: %v", err)
		return
	}

	// Extract router names and IDs and respond
	var routers []string
	for _, router := range routersList.Items {
		routers = append(routers, router.Name)
	}

	resp := struct {
		Routers []string `json:"routers"`
	}{
		Routers: routers,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteCloudRouterHandler(w http.ResponseWriter, r *http.Request) {
	var req CloudRouterRequest
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

	// Get the router ID based on the router name
	routerID := fmt.Sprintf("global/routers/%s", req.RouterName)

	// Delete the router
	temp, err := computeService.Routers.Delete(project, req.Region, routerID).Do()
	fmt.Println(temp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting cloud router: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Cloud router '%s' deleted successfully", req.RouterName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
