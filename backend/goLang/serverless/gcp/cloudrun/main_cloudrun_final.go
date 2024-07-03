package gcp_cloudrun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "btep.project/databaseConnection"
	"golang.org/x/oauth2"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

type CloudRunServiceRequest struct {
	AccountID int    `json:"accountId"`
	Token     string `json:"token"`
	ProjectID string `json:"projectId"`
	Location  string `json:"location"`
	Service   string `json:"service"`
}

type CloudRunServiceResponse struct {
	Message string `json:"message"`
}

func initRunService(token string) (*run.APIService, error) {
	ctx := context.Background()
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	service, err := run.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating Cloud Run service: %v", err)
	}
	service.BasePath = "https://run.googleapis.com/"

	return service, nil
}

func initFunctionsService(token string) (*cloudfunctions.Service, error) {
	ctx := context.Background()
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	Service, err := cloudfunctions.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating compute service: %v", err)
	}

	return Service, nil
}

func CreateCloudRunServiceHandler(w http.ResponseWriter, r *http.Request) {
	var req CloudRunServiceRequest
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

	// Initialize Cloud Run service
	runService, err := initRunService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing Cloud Run service: %v", err)
		return
	}

	service := &run.Service{
		ApiVersion: "serving.knative.dev/v1",
		Kind:       "Service",
		Metadata: &run.ObjectMeta{
			Name: "my-cloud-run-service",
		},
		Spec: &run.ServiceSpec{
			Template: &run.RevisionTemplate{
				Metadata: &run.ObjectMeta{
					Annotations: map[string]string{
						"autoscaling.knative.dev/maxScale": "5",
					},
				},
				Spec: &run.RevisionSpec{
					ContainerConcurrency: 80,
					Containers: []*run.Container{
						{
							Image: "gcr.io/my-project/my-image:latest",
							Ports: []*run.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create the Cloud Run service
	createCall := runService.Projects.Locations.Services.Create(cloudAccount.ProjectID.String, service)
	_, err = createCall.Do()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Cloud Run service: %v", err)
		return
	}

	// Send success response
	resp := CloudRunServiceResponse{Message: "Cloud Run service created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ListCloudRunServicesHandler(w http.ResponseWriter, r *http.Request) {

}

type DeleteCloudRunServiceRequest struct {
	AccountID int    `json:"accountId"`
	Token     string `json:"token"`
	Location  string `json:"location"`
	Service   string `json:"service"`
}

func DeleteCloudRunServiceHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteCloudRunServiceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Fetch cloud account details from the database
	cloudAccount, err := db.GetCloudAccountDetails(req.AccountID)
	fmt.Println(cloudAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting cloud account details: %v", err)
		return
	}

	// Initialize Cloud Run service
	runService, err := initRunService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing Cloud Run service: %v", err)
		return
	}

	// Delete the Cloud Run service
	_, err = runService.Projects.Locations.Services.Delete(req.Service).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting Cloud Run service: %v", err)
		return
	}

	// Send success response
	resp := CloudRunServiceResponse{Message: "Cloud Run service deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type CloutFuntionRequest struct {
	AccountID int    `json:"accountId"`
	Token     string `json:"token"`
	Location  string `json:"location"`
	Service   string `json:"service"`
}

func CreateCloudFunctionHandler(w http.ResponseWriter, r *http.Request) {
	var req CloutFuntionRequest
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

	// Initialize Cloud Functions service
	functionsService, err := initFunctionsService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing Cloud Functions service: %v", err)
		return
	}

	// Create the Cloud Function object using request parameters
	function := &cloudfunctions.CloudFunction{
		Name:             "projects/" + cloudAccount.ProjectID.String + "/locations/" + req.Location + "/functions/" + req.Service,
		Description:      req.Service + " function",
		EntryPoint:       req.Service + "EntryPoint",
		Runtime:          "nodejs14",
		SourceArchiveUrl: req.Service + ".zip",
	}

	// Create the Cloud Function
	createCall := functionsService.Projects.Locations.Functions.Create(cloudAccount.ProjectID.String, function)
	_, err = createCall.Do()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating Cloud Function: %v", err)
		return
	}

	// Send success response
	resp := CloudRunServiceResponse{Message: "Cloud Function created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ListCloudRunFunctionRequest struct {
	AccountID int    `json:"accountId"`
	Token     string `json:"token"`
	Location  string `json:"location"`
}

func ListCloudFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	var req ListCloudRunFunctionRequest
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

	// Initialize Cloud Functions service
	functionsService, err := initFunctionsService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing Cloud Functions service: %v", err)
		return
	}

	// List Cloud Functions
	listCall := functionsService.Projects.Locations.Functions.List("projects/" + cloudAccount.ProjectID.String + "/locations/" + req.Location)
	response, err := listCall.Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing Cloud Functions: %v", err)
		return
	}

	// Extract function names from the response
	var functions []string
	for _, function := range response.Functions {
		functions = append(functions, function.Name)
	}

	// Send success response
	resp := struct {
		Functions []string `json:"functions"`
	}{
		Functions: functions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type GetCloudRunFunctionRequest struct {
	AccountID int    `json:"accountId"`
	Token     string `json:"token"`
	Location  string `json:"location"`
	Service   string `json:"service"`
}

func GetCloudFunctionHandler(w http.ResponseWriter, r *http.Request) {
	var req GetCloudRunFunctionRequest
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

	// Initialize Cloud Functions service
	functionsService, err := initFunctionsService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing Cloud Functions service: %v", err)
		return
	}

	// Get Cloud Function details
	getCall := functionsService.Projects.Locations.Functions.Get("projects/" + cloudAccount.ProjectID.String + "/locations/" + req.Location + "/functions/" + req.Service)
	function, err := getCall.Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching Cloud Function: %v", err)
		return
	}

	// Send success response with Cloud Function details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(function)
}

type DeleteCloudFuntion struct {
	AccountID int    `json:"accountId"`
	Token     string `json:"token"`
	Location  string `json:"location"`
	Service   string `json:"service"`
}

func DeleteCloudFunctionHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteCloudFuntion
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

	// Initialize Cloud Functions service
	functionsService, err := initFunctionsService(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing Cloud Functions service: %v", err)
		return
	}

	// Delete the Cloud Function
	deleteCall := functionsService.Projects.Locations.Functions.Delete("projects/" + cloudAccount.ProjectID.String + "/locations/" + req.Location + "/functions/" + req.Service)
	_, err = deleteCall.Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting Cloud Function: %v", err)
		return
	}

	// Send success response
	resp := CloudRunServiceResponse{Message: "Cloud Function deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
