package gcp_network

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"google.golang.org/api/compute/v1"
)

type FirewallRuleRequest struct {
	ProjectID     int      `json:"projectId"`
	FirewallName  string   `json:"firewallName"`
	Network       string   `json:"network"`
	SourceCIDR    string   `json:"sourceCidr"`
	AllowedPorts  string   `json:"allowedPorts"`
	Direction     string   `json:"direction"`
	Description   string   `json:"description"`
	Priority      int64    `json:"priority"`
	SourceTags    []string `json:"sourceTags"` // Corrected to []string
	DestinationIP string   `json:"destinationIP"`
	Token         string   `json:"token"`
}

func CreateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req FirewallRuleRequest
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

	firewall := &compute.Firewall{
		Name:         req.FirewallName,
		Network:      req.Network,
		SourceRanges: []string{req.SourceCIDR},
		Allowed: []*compute.FirewallAllowed{
			{
				IPProtocol: "tcp",
				Ports:      []string{req.AllowedPorts},
			},
		},
		Direction:         req.Direction,
		Description:       req.Description,
		Priority:          req.Priority,
		SourceTags:        req.SourceTags,
		DestinationRanges: []string{req.DestinationIP},
	}

	project := cloudAccount.ProjectID.String
	computeService, err := initComputeService(req.Token)

	_, err = computeService.Firewalls.Insert(project, firewall).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating firewall rule: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Firewall rule created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func ListFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	var req FirewallRuleRequest
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

	// List firewall rules
	firewallRulesList, err := computeService.Firewalls.List(project).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing firewall rules: %v", err)
		return
	}

	// Extract firewall rule names and respond
	var firewallRules []string
	for _, rule := range firewallRulesList.Items {
		firewallRules = append(firewallRules, rule.Name)
	}

	resp := struct {
		FirewallRules []string `json:"firewallRules"`
	}{
		FirewallRules: firewallRules,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req FirewallRuleRequest
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

	// Delete the firewall rule
	temp, err := computeService.Firewalls.Delete(project, req.FirewallName).Do()
	fmt.Println(temp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting firewall rule: %v", err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Firewall rule '%s' deleted successfully", req.FirewallName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
