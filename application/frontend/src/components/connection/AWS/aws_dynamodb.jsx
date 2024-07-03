import axios from 'axios';

export const createDynamoDBItem = async (itemData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/createItem', itemData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const readDynamoDBItem = async (keyData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/readItem', keyData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteDynamoDBItem = async (keyData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/deleteItem', keyData);
    return response.data;
  } catch (error) {
    throw error;
  }
};


export const listDynamoDBItems = async (tableData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/listItems', tableData);
    return response.data;
  } catch (error) {
    throw error;
  }
}
export const updateDynamoDBItem = async (itemData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/updateItem', itemData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const createDynamoDBTable = async (tableData) => {
  const transformedObject = {
    tableName: tableData.tableName,
    region: tableData.region,
    accountID: parseInt(tableData.accountID),
    attributeDefinitions: [
        {
            attributeName: tableData.attributeName,
            attributeType: tableData.attributeType
        }
    ],
    keySchema: [
        {
            keyType: tableData.keySchema,
            attributeName: tableData.attributeName
        }
    ],
    provisionedThroughput: {
        ReadCapacityUnits: parseInt(tableData.provisionedThroughput),
        WriteCapacityUnits: parseInt(tableData.provisionedThroughput)
    }
};
  console.log("Table from api:",transformedObject );
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/createTable', transformedObject);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteDynamoDBTable = async (accountID,name,region) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/deleteTable', {
      accountID: accountID,
      tableName: name,
      region: region
    
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const updateDynamoDBTable = async (tableData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/updateTable', tableData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listDynamoDBTables = async (data) => {
  console.log(data);
  try {
    const response = await axios.post('http://localhost:8080/aws/dynamodb/listTables',{
      region: data.region,
      accountID: data.accountID
    });
    console.log(response.data);
    return response.data;
  } catch (error) {
    throw error;
  }
};
