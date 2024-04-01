package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// InstanceRequest represents the JSON request structure for EC2 instance operations
type InstanceRequest struct {
	InstanceType     string   `json:"instanceType"`
	AmiID            string   `json:"amiID"`
	KeyName          string   `json:"keyName"`
	SecurityGroupIDs []string `json:"securityGroupIDs"`
	SubnetID         string   `json:"subnetID"`
	Region           string   `json:"region"`
	Name             string   `json:"name"`
}

// InstanceResponse represents the JSON response structure for EC2 instance operations
type InstanceResponse struct {
	Message     string   `json:"message"`
	InstanceIDs []string `json:"instanceIDs,omitempty"`
}

// Initialize AWS session
var sess *session.Session
var svc *ec2.EC2

func main() {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"), // Change to your desired region
	})
	if err != nil {
		log.Fatal("Error initializing AWS session:", err)
	}
	svc = ec2.New(sess)

	// Define routes
	http.HandleFunc("/createInstance", CreateInstanceHandler)
	http.HandleFunc("/listInstances", ListInstancesHandler)
	http.HandleFunc("/terminateInstance", TerminateInstanceHandler)
	http.HandleFunc("/terminal", handleTerminal)

	// Start the HTTP server
	port := "8200"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// CreateInstanceHandler handles POST requests to create an EC2 instance
func CreateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var req InstanceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	// Initialize AWS session with the provided region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.Region), // Use the region from the request
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error initializing AWS session: %v", err)
		return
	}
	svc := ec2.New(sess)

	// Create an EC2 instance
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(req.AmiID),
		InstanceType:     aws.String(req.InstanceType),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          aws.String(req.KeyName),
		SecurityGroupIds: aws.StringSlice(req.SecurityGroupIDs),
		SubnetId:         aws.String(req.SubnetID),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(req.Name), // Set your desired instance name
					},
				},
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating EC2 instance: %v", err)
		return
	}

	instanceID := *runResult.Instances[0].InstanceId
	resp := InstanceResponse{Message: fmt.Sprintf("EC2 instance created with ID: %s", instanceID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListInstancesHandler handles GET requests to list EC2 instances
func ListInstancesHandler(w http.ResponseWriter, r *http.Request) {
	result, err := svc.DescribeInstances(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error listing EC2 instances: %v", err)
		return
	}

	// Extract instance IDs and respond
	var instanceIDs []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// Append the instance ID to the slice
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}

	// Create the response struct with instance IDs
	resp := InstanceResponse{Message: instanceIDs[0], InstanceIDs: instanceIDs}

	// Encode the response struct to JSON
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding JSON response: %v", err)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type InstanceRequest1 struct {
	InstanceID string `json:"instanceID"`
}

// TerminateInstanceHandler handles POST requests to terminate an EC2 instance
func TerminateInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var req InstanceRequest1
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	if req.InstanceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Instance ID is required")
		return
	}

	_, err = svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(req.InstanceID)},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error terminating EC2 instance: %v", err)
		return
	}

	resp := InstanceResponse{Message: fmt.Sprintf("EC2 instance terminated: %s", req.InstanceID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleTerminal(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	// Load the private key associated with the EC2 instance
	keyFile := "FirstTryApi.pem"
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println("Read key file:", err)
		return
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Println("Parse private key:", err)
		return
	}

	// Establish SSH connection to EC2 instance
	config := &ssh.ClientConfig{
		User: "AWSLinux", // Specify the SSH username for your EC2 instance
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Insecure; use proper host key validation in production
	}

	client, err := ssh.Dial("tcp", "13.126.233.119:22", config)
	if err != nil {
		log.Println("SSH Dial:", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Println("Session:", err)
		return
	}
	defer session.Close()

	// Set up pipes to relay data between WebSocket and SSH session
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				return
			}
			session.Stdout.Write(message)
		}
	}()

	go func() {
		for {
			output := make([]byte, 1024)
			n, err := session.Stdin.Read(output)
			if err != nil {
				log.Println("Output:", err)
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.TextMessage, output[:n]); err != nil {
					log.Println("Write:", err)
					return
				}
			}
		}
	}()

	select {}
}
