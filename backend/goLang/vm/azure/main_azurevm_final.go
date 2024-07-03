package azure_vms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

type VMImage struct {
	Publisher string `json:"publisher"`
	Offer     string `json:"offer"`
	SKU       string `json:"sku"`
	Version   string `json:"version"`
}

type VMRequest struct {
	VMName             string   `json:"vmName"`
	ResourceGroup      string   `json:"resourceGroup"`
	SubscriptionID     string   `json:"subscriptionID"`
	Token              string   `json:"token"`
	Image              VMImage  `json:"image"`
	VMSize             string   `json:"vmSize"`
	KeyPairName        string   `json:"keyPairName"`
	InboundPorts       []int    `json:"inboundPorts"`
	Region             string   `json:"region"`
	Zones              []string `json:"zones"`
	AccountID          int      `json:"accountID"`
	IdentityName       string   `json:"identityName"`
	NetworkInterfaceID string   `json:"networkInterfaceID"`
}

// ListVMsRequest represents the JSON request structure for listing VMs
type ListVMsRequest struct {
	Token     string `json:"token"`
	AccountID int    `json:"accountID"`
}

// DeleteVMRequest represents the JSON request structure for deleting a VM
type DeleteVMRequest struct {
	VMName         string `json:"vmName"`
	ResourceGroup  string `json:"resourceGroup"`
	SubscriptionID string `json:"subscriptionID"`
	Token          string `json:"token"`
}

// VMResponse represents the JSON response structure
type VMResponse struct {
	Message string `json:"message"`
}

type tokenAuthorizer struct {
	token string
}

func (ta tokenAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Header.Set("Authorization", "Bearer "+ta.token)
			return r, nil
		})
	}
}

func initComputeClient(subscriptionID string, token string) (compute.VirtualMachinesClient, error) {
	client := compute.NewVirtualMachinesClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{}
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}

// CreateVMHandler creates a new Azure VM
func CreateVMHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req VMRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Create Azure VM client
	client, err := initComputeClient(req.SubscriptionID, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure VM client: %v", err)
		return
	}
	networkInterfaceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s", req.SubscriptionID, req.ResourceGroup, req.NetworkInterfaceID)

	vmParameters := compute.VirtualMachine{
		Location: &req.Region,
		Name:     &req.VMName,
		Identity: &compute.VirtualMachineIdentity{
			Type: compute.ResourceIdentityTypeUserAssigned,
			UserAssignedIdentities: map[string]*compute.UserAssignedIdentitiesValue{
				fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s", req.SubscriptionID, req.ResourceGroup, req.IdentityName): {},
			},
		},
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(req.VMSize),
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: &compute.ImageReference{
					Publisher: &req.Image.Publisher,
					Offer:     &req.Image.Offer,
					Version:   &req.Image.Version,
					Sku:       &req.Image.SKU,
				},
			},
			OsProfile: &compute.OSProfile{
				ComputerName: &req.VMName,
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: &networkInterfaceID,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
	}

	// Create the VM
	_, err = client.CreateOrUpdate(context.Background(), req.ResourceGroup, req.VMName, vmParameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Azure VM: %v", err)
		return
	}

	// Send success response
	resp := VMResponse{Message: "Azure VM created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListVMsHandler lists Azure VMs across all accessible subscriptions
func ListVMsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req ListVMsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Get all subscription IDs
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)

	// subscriptionIDs, err := GetAllSubscriptionIDs(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to get subscription IDs: %v", err)
		return
	}

	// Create Azure VM client
	client, err := initComputeClient(cloudAccount.SubscriptionID.String, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure VM client: %v", err)
		return
	}
	responce, err := client.ListAll(context.Background(), "", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure VM client: %v", err)
		return
	}
	fmt.Println(responce)
	// List VMs across all subscriptions
	// var allVMs []compute.VirtualMachine
	// for _, subscriptionID := range subscriptionIDs {
	// 	vmClient := compute.NewVirtualMachinesClient(subscriptionID)
	// 	vmClient.Authorizer = client.Authorizer // Use the same authorizer
	// 	vms, err := vmClient.ListAll(context.Background(), "", "")
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Error listing VMs in subscription %s: %v\n", subscriptionID, err)
	// 		continue
	// 	}
	// 	allVMs = append(allVMs, vms.Values()...)
	// }

	// Send success response with all VMs list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responce)
}
func DeleteVMHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req DeleteVMRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Create Azure VM client
	client, err := initComputeClient(req.SubscriptionID, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure VM client: %v", err)
		return
	}

	// Delete the VM
	_, err = client.Delete(context.Background(), req.ResourceGroup, req.VMName, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting Azure VM: %v", err)
		return
	}

	// Send success response
	resp := VMResponse{Message: "Azure VM deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type TokenRequest struct {
	Token string `json:"token"`
}
type SubscriptionIDsResponse struct {
	SubscriptionIDs []string `json:"subscriptionIDs"`
}

// GetSubscriptionIDsHandler handles the request to retrieve subscription IDs
func GetSubscriptionIDsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req TokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Retrieve subscription IDs
	subscriptionIDs, err := GetAllSubscriptionIDs(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving subscription IDs: %v", err)
		return
	}

	// Create response
	resp := SubscriptionIDsResponse{
		SubscriptionIDs: subscriptionIDs,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func withBearerToken(token string) autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Header.Set("Authorization", "Bearer "+token)
			return r, nil
		})
	}
}
func GetAllSubscriptionIDs(token string) ([]string, error) {
	client := subscriptions.NewClient()
	client.RequestInspector = withBearerToken(token)

	var subscriptionIDs []string
	subList, err := client.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %v", err)
	}

	for _, subscription := range subList.Values() {
		subscriptionIDs = append(subscriptionIDs, *subscription.SubscriptionID)
	}

	return subscriptionIDs, nil
}
