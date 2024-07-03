// components/connection/provider.js

import { connectToAWS, getCredentialsArray as getAWSCredentials } from '@/components/connection/AWS/aws_connection';
import { connectToAzure , getAzureCredentialsArray as getAzureCredentials} from '@/components/connection/AZURE/azure_connection';

export const connectToProvider = async (provider) => {
    switch (provider) {
        case 'AWS':
            return await connectToAWS(getAWSCredentials());
        case 'Azure':
            return await connectToAzure(getAzureCredentials());
        // Add more cases for other providers if needed
        default:
            throw new Error(`Provider ${provider} not supported.`);
    }
};
