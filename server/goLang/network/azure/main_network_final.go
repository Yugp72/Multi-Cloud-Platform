package azure_network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest"
)

type NetworkRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	NetworkName    string `json:"networkName"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
	Prefix         string `json:"prefix"`
	Token          string `json:"token"`
}
type DeleteNetworkRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
	NetworkName    string `json:"networkName"`
	Token          string `json:"token"`
}

type ListNetworkRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
	Token          string `json:"token"`
}

type ListNetworkResponse struct {
	Message  string   `json:"message"`
	Networks []string `json:"networks,omitempty"`
}
type NetworkResponse struct {
	Message  string   `json:"message"`
	Networks []string `json:"networks"`
}
type tokenAuthorizer struct {
	token string
}

//	func (ta tokenAuthorizer) WithAuthorization() autorest.PrepareDecorator {
//		return func(p autorest.Preparer) autorest.Preparer {
//			return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
//				r.Header.Set("Authorization", "Bearer "+ta.token)
//				return r, nil
//			})
//		}
//	}
func initNetworkClient(subscriptionID string, token string) (network.VirtualNetworksClient, error) {
	client := network.NewVirtualNetworksClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{} // We manually insert the token
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}
func CreateNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var req NetworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initNetworkClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := network.VirtualNetwork{
		Location: &req.Location,
		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
			AddressSpace: &network.AddressSpace{
				AddressPrefixes: &[]string{req.Prefix},
			},
		},
	}
	ctx := context.Background()
	_, err = client.CreateOrUpdate(ctx, req.ResourceGroup, req.NetworkName, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create network: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(NetworkResponse{Message: "Network created successfully"})
}

func ListNetworksHandler(w http.ResponseWriter, r *http.Request) {
	var req ListNetworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initNetworkClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	networkList, err := client.List(ctx, req.ResourceGroup)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list networks: %v", err), http.StatusInternalServerError)
		return
	}

	var networks []string
	for _, vnet := range networkList.Values() {
		networks = append(networks, *vnet.Name)
	}

	resp := NetworkResponse{
		Message:  "Successfully retrieved networks",
		Networks: networks,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DeleteNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteNetworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initNetworkClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	future, err := client.Delete(ctx, req.ResourceGroup, req.NetworkName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete network: %v", err), http.StatusInternalServerError)
		return
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to complete network deletion operation: %v", err), http.StatusInternalServerError)
		return
	}

	resp := NetworkResponse{
		Message: "Network deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
