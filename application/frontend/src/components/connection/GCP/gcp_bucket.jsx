import axios from 'axios';

export const createGCSBucket = async (bucketData) => {
    bucketData.accountID = parseInt(bucketData.accountID, 10);
    console.log("Bucket Data: ",bucketData);
    try {
    const response = await axios.post('http://localhost:8080/gcp/gcs/createBucket', {
        accountID: bucketData.accountID,
        bucketName: bucketData.bucketName,
        token: bucketData.accessToken
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const uploadGCSObject = async (objectData) => {
  try {
    const formData = new FormData();
    formData.append('accountID', objectData.accountID);
    formData.append('bucketName', objectData.bucketName);
    formData.append('objectName', objectData.objectName);
    formData.append('file', objectData.file);

    const response = await axios.post('http://localhost:8080/gcp/gcs/uploadObject', formData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getGCSObject = async (objectData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/gcs/getObject', objectData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteGCSObject = async (objectData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/gcs/deleteObject', objectData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteGCSBucket = async (accountID,bucketName,accessToken) => {
    // console.log("Delete bucket1211:",bucketData);
  try {
    const response = await axios.post('http://localhost:8080/gcp/gcs/deleteBucket', {
        accountID:accountID,
        bucketName: bucketName,
        token: accessToken
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listGCSBuckets = async (accountID, accessToken) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/gcs/listBuckets', { accountID: accountID, token: accessToken });
    return response.data;
  } catch (error) {
    throw error;
  }
};
