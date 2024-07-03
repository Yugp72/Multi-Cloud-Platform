const { InteractiveBrowserCredential,ClientSecretCredential } = require('@azure/identity');
const { ComputeManagementClient } = require('@azure/arm-compute');
import { PublicClientApplication } from '@azure/msal-browser';
// const axios = require('axios');


// app.use(cors());


export const connectToAzure = async (credentialsArray) => {
    try {
        const connectedAccounts = [];
        
        for (const credentials of credentialsArray) {
            const { clientId, clientSecret, tenantId } = credentials; 
            
            const credential = new InteractiveBrowserCredential({ clientId, clientSecret, tenantId }); // Provide the credentials
            const client = new ComputeManagementClient(credential, credentials.subscriptionId);

            const virtualMachines = await listVirtualMachines(client);

            connectedAccounts.push({
                subscriptionId: credentials.subscriptionId,
                virtualMachines: virtualMachines,
                status: 'Connected',
            });
            
        }
        const storedAccounts = JSON.parse(sessionStorage.getItem('connectedAccounts') || '[]'); // Parse existing accounts or initialize as empty array
        const finalupdateAccounts = [...storedAccounts, ...connectedAccounts]; // Merge new and existing accounts

        sessionStorage.setItem('connectedAccounts', JSON.stringify(finalupdateAccounts)); // Update sessionStorage

            return { success: true, message: 'Connected to Azure successfully', connectedAccounts };
        } catch (error) {
            return { success: false, message: `Failed to connect to Azure: ${error.message}` };
        }
    };

//     export async function listVirtualMachines(tenantId, clientId, clientSecret, subscriptionId) {
//         console.log("tenantId: ", tenantId);
//         console.log("clientId: ", clientId);
//         console.log("clientSecret: ", clientSecret);
//         const clientSecretCredential = new ClientSecretCredential(
//             tenantId, // your Azure Active Directory Tenant ID
//             clientId, // your Azure Active Directory Client ID
//             "Ovq8Q~S8kDfqdvS3pFGj1V2sOs98vJPB8ZsLHc7c" // your Azure Active Directory Client Secret
//         );
    
//         const computeClient = new ComputeManagementClient(clientSecretCredential, subscriptionId);
    
//         const result = await computeClient.virtualMachines.listAll();
//         console.log("this is results: ",result);
    
//         console.log("List of VMs:");
//         for await (const vm of result) {
//             console.log(`Name: ${vm.name}, Location: ${vm.location}, PowerState: ${vm.instanceView.statuses[1].displayStatus}`);
//         }
//     }
    
//     listVirtualMachines("<tenantId>", "<clientId>", "<clientSecret>", "<subscriptionId>")
//         .catch((err) => console.error("Error listing VMs: ", err));

    

//     // export const initializeMsal = async (clientId, tenantId) => {
//     //     const msalConfig = {
//     //         auth: {
//     //             clientId: clientId,
//     //             authority: `https://login.microsoftonline.com/${tenantId}`
//     //         }
//     //     };

//     //     const msalInstance =  new PublicClientApplication(msalConfig);
//     //     console.log("msalInstance: ", msalInstance);

//     //     return msalInstance;
//     // };
//     // export const listVirtualMachines = async (credentials) => {
//     //     try {
//     //         const msalInstance = await initializeMsal(credentials.clientId, credentials.tenantId);
           
//     //         console.log("credentials: ", credentials);
    
//     //         const tokenRequest = {
//     //             scopes: ['https://management.azure.com/.default']
//     //         };
//     //         // if (!msalInstance.getAccount()) {
//     //         //     await initializeMsal(credentials.clientId, credentials.tenantId);
//     //         // }
//     //         // console.log("msalInstance.getAccount(): ", msalInstance.getAccount());
            
//     //         const tokenResponse =  await msalInstance.acquireTokenSilent(tokenRequest);
//     //         console.log("tokenResponse: ", tokenResponse);
    
//     //         const accessToken = tokenResponse.accessToken;
//     //         console.log("accessToken: ", accessToken);
    
//     //         const { subscriptionId } = credentials;
    
//     //         const apiVersion = '2024-03-01';
//     //         const url = `https://management.azure.com/subscriptions/${subscriptionId}/providers/Microsoft.Compute/virtualMachines?api-version=${apiVersion}`;
//     //         const headers = {
//     //             'Content-Type': 'application/json',
//     //             'Authorization': `Bearer ${accessToken}`
//     //         };
    
//     //         const response = await axios.get(url, { headers });
//     //         return response.data.value;
//     //     } catch (error) {
//     //         console.error("Error fetching VM instances:", error.message);
//     //         throw error;
//     //     }
//     // };
// export const createVirtualMachine = async (client, resourceGroupName, vmName, vmParams) => {
//     try {
//         const virtualMachine = await client.virtualMachines.createOrUpdate(resourceGroupName, vmName, vmParams);
//         return virtualMachine;
//     } catch (error) {
//         console.error('Error creating virtual machine:', error);
//         throw new Error('Error creating virtual machine');
//     }
// };

// // Function to stop a virtual machine
// export const stopVirtualMachine = async (client, resourceGroupName, vmName) => {
//     try {
//         const operation = await client.virtualMachines.powerOff(resourceGroupName, vmName);
//         return operation;
//     } catch (error) {
//         console.error('Error stopping virtual machine:', error);
//         throw new Error('Error stopping virtual machine');
//     }
// };

// // Function to restart a virtual machine
// export const restartVirtualMachine = async (client, resourceGroupName, vmName) => {
//     try {
//         const operation = await client.virtualMachines.start(resourceGroupName, vmName);
//         return operation;
//     } catch (error) {
//         console.error('Error restarting virtual machine:', error);
//         throw new Error('Error restarting virtual machine');
//     }
// };

// // Function to delete a virtual machine
// export const deleteVirtualMachine = async (client, resourceGroupName, vmName) => {
//     try {
//         const operation = await client.virtualMachines.beginDelete(resourceGroupName, vmName);
//         await operation.pollUntilFinished();
//         return operation;
//     } catch (error) {
//         console.error('Error deleting virtual machine:', error);
//         throw new Error('Error deleting virtual machine');
//     }
// };

// // Function to remove an account from the list of connected accounts
// export const removeAccount = (subscriptionId, connectedAccounts) => {
//     const updatedAccounts = connectedAccounts.filter(account => account.subscriptionId !== subscriptionId);
//     return updatedAccounts;
// };

// export const getAzureCredentialsArray= (credentialsArray) => {
//     return credentialsArray;
// };

// // Function to get connected Azure accounts
// export const getConnectedAzureAccounts = async (credentialsArray) => {
//     try {
//         const { connectedAccounts } = await connectToAzure(credentialsArray);
//         return connectedAccounts;
//     } catch (error) {
//         console.error('Error fetching connected Azure accounts:', error);
//         return [];
//     }
// };
