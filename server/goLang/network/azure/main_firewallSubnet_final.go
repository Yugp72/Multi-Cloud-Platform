package azure_network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest"
)

func (ta tokenAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Header.Set("Authorization", "Bearer "+ta.token)
			return r, nil
		})
	}
}
func initSubnetClient1(subscriptionID string, token string) (network.SubnetsClient, error) {
	client := network.NewSubnetsClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{} // We manually insert the token
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}

func initFirewallClient(subscriptionID string, token string) (network.AzureFirewallsClient, error) {
	client := network.NewAzureFirewallsClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{} // We manually insert the token
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}

type SubnetRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
	NetworkName    string `json:"networkName"`
	Prefix         string `json:"prefix"`
	Token          string `json:"token"`
	SubnetName     string `json:"subnetName"`
}

func CreateSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initSubnetClient1(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := network.Subnet{
		Name: &req.SubnetName,
		SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
			AddressPrefix: &req.Prefix,
		},
	}
	ctx := context.Background()
	_, err = client.CreateOrUpdate(ctx, req.ResourceGroup, req.NetworkName, req.SubnetName, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create subnet: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(NetworkResponse{Message: "Subnet created successfully"})
}

type FirewallRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
	Token          string `json:"token"`
	FirewallName   string `json:"firewallName"`
}

func CreateFirewallHandler(w http.ResponseWriter, r *http.Request) {
	var req FirewallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initFirewallClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := network.AzureFirewall{
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			Sku: &network.AzureFirewallSku{
				Name: network.AZFWHub,
			},
			ThreatIntelMode: network.AzureFirewallThreatIntelModeAlert,
		},
	}
	ctx := context.Background()
	_, err = client.CreateOrUpdate(ctx, req.ResourceGroup, req.FirewallName, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create firewall: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(NetworkResponse{Message: "Firewall created successfully"})
}

type DeleteSubnetRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
	NetworkName    string `json:"networkName"`
	Token          string `json:"token"`
	SubnetName     string `json:"subnetName"`
}

func DeleteSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteSubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initSubnetClient1(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	resp, err := client.Delete(ctx, req.ResourceGroup, req.NetworkName, req.SubnetName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete subnet: %v", err), http.StatusInternalServerError)
		return
	}

	if resp.Response().StatusCode != http.StatusAccepted {
		http.Error(w, fmt.Sprintf("Delete operation not accepted: %v", resp.Status), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(NetworkResponse{Message: "Subnet deleted successfully"})
}

type DeleteFirewallRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
	Token          string `json:"token"`
	FirewallName   string `json:"firewallName"`
}

func DeleteFirewallHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteFirewallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initFirewallClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	resp, err := client.Delete(ctx, req.ResourceGroup, req.FirewallName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete firewall: %v", err), http.StatusInternalServerError)
		return
	}

	if resp.Response().StatusCode != http.StatusAccepted {
		http.Error(w, fmt.Sprintf("Delete operation not accepted: %v", resp.Status), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(NetworkResponse{Message: "Firewall deleted successfully"})
}

// list firewall
type FirewallResponse struct {
	Message   string   `json:"message"`
	Firewalls []string `json:"firewalls"`
}

type ListFirewallRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	Token          string `json:"token"`
}

func ListFirewallHandler(w http.ResponseWriter, r *http.Request) {
	var req ListFirewallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initFirewallClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	firewallList, err := client.ListAll(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list firewalls: %v", err), http.StatusInternalServerError)
		return
	}

	var firewalls []string
	for _, firewall := range firewallList.Values() {
		firewalls = append(firewalls, *firewall.Name)
	}

	resp := FirewallResponse{
		Message:   "Successfully retrieved firewalls",
		Firewalls: firewalls,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type SubnetResponse struct {
	Message string   `json:"message"`
	Subnets []string `json:"subnets"`
}

type ListSubnetRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	NetworkName    string `json:"networkName"`
	Token          string `json:"token"`
}

func ListSubnetHandler(w http.ResponseWriter, r *http.Request) {
	var req ListSubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initSubnetClient1(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	subnetList, err := client.List(ctx, req.NetworkName, req.SubscriptionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list subnets: %v", err), http.StatusInternalServerError)
		return
	}

	var subnets []string
	for _, subnet := range subnetList.Values() {
		subnets = append(subnets, *subnet.Name)
	}

	resp := SubnetResponse{
		Message: "Successfully retrieved subnets",
		Subnets: subnets,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
