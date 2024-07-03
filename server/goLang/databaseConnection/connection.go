package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// CloudAccount represents a row in the CloudAccount table
type CloudAccount struct {
	AccountID      int
	UserID         int
	CloudProvider  string
	AccessKey      sql.NullString
	SecretKey      sql.NullString
	SubscriptionID sql.NullString
	TenantID       sql.NullString
	ClientID       sql.NullString
	ClientSecret   sql.NullString
	Region         sql.NullString
	AdditionalInfo sql.NullString
	ClientEmail    sql.NullString
	PrivateKey     sql.NullString
	ProjectID      sql.NullString
}

// GetCloudAccountDetails retrieves cloud account details from the database
func GetCloudAccountDetails(accountID int) (*CloudAccount, error) {
	db, err := sql.Open("mysql", "newuser:password@tcp(127.0.0.1:3307)/multicloud")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var account CloudAccount
	row := db.QueryRow("SELECT ClientEmail, PrivateKey, ProjectID, Region, AdditionalInformation, CloudProvider, AccessKey, SecretKey, SubscriptionID, TenantID, ClientID, ClientSecret FROM CloudAccount WHERE AccountID = ?", accountID)
	err = row.Scan(
		&account.ClientEmail,
		&account.PrivateKey,
		&account.ProjectID,
		&account.Region,
		&account.AdditionalInfo,
		&account.CloudProvider,
		&account.AccessKey,
		&account.SecretKey,
		&account.SubscriptionID,
		&account.TenantID,
		&account.ClientID,
		&account.ClientSecret,
	)
	fmt.Println(account.PrivateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows found for account ID %d", accountID)
		}
		return nil, err
	}

	// Handle NULL values and set them to empty strings
	if !account.ClientEmail.Valid {
		account.ClientEmail.String = ""
	}
	if !account.PrivateKey.Valid {
		account.PrivateKey.String = ""
	}
	if !account.ProjectID.Valid {
		account.ProjectID.String = ""
	}
	if !account.Region.Valid {
		account.Region.String = ""
	}
	if !account.AdditionalInfo.Valid {
		account.AdditionalInfo.String = ""
	}
	if !account.AccessKey.Valid {
		account.AccessKey.String = ""
	}
	if !account.SecretKey.Valid {
		account.SecretKey.String = ""
	}
	if !account.SubscriptionID.Valid {
		account.SubscriptionID.String = ""
	}
	if !account.TenantID.Valid {
		account.TenantID.String = ""
	}
	if !account.ClientID.Valid {
		account.ClientID.String = ""
	}
	if !account.ClientSecret.Valid {
		account.ClientSecret.String = ""
	}

	// Repeat this for other fields...

	return &account, nil
}
