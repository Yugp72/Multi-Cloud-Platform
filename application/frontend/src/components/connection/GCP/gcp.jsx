import React, { useState, useEffect } from "react";
import { TextInput } from "@mantine/core";
import { Text, Button } from "@mantine/core";
import { connectToGCP } from './gcp_connection';
import { gql, useMutation } from '@apollo/client';
import { jwtDecode } from 'jwt-decode';

const GCPConnection = () => {
  
  // Decode the token to get the UserID
  const decodedToken = jwtDecode(localStorage.getItem('token'));
  const UserID = decodedToken.userID;

  // State variables for input fields
  const [serviceAccounts, setServiceAccounts] = useState(['']);
  const [privateKeys, setPrivateKeys] = useState(['']);
  const [projectId, setProjectId] = useState(['']);
  const [loading, setLoading] = useState(false);
  const [messages, setMessages] = useState('');

  // GraphQL mutation to add a cloud account
  const ADD_CLOUD_ACCOUNT = gql`
    mutation AddCloudAccount($cloudAccountInput: CloudAccountInput!) {
      addCloudAccount(cloudAccountInput: $cloudAccountInput) {
        AccessKey
        AccountID
        AdditionalInformation
        ClientID
        ClientSecret
        CloudProvider
        Region
        SecretKey
        TenantID
        SubscriptionID
        UserID
        ClientEmail
        PrivateKey
        ProjectID
      }
    }
  `;

  const [addCloudAccount] = useMutation(ADD_CLOUD_ACCOUNT);

  // Effect hook to load stored access keys and secret keys
  useEffect(() => {
    const storedServiceAccounts = sessionStorage.getItem('serviceAccounts');
    const storedPrivateKeys = sessionStorage.getItem('privateKeys');
    const storedProjectIds = sessionStorage.getItem('projectIds');

    if (storedServiceAccounts && storedPrivateKeys && storedProjectIds) {
      setServiceAccounts(JSON.parse(storedServiceAccounts));
      setPrivateKeys(JSON.parse(storedPrivateKeys));
      setProjectId(JSON.parse(storedProjectIds));
    }
  }, []);

  // Effect hook to store access keys and secret keys in session storage
  useEffect(() => {
    sessionStorage.setItem('serviceAccounts', JSON.stringify(serviceAccounts));
    sessionStorage.setItem('privateKeys', JSON.stringify(privateKeys));
    sessionStorage.setItem('projectIds', JSON.stringify(projectId));
  }, [serviceAccounts, privateKeys, projectId]);

  // Function to handle connection to GCP
  const handleConnect = async () => {
    setLoading(true);
    setMessages('');

    if (serviceAccounts.some(account => account === '') || privateKeys.some(key => key === '') || projectId.some(id => id === '')) {
      setMessages('Please provide service account, private key, and project ID for each GCP account.');
      setLoading(false);
      return;
    }

    try {
      // Connect to GCP for each account
      const results = await Promise.all(serviceAccounts.map(async (serviceAccount, index) => {
        const credentials = {
          serviceAccount,
          privateKey: privateKeys[index],
          projectId: projectId[index]
        };

        const result = await connectToGCP(credentials); // Connect to GCP for this account

        // Prepare account details to be stored in the database
        const accountDetails = {
          ProjectID: serviceAccount,
          PrivateKey: privateKeys[index],
          ProjectID: projectId[index],
          UserID: UserID,
          CloudProvider: 'GCP',
          // Add other account details here if needed
        };

        // Add the cloud account to the database
        await addCloudAccount({ variables: { cloudAccountInput: accountDetails } });

        return result;
      }));

      setMessages('Connected to GCP successfully');

      return results;
    } catch (error) {
      setMessages('Failed to connect to GCP');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <Text className="title" size="lg">Connect to GCP Cloud Provider</Text>
      {serviceAccounts.map((account, index) => (
        <div key={index} className="inputGroup">
          <label className="label">Service Account</label>
          <TextInput
            className="input"
            placeholder="Enter your GCP service account"
            value={account}
            onChange={(event) => {
              const newServiceAccounts = [...serviceAccounts];
              newServiceAccounts[index] = event.target.value;
              setServiceAccounts(newServiceAccounts);
            }}
            disabled={loading}
          />
          <label className="label">Private Key</label>
          <TextInput
            className="input"
            placeholder="Enter your GCP private key"
            value={privateKeys[index]}
            onChange={(event) => {
              const newPrivateKeys = [...privateKeys];
              newPrivateKeys[index] = event.target.value;
              setPrivateKeys(newPrivateKeys);
            }}
            disabled={loading}
          />
          <label className="label">Project ID</label>
          <TextInput
            className="input"
            placeholder="Enter your GCP project ID"
            value={projectId[index]}
            onChange={(event) => {
              const newProjectId = [...projectId];
              newProjectId[index] = event.target.value;
              setProjectId(newProjectId);
            }}
            disabled={loading}
          />
        </div>
      ))}
      <Button className="button" onClick={handleConnect} disabled={loading}>
        {loading ? 'Connecting...' : 'Connect'}
      </Button>
      <Text className="message" color={messages.includes('Failed') ? 'red' : 'green'}>
        {messages}
      </Text>
    </div>
  );
};

export default GCPConnection;
