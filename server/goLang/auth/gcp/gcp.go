package connection

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	googleAuth "golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8081/auth/google/callback",
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

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	token, err := googleOauthConfig.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
		return
	}

	if err != nil {
		log.Fatalf("Cannot create Firestore client: %v", err)
		return
	}
	fmt.Printf("Token: %v", token.AccessToken)
	http.Redirect(w, r, "http://localhost:3000/cloud", http.StatusTemporaryRedirect)

	// w.WriteHeader(http.StatusPermanentRedirect)
	// w.Header().Set("Location", "http://localhost:3000/cloud")
}
