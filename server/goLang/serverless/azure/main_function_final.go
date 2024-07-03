package azure_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

type FunctionAppRequest struct {
	SubscriptionID        string            `json:"subscriptionID"`
	ResourceGroup         string            `json:"resourceGroup"`
	AppServicePlan        string            `json:"appServicePlan"`
	Location              string            `json:"location"`
	FunctionAppName       string            `json:"functionAppName"`
	Environment           map[string]string `json:"environment"`
	ClientAffinityEnabled bool              `json:"clientAffinityEnabled"`
	Token                 string            `json:"token"`
}

type FunctionAppResponse struct {
	Message      string   `json:"message"`
	FunctionApps []string `json:"functionApps,omitempty"`
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

func initFunctionAppClient(subscriptionID string, token string) (web.AppsClient, error) {
	client := web.NewAppsClient(subscriptionID)
	client.Authorizer = autorest.NullAuthorizer{}
	client.RequestInspector = tokenAuthorizer{token: token}.WithAuthorization()
	return client, nil
}

func CreateFunctionAppHandler(w http.ResponseWriter, r *http.Request) {
	var req FunctionAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initFunctionAppClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize client: %v", err), http.StatusInternalServerError)
		return
	}

	siteProperties := &web.SiteProperties{
		ServerFarmID:          to.StringPtr(req.AppServicePlan),
		ClientAffinityEnabled: to.BoolPtr(req.ClientAffinityEnabled),
	}

	// Setting environment variables if provided
	if len(req.Environment) > 0 {
		appSettings := make([]web.NameValuePair, 0, len(req.Environment))
		for key, value := range req.Environment {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  to.StringPtr(key),
				Value: to.StringPtr(value),
			})
		}
		siteProperties.SiteConfig = &web.SiteConfig{
			AppSettings: &appSettings,
		}
	}

	params := web.Site{
		Location:       to.StringPtr(req.Location),
		SiteProperties: siteProperties,
	}

	_, err = client.CreateOrUpdate(context.Background(), req.ResourceGroup, req.FunctionAppName, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Function App: %v", err), http.StatusInternalServerError)
		return
	}

	resp := FunctionAppResponse{Message: "Function app created successfully"}
	json.NewEncoder(w).Encode(resp)
}

func ListFunctionAppsHandler(w http.ResponseWriter, r *http.Request) {
	var req FunctionAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initFunctionAppClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	includeSlots := false
	result, err := client.ListByResourceGroup(ctx, req.ResourceGroup, &includeSlots)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing Function Apps: %v", err), http.StatusInternalServerError)
		return
	}

	var functionApps []string
	for _, app := range result.Values() {
		functionApps = append(functionApps, *app.Name)
	}

	resp := FunctionAppResponse{FunctionApps: functionApps}
	json.NewEncoder(w).Encode(resp)
}

func DeleteFunctionAppHandler(w http.ResponseWriter, r *http.Request) {
	var req FunctionAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := initFunctionAppClient(req.SubscriptionID, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deleteMetrics := true
	deleteEmptyServerFarm := true
	_, err = client.Delete(context.Background(), req.ResourceGroup, req.FunctionAppName, &deleteMetrics, &deleteEmptyServerFarm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting Function App: %v", err), http.StatusInternalServerError)
		return
	}

	resp := FunctionAppResponse{Message: "Function app deleted successfully"}
	json.NewEncoder(w).Encode(resp)
}
