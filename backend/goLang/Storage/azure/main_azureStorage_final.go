package azure_storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	db "btep.project/databaseConnection"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

type StorageAccountRequest struct {
	AccountName    string `json:"accountName"`
	ResourceGroup  string `json:"resourceGroup"`
	SubscriptionID string `json:"subscriptionID"`
	Token          string `json:"token"`
	Location       string `json:"location"`
	StorageType    string `json:"storageType"`
	AccessTier     string `json:"accessTier"`
}

type StorageAccountResponse struct {
	Message string `json:"message"`
}

func CreateStorageAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req StorageAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Create Azure Storage account client
	client, err := initStorageClient(req.SubscriptionID, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure Storage account client: %v", err)
		return
	}

	// Create storage account parameters
	accountParameters := storage.AccountCreateParameters{
		Location: &req.Location,
		Sku: &storage.Sku{
			Name: storage.SkuName(req.StorageType),
		},

		Kind: storage.StorageV2,
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			AccessTier: storage.AccessTier(req.AccessTier),
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: to.BoolPtr(true),
					},
				},
				KeySource: storage.KeySource(storage.MicrosoftStorage),
			},
		},

		Tags: map[string]*string{},
	}

	// Create the storage account
	_, err = client.Create(context.Background(), req.ResourceGroup, req.AccountName, accountParameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Azure Storage account: %v", err)
		return
	}

	// Send success response
	resp := StorageAccountResponse{Message: "Azure Storage account created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteStorageAccountHandler deletes an Azure Storage account
func DeleteStorageAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req StorageAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Create Azure Storage account client
	client, err := initStorageClient(req.SubscriptionID, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure Storage account client: %v", err)
		return
	}

	// Delete the storage account
	_, err = client.Delete(context.Background(), req.ResourceGroup, req.AccountName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting Azure Storage account: %v", err)
		return
	}

	// Send success response
	resp := StorageAccountResponse{Message: "Azure Storage account deleted successfully"}
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

func initStorageClient(subscriptionID string, token string) (storage.AccountsClient, error) {
	client := storage.NewAccountsClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{}
	client.RequestInspector = withBearerToken(token)
	return client, nil
}

type ListStorageAccountRequest struct {
	AccountID int    `json:"accountID"`
	Token     string `json:"token"`
}

func ListStorageAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListStorageAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	client, err := initStorageClient(cloudAccount.SubscriptionID.String, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure Storage account client: %v", err)
		return
	}
	fmt.Println("client:", client)

	accounts, err := client.ListComplete(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing Azure Storage accounts: %v", err)
		return
	}
	fmt.Println("accounts:", accounts)

	// Convert accounts to JSON and send as response
	jsonResponse, err := json.Marshal(accounts)
	fmt.Println("jsonResponse:", jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func UpdateStorageAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req StorageAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	client, err := initStorageClient(req.SubscriptionID, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure Storage account client: %v", err)
		return
	}

	// Define update parameters
	updateParams := storage.AccountUpdateParameters{
		Sku: &storage.Sku{
			Name: storage.SkuName(req.StorageType),
		},
		Tags: map[string]*string{},
		// Add other properties to update as needed
	}

	_, err = client.Update(context.Background(), req.ResourceGroup, req.AccountName, updateParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating Azure Storage account: %v", err)
		return
	}

	// Send success response
	resp := StorageAccountResponse{Message: "Azure Storage account updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type UploadObjectRequest struct {
	AccountName    string                `json:"accountName"`
	ResourceGroup  string                `json:"resourceGroup"`
	SubscriptionID string                `json:"subscriptionID"`
	Token          string                `json:"token"`
	ContainerName  string                `json:"containerName"`
	ObjectName     string                `json:"objectName"`
	BlobName       string                `json:"blobName"`
	file           *multipart.FileHeader `json:"file"`
}

type UploadObjectResponse struct {
	Message string `json:"message"`
}

func UploadBlobHandler(w http.ResponseWriter, r *http.Request) {
	var req UploadObjectRequest
	req.AccountName = r.FormValue("accountName")
	req.Token = r.FormValue("token")
	req.ContainerName = r.FormValue("containerName")
	req.ObjectName = r.FormValue("objectName")
	req.BlobName = r.FormValue("blobName")
	req.SubscriptionID = r.FormValue("subscriptionID")
	req.ResourceGroup = r.FormValue("resourceGroup")
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "Invalid request body", err)
	// 	return
	// }

	client, err := initStorageClient(req.SubscriptionID, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to initialize Azure Storage account client: %v", err)
		return
	}

	containerURL := getContainerURL(client, req.ResourceGroup, req.AccountName)
	file, fileHeader, err := r.FormFile("file")
	blobURL := containerURL.NewBlockBlobURL(fileHeader.Filename)

	fmt.Println(file)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to open file: %v", err)
		return
	}
	defer file.Close()

	fileBytes := new(bytes.Buffer)
	_, err = fileBytes.ReadFrom(file)
	fileBytes1 := fileBytes.Bytes()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to read file: %v", err)
		return
	}

	_, err = azblob.UploadBufferToBlockBlob(context.Background(), fileBytes1, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to upload blob: %v", err)
		return
	}

	resp := StorageAccountResponse{Message: "Blob uploaded successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getContainerURL(client storage.AccountsClient, resourceGroup, accountName string) azblob.ContainerURL {
	keys, err := client.ListKeys(context.Background(), resourceGroup, accountName, "")
	if err != nil {
		fmt.Println("Failed to get storage account keys:", err)
		os.Exit(1)
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, *(*keys.Keys)[0].Value)
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	url, err := url.Parse("https://" + accountName + ".blob.core.windows.net")
	serviceURL := azblob.NewServiceURL(*url, p)
	containerURL := serviceURL.NewContainerURL("vatsaldp")

	return containerURL
}

type GetObjectRequest struct {
	AccountName    string `json:"accountName"`
	Token          string `json:"token"`
	ContainerName  string `json:"containerName"`
	ObjectName     string `json:"objectName"`
	SubscriptionID string `json:"subscriptionID"`
	ResourceGroup  string `json:"resourceGroup"`
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	var req GetObjectRequest
	req.AccountName = r.FormValue("accountName")
	req.Token = r.FormValue("token")
	req.ContainerName = r.FormValue("containerName")
	req.ObjectName = r.FormValue("objectName")
	req.SubscriptionID = r.FormValue("subscriptionID")
	req.ResourceGroup = r.FormValue("resourceGroup")

	// Initialize Azure Blob Storage client
	client, err := initStorageClient(req.SubscriptionID, req.Token)
	if err != nil {
		handleError(w, fmt.Sprintf("Failed to initialize Azure Storage account client: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a container URL
	containerURL := getContainerURL(client, req.ResourceGroup, req.AccountName)

	// Create a blob URL within the container
	blobURL := containerURL.NewBlockBlobURL(req.ObjectName)

	// Download the blob content
	downloadResponse, err := blobURL.Download(context.Background(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		handleError(w, fmt.Sprintf("Failed to download object from Azure Storage: %v", err), http.StatusInternalServerError)
		return
	}
	defer downloadResponse.Body(azblob.RetryReaderOptions{}).Close()

	// Set response headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+req.ObjectName)

	// Write the blob content to the response writer
	if _, err := io.Copy(w, downloadResponse.Body(azblob.RetryReaderOptions{})); err != nil {
		handleError(w, fmt.Sprintf("Failed to write object data to response: %v", err), http.StatusInternalServerError)
		return
	}
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, errMsg)
}
