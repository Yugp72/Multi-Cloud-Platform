package aws_eks

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

type ClusterRequest struct {
	AccountID         int      `json:"accountId"`
	Region            string   `json:"region"`
	ClusterID         string   `json:"clusterId"`
	KubernetesVersion string   `json:"kubernetesVersion"`
	RoleARN           string   `json:"roleARN"`
	Subnets           []string `json:"subnets"`
	ClusterLogging    string   `json:"clusterLogging"`
	EncryptionConfig  string   `json:"encryptionConfig"`
}

func CreateEKSHandler(w http.ResponseWriter, r *http.Request) {
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create EKS service client
	svc := eks.New(sess)

	var logSetups []*eks.LogSetup
	for _, logType := range req.ClusterLogging {
		logTypeStr := string(logType) // Convert rune to string
		logSetup := &eks.LogSetup{
			Enabled: aws.Bool(true),
			Types:   []*string{aws.String(logTypeStr)},
		}
		logSetups = append(logSetups, logSetup)
	}
	// Create the EKS cluster with input parameters
	_, err = svc.CreateCluster(&eks.CreateClusterInput{
		Name:               aws.String(req.ClusterID),
		Version:            aws.String(req.KubernetesVersion),
		RoleArn:            aws.String(req.RoleARN),
		ResourcesVpcConfig: &eks.VpcConfigRequest{SubnetIds: aws.StringSlice(req.Subnets)},
		EncryptionConfig:   []*eks.EncryptionConfig{{Provider: &eks.Provider{KeyArn: aws.String(req.EncryptionConfig)}}},
		Logging:            &eks.Logging{ClusterLogging: logSetups},
		ClientRequestToken: nil,
		KubernetesNetworkConfig: &eks.KubernetesNetworkConfigRequest{
			ServiceIpv4Cidr: aws.String("")},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating EKS cluster: %v", err)
		return
	}

	// Send success response
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "EKS cluster created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

type UpdateEKSRequest struct {
	AccountID         int    `json:"accountId"`
	ClusterID         string `json:"clusterId"`
	KubernetesVersion string `json:"kubernetesVersion"`
	Region            string `json:"region"`
}

func UpdateEKSHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateEKSRequest
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create EKS service client
	svc := eks.New(sess)

	// Prepare parameters for updating the EKS cluster
	updateInput := &eks.UpdateClusterVersionInput{
		Name:               aws.String(req.ClusterID),
		Version:            aws.String(req.KubernetesVersion),
		ClientRequestToken: nil,
	}

	// Update the EKS cluster
	_, err = svc.UpdateClusterVersion(updateInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating EKS cluster: %v", err)
		return
	}

	// Send success response
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "EKS cluster updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type DeleteEKSRequest struct {
	ClusterID string `json:"clusterId"`
	AccountID int    `json:"accountId"`
	Region    string `json:"region"`
}

func DeleteEKSHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteEKSRequest
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create EKS service client
	svc := eks.New(sess)

	// Prepare parameters for deleting the EKS cluster
	deleteInput := &eks.DeleteClusterInput{
		Name: aws.String(req.ClusterID),
	}

	// Delete the EKS cluster
	_, err = svc.DeleteCluster(deleteInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting EKS cluster: %v", err)
		return
	}

	// Send success response
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "EKS cluster deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListEKSRequest struct {
	AccountID int    `json:"accountId"`
	Region    string `json:"region"`
}

func ListEKSHandler(w http.ResponseWriter, r *http.Request) {
	var req ListEKSRequest
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region),
		Credentials: credentials.NewStaticCredentials(
			cloudAccount.AccessKey.String,
			cloudAccount.SecretKey.String,
			"",
		),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session")
		return
	}

	// Create EKS service client
	svc := eks.New(sess)

	// List EKS clusters
	result, err := svc.ListClusters(&eks.ListClustersInput{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing EKS clusters: %v", err)
		return
	}

	// Extract cluster names from the response
	var clusters []string
	for _, cluster := range result.Clusters {
		clusters = append(clusters, *cluster)
	}

	// Send success response
	resp := struct {
		Clusters []string `json:"clusters"`
	}{
		Clusters: clusters,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
