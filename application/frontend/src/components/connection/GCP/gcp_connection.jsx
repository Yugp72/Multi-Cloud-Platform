// import { google } from 'googleapis';

// export const connectToGCP = async (credentialsArray) => {
//     try {
//         const results = [];
//         const connectedAccounts = [];

//         for (const credentials of credentialsArray) {
//             const auth = new google.auth.GoogleAuth({
//                 credentials: {
//                     client_email: credentials.clientEmail,
//                     private_key: credentials.privateKey,
//                 },
//                 scopes: ['https://www.googleapis.com/auth/cloud-platform'],
//             });

//             const client = google.cloud({
//                 version: 'v1',
//                 auth,
//             });

//             results.push({ success: true, message: 'Connected to GCP successfully' });

//             const accountDetails = {
//                 accountId: credentials.clientEmail,
//                 region: credentials.region || 'global',
//                 status: 'Connected',
//             };
//             connectedAccounts.push(accountDetails);
//         }

//         const storedAccounts = JSON.parse(sessionStorage.getItem('connectedAccounts') || '[]');
//         const finalupdateAccounts = [...storedAccounts, ...connectedAccounts];
//         sessionStorage.setItem('connectedAccounts', JSON.stringify(finalupdateAccounts));

//         await fetchGCPInstances(credentialsArray);

//         return { results, connectedAccounts };
//     } catch (error) {
//         const errorMessages = [{ success: false, message: `Failed to connect to GCP: ${error.message}` }];
//         return { results: errorMessages, connectedAccounts: [] };
//     }
// };

// export const listGCPInstances = async (credentials) => {
//     try {
//         const auth = new google.auth.GoogleAuth({
//             credentials: {
//                 client_email: credentials.clientEmail,
//                 private_key: credentials.privateKey,
//             },
//             scopes: ['https://www.googleapis.com/auth/cloud-platform'],
//         });

//         const client = google.compute({
//             version: 'v1',
//             auth,
//         });

//         const data = await client.instances.aggregatedList();
//         const instances = data.data.items;

//         const instanceNames = [];
//         for (const regionInstances of Object.values(instances)) {
//             for (const instance of regionInstances.instances || []) {
//                 instanceNames.push(instance.name);
//             }
//         }

//         return instanceNames;
//     } catch (error) {
//         console.error('Error listing GCP instances:', error);
//         return [];
//     }
// };

// const fetchGCPInstances = async (credentialsArray) => {
//     try {
//         for (const credentials of credentialsArray) {
//             const instances = await listGCPInstances(credentials);
//             console.log(`Fetched GCP instances for account ${credentials.clientEmail}:`, instances);
//             // Handle the fetched instances as needed
//         }
//     } catch (error) {
//         console.error('Error fetching GCP instances:', error);
//     }
// };

// // Add other GCP functions such as createGCPInstance, stopGCPInstance, etc.
