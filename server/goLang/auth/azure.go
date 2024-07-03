package auth_gcp1

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"

	db "btep.project/databaseConnection"
	"github.com/gorilla/mux"
)

var (
	clientID    = "YOUR_CLIENT_ID_HERE"
	redirectURI = "http://localhost:8081/auth/azure/callback"
	authURL     = "https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-/oauth2/v2.0/authorize"
	tokenURL    = "https://login.microsoftonline.com/8059c3d5-a962--ef7c9211422a/oauth2/v2.0/token"
	scope       = "https://management.azure.com/.default"
	cookieName  = "codeVerifier"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth/azure/login", handleAzureLogin).Methods("GET")
	r.HandleFunc("/auth/azure/callback", handleAzureCallback).Methods("GET")
	http.Handle("/", r)

	fmt.Println("Server listening on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func generateCodeVerifier(accountID int) string {
	// Get the cloud account details using accountID
	cloudAccount, err := db.GetCloudAccountDetails(accountID)
	if err != nil {
		log.Fatalf("Error getting cloud account details: %v", err)
	}

	// Concatenate account details to create a unique verifier
	verifierData := fmt.Sprintf("%s-%s-%s", cloudAccount.ClientID, cloudAccount.ClientSecret)

	// Hash the verifier data
	h := sha256.Sum256([]byte(verifierData))

	// Encode the hash as base64
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func handleAzureLogin(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("accountID")
	if accountIDStr == "" {
		http.Error(w, "Account ID not provided", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	codeVerifier := generateCodeVerifier(accountID)

	state := generateRandomString(16)

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", scope)
	params.Set("code_challenge", generateCodeChallenge(codeVerifier))
	params.Set("code_challenge_method", "S256")
	params.Set("state", state)

	redirectURL := authURL + "?" + params.Encode()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func handleAzureCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	codeVerifier := r.URL.Query().Get("codeVerifier")

	tokenParams := url.Values{}
	tokenParams.Set("client_id", clientID)
	tokenParams.Set("code", code)
	tokenParams.Set("redirect_uri", redirectURI)
	tokenParams.Set("grant_type", "authorization_code")
	tokenParams.Set("code_verifier", codeVerifier)
	tokenParams.Set("scope", scope)

	// Handle token exchange and response as needed...
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[idx.Int64()]
	}
	return string(b)
}
