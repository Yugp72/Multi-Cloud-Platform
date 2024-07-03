import axios from 'axios';



export const createECS = async (clusterName, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/ecs/createCluster', {
            clusterName: clusterName,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const deleteECS = async (clusterName, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/ecs/deleteCluster', {
            clusterName: clusterName,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const listECS = async (accountID,region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/ecs/listClusters',{
            accountID: accountID,
            region: region
        
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

// Similarly, define functions for other AWS ECS APIs

export const createEKS = async (clusterName, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/eks/createCluster', {
            clusterName: clusterName,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const deleteEKS = async (clusterName, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/eks/deleteCluster', {
            clusterName: clusterName,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const listEKS = async (accountID, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/eks/listClusters',{
            accountID: accountID,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const createLambda = async (functionName, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/lambda/createFunction', {
            functionName: functionName,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const deleteLambda = async (functionName, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/lambda/deleteFunction', {
            functionName: functionName,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const listLambda = async () => {
    try {
        const response = await axios.post('http://localhost:8080/aws/lambda/listFunctions');
        return response.data;
    } catch (error) {
        throw error;
    }
}

