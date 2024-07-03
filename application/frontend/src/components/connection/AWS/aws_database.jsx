import axios from 'axios';

// Function to create an item in DynamoDB
const createDynamoDBItem = async (itemData) => {
    try {
        const response = await axios.post("/aws/dynamodb/createItem", itemData);
        return response.data;
    } catch (error) {
        console.error("Error creating item in DynamoDB:", error);
        throw error;
    }
};

// Function to read an item from DynamoDB
const readDynamoDBItem = async (itemId) => {
    try {
        const response = await axios.get(`/aws/dynamodb/readItem?id=${itemId}`);
        return response.data;
    } catch (error) {
        console.error("Error reading item from DynamoDB:", error);
        throw error;
    }
};

// Function to delete an item from DynamoDB
const deleteDynamoDBItem = async (itemId) => {
    try {
        const response = await axios.post("/aws/dynamodb/deleteItem", { id: itemId });
        return response.data;
    } catch (error) {
        console.error("Error deleting item from DynamoDB:", error);
        throw error;
    }
};

// Function to update an item in DynamoDB
const updateDynamoDBItem = async (itemId, updatedData) => {
    try {
        const response = await axios.put(`/aws/dynamodb/updateItem?id=${itemId}`, updatedData);
        return response.data;
    } catch (error) {
        console.error("Error updating item in DynamoDB:", error);
        throw error;
    }
};

// Function to create a table in DynamoDB
const createDynamoDBTable = async (tableData) => {
    try {
        const response = await axios.post("/aws/dynamodb/createTable", tableData);
        return response.data;
    } catch (error) {
        console.error("Error creating DynamoDB table:", error);
        throw error;
    }
};

// Function to delete a table from DynamoDB
const deleteDynamoDBTable = async (tableName) => {
    try {
        const response = await axios.post("/aws/dynamodb/deleteTable", { name: tableName });
        return response.data;
    } catch (error) {
        console.error("Error deleting DynamoDB table:", error);
        throw error;
    }
};

// Function to update a table in DynamoDB
const updateDynamoDBTable = async (tableName, updatedData) => {
    try {
        const response = await axios.put(`/aws/dynamodb/updateTable?name=${tableName}`, updatedData);
        return response.data;
    } catch (error) {
        console.error("Error updating DynamoDB table:", error);
        throw error;
    }
};

// Function to list all tables in DynamoDB
const listDynamoDBTables = async () => {
    try {
        const response = await axios.get("/aws/dynamodb/listTables");
        return response.data;
    } catch (error) {
        console.error("Error listing DynamoDB tables:", error);
        throw error;
    }
};

export {
    createDynamoDBItem,
    readDynamoDBItem,
    deleteDynamoDBItem,
    updateDynamoDBItem,
    createDynamoDBTable,
    deleteDynamoDBTable,
    updateDynamoDBTable,
    listDynamoDBTables,
};
