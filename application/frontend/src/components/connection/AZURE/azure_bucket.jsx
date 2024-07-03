import axios from 'axios';

export const createAzureStorageAccount = async (accountData) => {
  try {
    const response = await axios.post('http://localhost:8080/azure/storage/createAccount', accountData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteAzureStorageAccount = async (accountData) => {
  try {
    const response = await axios.post('http://localhost:8080/azure/storage/deleteAccount', accountData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const updateAzureStorageAccount = async (accountData) => {
  try {
    const response = await axios.post('http://localhost:8080/azure/storage/updateStorageAccount', accountData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listAzureStorageAccounts = async (accountID, token) => {
  try {
    const response = await axios.post('http://localhost:8080/azure/storage/listAccounts', {
        accountID: accountID,
        token: token
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getAzureStorageObjects = async (objectData) => {
  try {
    const response = await axios.post('http://localhost:8080/azure/storage/getObjects', objectData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const uploadAzureStorageObjects = async (objectData) => {
  try {
    const response = await axios.post('http://localhost:8080/azure/storage/uploadObjects', objectData);
    return response.data;
  } catch (error) {
    throw error;
  }
};
