import AWS from 'aws-sdk';
const axios = require('axios');

const apiClient = axios.create({
    baseURL: 'http://localhost:8080',
});

export const connectToAWS = async (credentialsArray) => {
    //condition that if the user already have a same connectAccount then it would not be again connect

    try {
        const results = [];
        const connectedAccounts = [];

        for (const credentials of credentialsArray) {
            const first=AWS.config.update({
                accessKeyId: credentials.accessKeyId,
                secretAccessKey: credentials.secretAccessKey,
                region: credentials.region || 'ap-south-1',
            });

            results.push({ success: true, message: 'Connected to AWS successfully'});
            const accountDetails = {
                accountId: credentials.accessKeyId,
                region: credentials.region || 'ap-south-1',
                status: 'Connected',
            };
            connectedAccounts.push(accountDetails);
        }
        
        const storedAccounts = JSON.parse(sessionStorage.getItem('connectedAccounts') || '[]'); // Parse existing accounts or initialize as empty array
        const finalupdateAccounts = [...storedAccounts, ...connectedAccounts]; // Merge new and existing accounts
        
        sessionStorage.setItem('connectedAccounts', JSON.stringify(finalupdateAccounts)); // Update sessionStorage
        
        await fetchEC2Instances(credentialsArray);


        return { results, connectedAccounts };
    } catch (error) {
        const errorMessages = [{ success: false, message: `Failed to connect to AWS: ${error.message}` }];
        return { results: errorMessages, connectedAccounts: [] };
    }
};


export const getCredentialsArray = (credentialsArray) => {
    return credentialsArray;
};


