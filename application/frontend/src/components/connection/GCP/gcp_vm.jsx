import { TokenFileWebIdentityCredentials } from 'aws-sdk';
import axios from 'axios';

// Function to create a GCP instance
export const createGCPInstance = async (instanceRequest) => {
  instanceRequest.accountID = parseInt(instanceRequest.accountID,10);
  try {
    const response = await axios.post('http://localhost:8080/gcp/vm/createInstance', instanceRequest);
    return response.data;
  } catch (error) {
    throw error;
  }
};

// Function to list GCP instances
export const listGCPInstances = async (accountID, accessToken) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/vm/listInstances', {
      accountID:accountID,
      token: accessToken
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

// Function to terminate a GCP instance
export const terminateGCPInstance = async (accountID, instanceID, token) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/vm/terminateInstance', terminateRequest);
    return response.data;
  } catch (error) {
    throw error;
  }
};
