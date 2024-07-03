package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// CloudAccount represents a row in the CloudAccount table
type CloudAccount struct {
	AccountID      int
	UserID         int
	CloudProvider  string
	AccessKey      *string
	SecretKey      *string
	SubscriptionID *string
	TenantID       *string
	ClientID       *string
	ClientSecret   *string
	Region         *string
	AdditionalInfo *string
	ClientEmail    *string
	PrivateKey     *string
	ProjectID      *string
}

// GetCloudAccountDetails retrieves ClientEmail, PrivateKey, and ProjectID for a given AccountID
func GetCloudAccountDetails(accountID int) ([]string, error) {
	// Connect to the MySQL database
	db, err := sql.Open("mysql", "newuser:password@tcp(127.0.0.1:3307)/multicloud")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var clientEmail, privateKey, projectID, region, additionalInfo, cloudProvider, accessKey, secretKey, subscriptionID, tenantID, clientID, clientSecret string
	row := db.QueryRow("SELECT ClientEmail, PrivateKey, ProjectID FROM CloudAccount WHERE AccountID = ?", accountID)
	err = row.Scan(&clientEmail, &privateKey, &projectID, &region, &additionalInfo, &cloudProvider, &accessKey, &secretKey, &subscriptionID, &tenantID, &clientID, &clientSecret)
	if err != nil {
		return nil, err
	}

	return []string{clientEmail, privateKey, projectID, region, additionalInfo, cloudProvider, accessKey, secretKey, subscriptionID, tenantID, clientID, clientSecret}, nil
}
