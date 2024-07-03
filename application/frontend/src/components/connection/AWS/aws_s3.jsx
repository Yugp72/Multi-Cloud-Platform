import axios from 'axios';


export const createS3Bucket = async (bucketData) => {
  bucketData.accountID = parseInt(bucketData.accountID, 10);
  try {
    const response = await axios.post('http://localhost:8080/aws/s3/createBucket', {
      accountID: bucketData.accountID,
      region: bucketData.region,
      bucketName: bucketData.bucketName
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};


export const uploadS3Object = async (objectData) => {
  try {
    const formData = new FormData();
    formData.append('accountID', parseInt(objectData.accountID, 10));
    formData.append('region', objectData.region);
    formData.append('bucketName', objectData.bucketName);
    formData.append('objectKey', objectData.objectKey);
    formData.append('content', objectData.content);

    const response = await axios.post('http://localhost:8080/aws/s3/uploadObject', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};


export const getS3Object = async (objectData) => {
  try {
    const response = await axios.get('http://localhost:8080/aws/s3/getObject', { data: objectData });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteS3Object = async (objectData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/s3/deleteObject', objectData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteS3Bucket = async (accountID, name,region) => {
  console.log("Delete bucket1211:",accountID,name,region)
  try {
    const response = await axios.post('http://localhost:8080/aws/s3/deleteBucket', {
      accountID: accountID,
      bucketName: name,
      region: region
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listS3Buckets = async (accountID, region) => {
  try {
    // console.log('listBucketRequest:', listBucketRequest);
    const response = await axios.post('http://localhost:8080/aws/s3/listBuckets', {accountID:accountID, region:region});
    console.log(response.data);
    return response.data.buckets;
  } catch (error) {
    throw error;
  }
};
