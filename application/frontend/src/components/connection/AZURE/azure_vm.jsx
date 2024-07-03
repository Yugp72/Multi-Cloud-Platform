import axios from 'axios';

// Function to create an Azure VM instance
export const createAzureVMInstance = async (instanceData) => {
  instanceData.accountID = parseInt(instanceData.accountID, 10);
  try {
    const response = await axios.post('http://localhost:8080/azure/vm/createInstance', {
      accountID: instanceData.accountID,
      // Include other necessary data for creating the VM instance
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

// Function to list Azure VM instances
export const listAzureVMInstances = async (accountID,AccessTokenAZURE) => {
  accountID = parseInt(accountID, 10);
  try {
    const response = await axios.post('http://localhost:8080/azure/vm/listInstances', {
      accountID: accountID,
      token: AccessTokenAZURE
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

// Function to terminate an Azure VM instance
export const terminateAzureVMInstance = async (instanceID) => {
  instanceID = parseInt(instanceID, 10);
  try {
    const response = await axios.post('http://localhost:8080/azure/vm/terminateInstance', {
      instanceID: instanceID
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};
