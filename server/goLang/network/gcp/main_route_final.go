package gcp_network

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"google.golang.org/api/compute/v1"
)

type RouteRequest struct {
	ProjectID   int      `json:"projectId"`
	RouteName   string   `json:"routeName"`
	Destination string   `json:"destination"`
	NextHopIP   string   `json:"nextHopIp"`
	Network     string   `json:"network"`
	Priority    int64    `json:"priority"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"` // Corrected to []string
	Token       string   `json:"token"`
}

func CreateRouteHandler(w http.ResponseWriter, r *http.Request) {
	var req RouteRequest
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

	route := &compute.Route{
		Name:        req.RouteName,
		DestRange:   req.Destination,
		NextHopIp:   req.NextHopIP,
		Network:     req.Network,
		Priority:    req.Priority,
		Description: req.Description,
		Tags:        req.Tags,
	}

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)

	_, err = computeService.Routes.Insert(project, route).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating route: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Route created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteRouteHandler(w http.ResponseWriter, r *http.Request) {
	var req RouteRequest
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

	// Delete the route
	_, err = computeService.Routes.Delete(project, req.RouteName).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting route: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Route deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListRouteRequest struct {
	ProjectID int    `json:"projectId"`
	Token     string `json:"token"`
}

func ListRoutesHandler(w http.ResponseWriter, r *http.Request) {
	var req ListRouteRequest
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

	routes, err := computeService.Routes.List(project).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing routes: %v", err)
		return
	}

	resp := struct {
		Routes []*compute.Route `json:"routes"`
	}{
		Routes: routes.Items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
