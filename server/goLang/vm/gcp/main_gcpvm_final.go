package gcp_compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db "btep.project/databaseConnection"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// InstanceRequest represents the JSON request structure for GCP instance operations
type InstanceRequest struct {
	MachineType string `json:"machineType"`
	Image       string `json:"image"`
	Zone        string `json:"zone"`
	Name        string `json:"name"`
	AccountID   int    `json:"accountID"`
	Token       string `json:"token"`
}

// InstanceResponse represents the JSON response structure for GCP instance operations
type InstanceResponse struct {
	Message     string   `json:"message"`
	InstanceIDs []string `json:"instanceIDs,omitempty"`
}

func initComputeService(token string) (*compute.Service, error) {
	ctx := context.Background()
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	computeService, err := compute.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating compute service: %v", err)
	}

	return computeService, nil
}

// CreateInstanceHandler handles POST requests to create a GCP instance
func CreateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var req InstanceRequest
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
	computeService, err := initComputeService(req.Token)

	// Generate a unique name for the instance
	instanceName := generateValidInstanceName()

	// Create the GCP instance
	_, err = computeService.Instances.Insert(cloudAccount.ProjectID.String, req.Zone, &compute.Instance{
		Name:        instanceName,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", req.Zone, req.MachineType),
		Disks: []*compute.AttachedDisk{
			{
				Source:     fmt.Sprintf("projects/%s/zones/%s/disks/%s", cloudAccount.ProjectID.String, req.Zone, instanceName),
				Boot:       true,
				AutoDelete: true,
			},
		},
	}).Do()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating GCP instance: %v", err)
		return
	}

	resp := InstanceResponse{Message: fmt.Sprintf("GCP instance created with name: %s", instanceName)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Function to generate a valid GCP instance name
func generateValidInstanceName() string {
	return "instance-" + uuid.New().String()[:8] // Example: instance-abcdef12
}

// InstanceListRequest represents the JSON request structure for listing GCP instances
type InstanceListRequest struct {
	AccountID int    `json:"accountID"`
	Token     string `json:"token"`
	Zone      string `json:"zone"`
}

// TerminalRequest represents the JSON request structure for terminating a GCP instance
type TerminalRequest struct {
	AccountID  int    `json:"accountID"`
	InstanceID string `json:"instanceID"`
	Token      string `json:"token"`
	Zone       string `json:"zone"`
}

type InstanceDetails struct {
	InstanceID       string   `json:"instance_id"`        // Use the GCP instance name or ID
	InstanceType     string   `json:"instance_type"`      // Specify the machine type or instance size
	LaunchTime       string   `json:"launch_time"`        // Creation timestamp or start time of the instance
	PrivateIPAddress string   `json:"private_ip_address"` // Private IP address of the instance
	PublicIPAddress  string   `json:"public_ip_address"`  // Public IP address of the instance
	AvailabilityZone string   `json:"availability_zone"`  // GCP doesn't have zones exactly like AWS, but you can derive this information from the instance's selfLink
	StateCode        string   `json:"state_code"`         // State of the instance (e.g., RUNNING, TERMINATED)
	SecurityGroups   []string `json:"security_groups"`    // GCP uses firewall rules for security, so you can list the associated firewall rules here
	PlatformDetails  string   `json:"platform_details"`   // Any specific platform details related to the instance
	Tags             []string `json:"tags"`               // Tags or labels associated with the instance
}

type ListRequest struct {
	AccountID int    `json:"accountID"`
	Token     string `json:"token"`
}

func ListInstancesHandler(w http.ResponseWriter, r *http.Request) {
	var req ListRequest
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

	computeService, err := initComputeService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing GCP compute service: %v", err)
		return
	}

	// List GCP instances
	instancesList, err := computeService.Instances.AggregatedList(cloudAccount.ProjectID.String).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing GCP instances: %v", err)
		return
	}

	// Extract instance details
	var instances []InstanceDetails
	for _, zone := range instancesList.Items {
		for _, instance := range zone.Instances {
			instanceDetail := InstanceDetails{
				InstanceID:       instance.Name,
				InstanceType:     extractInstanceType(instance.MachineType),
				LaunchTime:       instance.CreationTimestamp,
				PrivateIPAddress: instance.NetworkInterfaces[0].NetworkIP,              // Assuming there's only one network interface
				PublicIPAddress:  instance.NetworkInterfaces[0].AccessConfigs[0].NatIP, // Assuming there's only one access config
				AvailabilityZone: extractZoneFromSelfLink(instance.SelfLink),           // Extract zone from the self link
				StateCode:        instance.Status,
			}
			instances = append(instances, instanceDetail)
		}
	}

	// Encode instance details into JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(instances)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding JSON response: %v", err)
		return
	}
}

func extractInstanceType(fullTypeURL string) string {
	parts := strings.Split(fullTypeURL, "/")
	return parts[len(parts)-1]
}

// Function to extract zone from GCP instance self link
func extractZoneFromSelfLink(selfLink string) string {
	parts := strings.Split(selfLink, "/")
	return parts[len(parts)-3] // Assuming the zone is the third-to-last part of the self link
}

// TerminateInstanceHandler handles POST requests to terminate a GCP instance
func TerminateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var req TerminalRequest
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
	computeService, err := initComputeService(req.Token)
	// Delete the GCP instance
	op, err := computeService.Instances.Delete(cloudAccount.ProjectID.String, req.Zone, req.InstanceID).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error terminating GCP instance: %v", err)
		return
	}
	fmt.Println(op)

	resp := InstanceResponse{Message: fmt.Sprintf("GCP instance terminated: %s", req.InstanceID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
