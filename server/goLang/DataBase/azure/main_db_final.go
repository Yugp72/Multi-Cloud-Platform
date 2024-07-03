package azure_cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cosmos-db/mgmt/documentdb"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

type CosmosDBAccountRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	Token          string `json:"token"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
	AccountName    string `json:"accountName"`
	DatabaseName   string `json:"databaseName"`
	CollectionName string `json:"collectionName"`
}

type CosmosDBAccountResponse struct {
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

func initCosmosDBClient(subscriptionID string, token string) (documentdb.DatabaseAccountsClient, error) {
	client := documentdb.NewDatabaseAccountsClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{} // Set to NullAuthorizer since we are manually inserting the token
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}

func initSQLClient(subscriptionID string, token string) (documentdb.SQLResourcesClient, error) {
	client := documentdb.NewSQLResourcesClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{} // Set to NullAuthorizer since we are manually inserting the token
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}

// DeleteCosmosDBAccountHandler deletes a Cosmos DB account
func DeleteCosmosDBAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req CosmosDBAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	client, err := initCosmosDBClient(req.SubscriptionID, req.Token)

	// Delete the Cosmos DB account
	_, err = client.Delete(context.Background(), req.ResourceGroup, req.AccountName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting Cosmos DB account: %v", err)
		return
	}

	// Send success response
	resp := CosmosDBAccountResponse{Message: "Cosmos DB account deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListCosmosDBAccountsHandler lists Cosmos DB accounts in a resource group
func ListCosmosDBAccountsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req CosmosDBAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	client, err := initCosmosDBClient(req.SubscriptionID, req.Token)

	// List Cosmos DB accounts
	accounts, err := client.ListByResourceGroup(context.Background(), req.ResourceGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing Cosmos DB accounts: %v", err)
		return
	}

	// Convert accounts to JSON and send as response
	jsonResponse, err := json.Marshal(accounts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func CreateCosmosDBAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req CosmosDBAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	client, err := initCosmosDBClient(req.SubscriptionID, req.Token)

	accountParameters := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: &req.Location,
		Kind:     documentdb.DatabaseAccountKindGlobalDocumentDB, // Adjust accordingly based on the API kind
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType: to.StringPtr("Standard"), // or another offer type as per your requirement
			Locations: &[]documentdb.Location{
				{
					FailoverPriority: to.Int32Ptr(0),
					LocationName:     to.StringPtr(req.Location),
				},
			},
		},
	}

	_, err = client.CreateOrUpdate(context.Background(), req.ResourceGroup, req.AccountName, accountParameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Cosmos DB account: %v", err)
		return
	}

	resp := CosmosDBAccountResponse{Message: "Cosmos DB account created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleError is a utility function to handle errors
func handleError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, message)
}

type CosmosDBDatabaseRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	Token          string `json:"token"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
	AccountName    string `json:"accountName"`
	DatabaseName   string `json:"databaseName"`
	DatabaseID     string `json:"databaseID"`
}

type CosmosDBDatabaseResponse struct {
	Message string `json:"message"`
}

func CreateCosmosDBDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	var req CosmosDBDatabaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := initSQLClient(req.SubscriptionID, req.Token)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to initialize Cosmos DB client: %v", err))
		return
	}

	err = createDatabase(context.Background(), &client, req.AccountName, req.ResourceGroup, req.DatabaseName, req.Location)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating Cosmos DB database: %v", err))
		return
	}

	resp := CosmosDBDatabaseResponse{Message: "Cosmos DB database created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func createDatabase(ctx context.Context, databaseClient *documentdb.SQLResourcesClient, accountName, resourceGroup, dbName, location string) error {
	databaseProperties := documentdb.SQLDatabaseResource{
		ID: to.StringPtr(dbName),
	}
	parameters := documentdb.SQLDatabaseCreateUpdateParameters{
		SQLDatabaseCreateUpdateProperties: &documentdb.SQLDatabaseCreateUpdateProperties{
			Resource: &databaseProperties,
			Options: &documentdb.CreateUpdateOptions{
				Throughput: to.Int32Ptr(400),
			},
		},
	}
	_, err := databaseClient.CreateUpdateSQLDatabase(ctx, resourceGroup, accountName, dbName, parameters)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCosmosDBDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	var req CosmosDBDatabaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := initSQLClient(req.SubscriptionID, req.Token)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to initialize Cosmos DB client: %v", err))
		return
	}

	err = deleteDatabase(context.Background(), &client, req.AccountName, req.ResourceGroup, req.DatabaseName)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Cosmos DB database: %v", err))
		return
	}

	resp := CosmosDBDatabaseResponse{Message: "Cosmos DB database deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func deleteDatabase(ctx context.Context, databaseClient *documentdb.SQLResourcesClient, accountName, resourceGroup, dbName string) error {
	_, err := databaseClient.DeleteSQLDatabase(ctx, resourceGroup, accountName, dbName)
	if err != nil {
		return err
	}
	return nil
}

type CosmosDBContainerRequest struct {
	SubscriptionID string `json:"subscriptionID"`
	Token          string `json:"token"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
	AccountName    string `json:"accountName"`
	DatabaseName   string `json:"databaseName"`
	ContainerName  string `json:"containerName"`
}

type CosmosDBContainerResponse struct {
	Message string `json:"message"`
}

func CreateCosmosDBContainerHandler(w http.ResponseWriter, r *http.Request) {
	var req CosmosDBContainerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := initSQLClient(req.SubscriptionID, req.Token)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to initialize Cosmos DB client: %v", err))
		return
	}

	err = createContainer(context.Background(), &client, req.AccountName, req.ResourceGroup, req.DatabaseName, req.ContainerName)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating Cosmos DB container: %v", err))
		return
	}

	resp := CosmosDBContainerResponse{Message: "Cosmos DB container created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func createContainer(ctx context.Context, containerClient *documentdb.SQLResourcesClient, accountName, resourceGroup, dbName, containerName string) error {
	containerProperties := documentdb.SQLContainerResource{
		ID: to.StringPtr(containerName),
		PartitionKey: &documentdb.ContainerPartitionKey{
			Paths: &[]string{"/partitionKey"}, // Specify the partition key path
			Kind:  documentdb.PartitionKindHash,
		},
	}
	parameters := documentdb.SQLContainerCreateUpdateParameters{
		SQLContainerCreateUpdateProperties: &documentdb.SQLContainerCreateUpdateProperties{
			Resource: &containerProperties,
		},
	}

	_, err := containerClient.CreateUpdateSQLContainer(ctx, resourceGroup, accountName, dbName, containerName, parameters)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCosmosDBContainerHandler(w http.ResponseWriter, r *http.Request) {
	var req CosmosDBContainerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := initSQLClient(req.SubscriptionID, req.Token)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to initialize Cosmos DB client: %v", err))
		return
	}

	err = deleteContainer(context.Background(), &client, req.AccountName, req.ResourceGroup, req.DatabaseName, req.ContainerName)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Cosmos DB container: %v", err))
		return
	}

	resp := CosmosDBContainerResponse{Message: "Cosmos DB container deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func deleteContainer(ctx context.Context, containerClient *documentdb.SQLResourcesClient, accountName, resourceGroup, dbName, containerName string) error {
	_, err := containerClient.DeleteSQLContainer(ctx, resourceGroup, accountName, dbName, containerName)
	if err != nil {
		return err
	}
	return nil
}
