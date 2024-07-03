package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	db "btep.project/databaseConnection"

	aws_dynamodb "btep.project/DataBase/aws"
	azure_cosmosdb "btep.project/DataBase/azure"
	gcp_firebase "btep.project/DataBase/gcp"
	aws_s3 "btep.project/Storage/aws"
	azure_storage "btep.project/Storage/azure"
	gcp_gcs "btep.project/Storage/gcp"
	aws_vpc "btep.project/network/aws"
	azure_network "btep.project/network/azure"
	gcp_network "btep.project/network/gcp"
	aws_ecs "btep.project/serverless/aws/ecs"
	aws_eks "btep.project/serverless/aws/eks"
	aws_lambda "btep.project/serverless/aws/lambda"
	azure_functions "btep.project/serverless/azure"
	gcp_cloudrun "btep.project/serverless/gcp/cloudrun"
	gcp_kubernetes "btep.project/serverless/gcp/gke"
	aws_ec2 "btep.project/vm/aws"
	azure_vms "btep.project/vm/azure"
	gcp_compute "btep.project/vm/gcp"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	googleAuth "golang.org/x/oauth2/google"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

	router.HandleFunc("/auth/google/login", handleGoogleLogin).Methods("POST")
	router.HandleFunc("/auth/google/callback", handleGoogleCallback).Methods("GET")

	// router.HandleFunc("/gcp/auth", auth_gcp.handleGetAccessToken).Methods("POST")
	router.HandleFunc("/auth/azure/login", handleAzureLogin).Methods("POST")
	router.HandleFunc("/auth/azure/callback", handleAzureCallback).Methods("GET")
	// AWS S3
	router.HandleFunc("/aws/s3/createBucket", aws_s3.CreateBucketHandler).Methods("POST")
	router.HandleFunc("/aws/s3/uploadObject", aws_s3.UploadObjectHandler).Methods("POST")
	router.HandleFunc("/aws/s3/getObject", aws_s3.GetObjectHandler).Methods("GET")
	router.HandleFunc("/aws/s3/deleteObject", aws_s3.DeleteObjectHandler).Methods("POST")
	router.HandleFunc("/aws/s3/deleteBucket", aws_s3.DeleteBucketHandler).Methods("POST")
	router.HandleFunc("/aws/s3/listBuckets", aws_s3.ListBucketsHandler).Methods("POST")

	// GCP GCS
	router.HandleFunc("/gcp/gcs/createBucket", gcp_gcs.CreateBucketHandler).Methods("POST")
	router.HandleFunc("/gcp/gcs/uploadObject", gcp_gcs.UploadObjectHandler).Methods("POST")
	router.HandleFunc("/gcp/gcs/getObject", gcp_gcs.GetObjectHandler).Methods("POST")
	router.HandleFunc("/gcp/gcs/deleteObject", gcp_gcs.DeleteObjectHandler).Methods("POST")
	router.HandleFunc("/gcp/gcs/deleteBucket", gcp_gcs.DeleteBucketHandler).Methods("POST")
	router.HandleFunc("/gcp/gcs/listBuckets", gcp_gcs.ListBucketsHandler).Methods("POST")

	// Azure Storage
	router.HandleFunc("/azure/storage/createAccount", azure_storage.CreateStorageAccountHandler).Methods("POST")
	router.HandleFunc("/azure/storage/deleteAccount", azure_storage.DeleteStorageAccountHandler).Methods("POST")
	router.HandleFunc("/azure/storage/updateStorageAccount", azure_storage.UpdateStorageAccountHandler).Methods("POST")
	router.HandleFunc("/azure/storage/listAccounts", azure_storage.ListStorageAccountsHandler).Methods("POST")
	router.HandleFunc("/azure/storage/getObjects", azure_storage.GetObjectHandler).Methods("POST")
	router.HandleFunc("/azure/storage/uploadObjects", azure_storage.UploadBlobHandler).Methods("POST")

	// EC2
	router.HandleFunc("/aws/ec2/createInstance", aws_ec2.CreateInstanceHandler).Methods("POST")
	router.HandleFunc("/aws/ec2/listInstances", aws_ec2.ListInstancesHandler).Methods("POST")
	router.HandleFunc("/aws/ec2/terminateInstance", aws_ec2.TerminateInstanceHandler).Methods("POST")

	// GCP VM
	router.HandleFunc("/gcp/vm/createInstance", gcp_compute.CreateInstanceHandler).Methods("POST")
	router.HandleFunc("/gcp/vm/listInstances", gcp_compute.ListInstancesHandler).Methods("POST")
	router.HandleFunc("/gcp/vm/terminateInstance", gcp_compute.TerminateInstanceHandler).Methods("POST")

	// Azure VM
	router.HandleFunc("/azure/vm/createInstance", azure_vms.CreateVMHandler).Methods("POST")
	router.HandleFunc("/azure/vm/listInstances", azure_vms.ListVMsHandler).Methods("POST")
	router.HandleFunc("/azure/vm/terminateInstance", azure_vms.DeleteVMHandler).Methods("POST")

	// DynamoDB
	router.HandleFunc("/aws/dynamodb/createItem", aws_dynamodb.CreateItemHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/readItem", aws_dynamodb.ReadItemHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/deleteItem", aws_dynamodb.DeleteItemHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/updateItem", aws_dynamodb.UpdateItemHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/createTable", aws_dynamodb.CreateTableHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/deleteTable", aws_dynamodb.DeleteTableHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/updateTable", aws_dynamodb.UpdateTableHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/listTables", aws_dynamodb.ListTablesHandler).Methods("POST")
	router.HandleFunc("/aws/dynamodb/listItems", aws_dynamodb.ListItemsHandler).Methods("POST")

	// GCP Firebase
	router.HandleFunc("/gcp/firebase/createTable", gcp_firebase.CreateTableHandler).Methods("POST")
	router.HandleFunc("/gcp/firebase/deleteTable", gcp_firebase.DeleteTableHandler).Methods("POST")
	router.HandleFunc("/gcp/firebase/updateTable", gcp_firebase.UpdateTableHandler).Methods("PUT")
	router.HandleFunc("/gcp/firebase/listTables", gcp_firebase.ListTablesHandler).Methods("GET")

	// Azure Cosmos DB
	router.HandleFunc("/azure/cosmos/createAccount", azure_cosmosdb.CreateCosmosDBAccountHandler).Methods("POST")
	router.HandleFunc("/azure/cosmos/deleteAccount", azure_cosmosdb.DeleteCosmosDBAccountHandler).Methods("POST")
	router.HandleFunc("/azure/cosmos/createContainer", azure_cosmosdb.CreateCosmosDBContainerHandler).Methods("POST")
	router.HandleFunc("/azure/cosmos/deleteContainer", azure_cosmosdb.DeleteCosmosDBContainerHandler).Methods("POST")
	router.HandleFunc("/azure/cosmos/createdatabase", azure_cosmosdb.CreateCosmosDBDatabaseHandler).Methods("POST")
	router.HandleFunc("/azure/cosmos/deleteDatabase", azure_cosmosdb.DeleteCosmosDBDatabaseHandler).Methods("POST")
	router.HandleFunc("/azure/cosmos/listAccounts", azure_cosmosdb.ListCosmosDBAccountsHandler).Methods("GET")

	// GCP Network
	router.HandleFunc("/gcp/network/createNetwork", gcp_network.CreateNetworkHandler).Methods("POST")
	router.HandleFunc("/gcp/network/listNetworks", gcp_network.ListNetworksHandler).Methods("GET")
	router.HandleFunc("/gcp/network/deleteNetwork", gcp_network.DeleteNetworkHandler).Methods("POST")
	router.HandleFunc("/gcp/network/createSubnet", gcp_network.CreateSubnetHandler).Methods("POST")
	router.HandleFunc("/gcp/network/deleteSubnet", gcp_network.DeleteSubnetHandler).Methods("POST")
	router.HandleFunc("/gcp/network/listSubnets", gcp_network.ListSubnetsHandler).Methods("GET")
	router.HandleFunc("/gcp/network/createRoute", gcp_network.CreateRouteHandler).Methods("POST")
	router.HandleFunc("/gcp/network/deleteRoute", gcp_network.DeleteRouteHandler).Methods("POST")
	router.HandleFunc("/gcp/network/listRoutes", gcp_network.ListRoutesHandler).Methods("GET")
	router.HandleFunc("/gcp/firewall/createFirewallRule", gcp_network.CreateFirewallRuleHandler).Methods("POST")
	router.HandleFunc("/gcp/firewall/deleteFirewallRule", gcp_network.DeleteFirewallRuleHandler).Methods("POST")
	router.HandleFunc("/gcp/firewall/listFirewallRules", gcp_network.ListFirewallRulesHandler).Methods("GET")
	router.HandleFunc("/gcp/router/createCloudRouter", gcp_network.CreateCloudRouterHandler).Methods("POST")
	router.HandleFunc("/gcp/router/deleteCloudRouter", gcp_network.DeleteCloudRouterHandler).Methods("POST")
	router.HandleFunc("/gcp/router/listCloudRouters", gcp_network.ListCloudRoutersHandler).Methods("GET")

	// AWS Network
	router.HandleFunc("/aws/network/createVPC", aws_vpc.CreateVPCHandler).Methods("POST")
	router.HandleFunc("/aws/network/listVPCs", aws_vpc.ListVPCsHandler).Methods("POST")
	router.HandleFunc("/aws/network/deleteVPC", aws_vpc.DeleteVPCHandler).Methods("POST")
	router.HandleFunc("/aws/network/createSubnet", aws_vpc.CreateSubnetHandler).Methods("POST")
	router.HandleFunc("/aws/network/updateSubnet", aws_vpc.UpdateSubnetHandler).Methods("PUT")
	router.HandleFunc("/aws/network/deleteSubnet", aws_vpc.DeleteSubnetHandler).Methods("POST")
	router.HandleFunc("/aws/network/listSubnet", aws_vpc.ListSubnetsHandler).Methods("POST")
	router.HandleFunc("/aws/network/createRouteTable", aws_vpc.CreateRouteTableHandler).Methods("POST")
	router.HandleFunc("/aws/network/deleteRouteTable", aws_vpc.DeleteRouteTableHandler).Methods("POST")
	router.HandleFunc("/aws/network/listRouteTable", aws_vpc.ListRouteTableHandler).Methods("POST")
	router.HandleFunc("/aws/network/createInternetGateway", aws_vpc.CreateInternetGatewayHandler).Methods("POST")
	router.HandleFunc("/aws/network/attachInternetGateway", aws_vpc.AttachInternetGatewayHandler).Methods("POST")
	router.HandleFunc("/aws/network/detachInternetGateway", aws_vpc.DetachInternetGatewayHandler).Methods("POST")
	router.HandleFunc("/aws/network/deleteInternetGateway", aws_vpc.DeleteInternetGatewayHandler).Methods("POST")
	router.HandleFunc("/aws/network/listInternetGateway", aws_vpc.ListAllInternetGatewaysHandler).Methods("POST")

	// Azure Network
	router.HandleFunc("/azure/network/createVNet", azure_network.CreateNetworkHandler).Methods("POST")
	router.HandleFunc("/azure/network/listVNet", azure_network.ListNetworksHandler).Methods("GET")
	router.HandleFunc("/azure/network/deleteVNet", azure_network.DeleteNetworkHandler).Methods("POST")
	router.HandleFunc("/azure/network/createSubnet", azure_network.CreateSubnetHandler).Methods("POST")
	router.HandleFunc("/azure/network/deleteSubnet", azure_network.DeleteSubnetHandler).Methods("POST")
	router.HandleFunc("/azure/network/createFirewall", azure_network.CreateFirewallHandler).Methods("POST")
	router.HandleFunc("/azure/network/deleteFirewall", azure_network.DeleteFirewallHandler).Methods("POST")
	router.HandleFunc("/azure/network/listFirewalls", azure_network.ListFirewallHandler).Methods("GET")
	router.HandleFunc("/azure/network/listSubnets", azure_network.ListSubnetHandler).Methods("GET")

	// Serverless AWS ECS
	router.HandleFunc("/aws/ecs/createService", aws_ecs.CreateServiceHandler).Methods("POST")
	router.HandleFunc("/aws/ecs/deleteService", aws_ecs.DeleteServiceHandler).Methods("POST")
	router.HandleFunc("/aws/ecs/updateService", aws_ecs.UpdateServiceHandler).Methods("POST")
	router.HandleFunc("/aws/ecs/listServices", aws_ecs.ListServicesHandler).Methods("POST")
	router.HandleFunc("/aws/ecs/createCluster", aws_ecs.CreateECSHandler).Methods("POST")
	router.HandleFunc("/aws/ecs/deleteCluster", aws_ecs.DeleteECSHandler).Methods("POST")
	router.HandleFunc("/aws/ecs/listClusters", aws_ecs.ListClustersHandler).Methods("POST")

	// Serverless AWS EKS
	router.HandleFunc("/aws/eks/updateCluster", aws_eks.UpdateEKSHandler).Methods("PUT")
	router.HandleFunc("/aws/eks/createCluster", aws_eks.CreateEKSHandler).Methods("POST")
	router.HandleFunc("/aws/eks/deleteCluster", aws_eks.DeleteEKSHandler).Methods("POST")
	router.HandleFunc("/aws/eks/listClusters", aws_eks.ListEKSHandler).Methods("POST")

	// Lambda
	router.HandleFunc("/aws/lambda/createFunction", aws_lambda.CreateLambdaFunctionHandler).Methods("POST")
	router.HandleFunc("/aws/lambda/listFunctions", aws_lambda.ListLambdaFunctionsHandler).Methods("GET")
	router.HandleFunc("/aws/lambda/deleteFunction", aws_lambda.DeleteLambdaFunctionHandler).Methods("POST")

	// Serverless GCP Cloud Functions
	router.HandleFunc("/gcp/cloudfunctions/createcloudRunService", gcp_cloudrun.CreateCloudRunServiceHandler).Methods("POST")
	router.HandleFunc("/gcp/cloudfunctions/deletecloudRunService", gcp_cloudrun.DeleteCloudRunServiceHandler).Methods("POST")
	router.HandleFunc("/gcp/cloudfunctions/listcloudRunService", gcp_cloudrun.ListCloudRunServicesHandler).Methods("GET")
	router.HandleFunc("/gcp/cloudfunctions/createCloudFunction", gcp_cloudrun.CreateCloudFunctionHandler).Methods("POST")
	router.HandleFunc("/gcp/cloudfunctions/deleteCloudFunction", gcp_cloudrun.DeleteCloudFunctionHandler).Methods("POST")
	router.HandleFunc("/gcp/cloudfunctions/listCloudFunction", gcp_cloudrun.ListCloudFunctionsHandler).Methods("GET")

	// Serverless GCP GKE
	router.HandleFunc("/gcp/gke/createCluster", gcp_kubernetes.CreateClusterHandler).Methods("POST")
	router.HandleFunc("/gcp/gke/deleteCluster", gcp_kubernetes.DeleteClusterHandler).Methods("POST")
	router.HandleFunc("/gcp/gke/listClusters", gcp_kubernetes.ListClustersHandler).Methods("GET")

	// Azure Functions
	router.HandleFunc("/azure/functions/createFunction", azure_functions.CreateFunctionAppHandler).Methods("POST")
	router.HandleFunc("/azure/functions/deleteFunction", azure_functions.DeleteFunctionAppHandler).Methods("POST")
	router.HandleFunc("/azure/functions/listFunctions", azure_functions.ListFunctionAppsHandler).Methods("GET")

	// Setup CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Adjust this to the appropriate domains in production
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization"}),
	)(router)

	// Apply the CORS middleware to the router
	// corsRouter := handlers.CORS(cors)(router)

	// Start the server
	log.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", cors)
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "570397565376-ks7fmsvgrma2c9gm2k8lfa5tjhqdpala.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{
			"https://www.googleapis.com/auth/datastore",
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/devstorage.read_write",
		},

		Endpoint: googleAuth.Endpoint,
	}
)

type googleToken struct {
	AccountID int `json:"accountID"`
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	//googleToken
	var req googleToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}
	fmt.Println("req done", req)
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}
	googleOauthConfig.ClientSecret = cloudAccount.ClientSecret.String
	authURL := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	authCode := r.URL.Query().Get("code")
	fmt.Println("googleOauthConfig:", googleOauthConfig.ClientSecret)
	token, err := googleOauthConfig.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
		return
	}
	fmt.Println("token: ", token.AccessToken)

	http.SetCookie(w, &http.Cookie{Name: "accessTokenGCP", Value: token.AccessToken})

	w.Header().Set("Authorization", "Bearer "+token.AccessToken)
	http.Redirect(w, r, "http://localhost:3000/cloud", http.StatusTemporaryRedirect)

}

var (
	clientID     = "YOUR_CLIENT_ID_HERE"
	clientSecret = "YOUR_CLIENT"
	redirectURI  = "http://localhost:8080/auth/azure/callback"
	authURL      = "https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/authorize"
	tokenURL     = "https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/token"
	scope        = "https://management.azure.com/user_impersonation"
	cookieName   = "codeVerifier"
)

func handleAzureLogin(w http.ResponseWriter, r *http.Request) {
	var req googleToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}
	fmt.Println("req done", req)
	accountID := req.AccountID

	state := "abcdefghijklmnop"
	cloudAccount, err := db.GetCloudAccountDetails(accountID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert accountIDStr to int: %v", err), http.StatusInternalServerError)
		return
	}
	clientID = cloudAccount.ClientID.String
	clientSecret = cloudAccount.ClientSecret.String

	params := url.Values{}
	params.Set("client_id", cloudAccount.ClientID.String)
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", scope)
	params.Set("code_challenge", "IdbruyXUD2XUVuNsvyuss8gxfag0PMLVjjoOzyZCaIY")
	params.Set("code_challenge_method", "S256")
	params.Set("state", state)

	redirectURL := authURL + "?" + params.Encode()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func handleAzureCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tokenParams := url.Values{}
	tokenParams.Set("client_id", clientID)
	tokenParams.Set("client_secret", clientSecret)
	tokenParams.Set("code", code)
	tokenParams.Set("redirect_uri", redirectURI)
	tokenParams.Set("grant_type", "authorization_code")
	tokenParams.Set("code_verifier", "u7asE4uSEPKu3eRl18_aJylYVXoTeaZWse1U-2mky5FDl881RtPHAagcKEBIIlcx7919ZcR5LHhD4FjP0ocVZ7rk3LJFmO_V42-evDi7A_ZbbaCvWNpv9MLjrTXFJZPv")
	tokenParams.Set("scope", scope)
	reqBody := strings.NewReader(tokenParams.Encode())

	req, err := http.NewRequest("POST", tokenURL, reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create POST request: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make POST request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode response body: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("respBody:", respBody)
	accessToken1, ok := respBody["access_token"].(string)
	if !ok {
		http.Error(w, "Access token not found in response", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "accessTokenAZURE", Value: string(accessToken1)})
	fmt.Println("token:", accessToken1)
	// w.Write(respBody)
	http.Redirect(w, r, "http://localhost:3000/cloud", http.StatusTemporaryRedirect)
}

// func generateCodeChallenge(verifier string) string {
// 	h := sha256.Sum256([]byte(verifier))
// 	return base64.RawURLEncoding.EncodeToString(h[:])
// }

// func generateRandomString(length int) string {
// 	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	b := make([]byte, length)
// 	for i := range b {
// 		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
// 		b[i] = chars[idx.Int64()]
// 	}
// 	return string(b)
// }
