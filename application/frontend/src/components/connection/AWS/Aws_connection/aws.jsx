import React, { useState, useEffect } from "react";
import { TextInput } from "@mantine/core";
import { Text, Button } from "@mantine/core";
import { connectToAWS } from '../aws_connection';
import { gql, useMutation } from '@apollo/client';
import { jwtDecode } from 'jwt-decode';



const AWSCloudConnection = () => {
  
  const decodedToken = jwtDecode(localStorage.getItem('token'));
  console.log("prinding decodedToken" , decodedToken);
  const UserID = decodedToken.userID;
  console.log("printing userid",UserID);
  const [accessKeys, setAccessKeys] = useState(['']);
  const [secretKeys, setSecretKeys] = useState(['']);
  const [region, setRegion] = useState(['']);
  const [loading, setLoading] = useState(false);
  const [messages, setMessages] = useState('');
  
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
        ClientEmail,
        PrivateKey,
        ProjectID
      }
    }
  `;
const [addCloudAccount] = useMutation(ADD_CLOUD_ACCOUNT);


  useEffect(() => {
    const storedAccessKeys = sessionStorage.getItem('accessKeys');
    const storedSecretKeys = sessionStorage.getItem('secretKeys');
    if (storedAccessKeys && storedSecretKeys) {
      setAccessKeys(JSON.parse(storedAccessKeys));
      setSecretKeys(JSON.parse(storedSecretKeys));
    }
  }, []);

  useEffect(() => {
    sessionStorage.setItem('accessKeys', JSON.stringify(accessKeys));
    sessionStorage.setItem('secretKeys', JSON.stringify(secretKeys));
  }, [accessKeys, secretKeys]);

  const handleConnect = async () => {
    setLoading(true);
    setMessages('');

    if (accessKeys.some(key => key === '') || secretKeys.some(key => key === '') || region.some(r => r === '')) {
      setMessages('Please provide access key, secret key, and region for each AWS account.');
      setLoading(false);
      return;
    }

    try {
      const storedAccounts = JSON.parse(sessionStorage.getItem('connectedAccounts') || '[]');

      const existingAccounts = storedAccounts.filter(storedAccount =>
        accessKeys.includes(storedAccount.AccessKey) && secretKeys.includes(storedAccount.SecretKey)
      );

      if (existingAccounts.length > 0) {
        setMessages('You have already connected to one or more of these AWS accounts.');
        setLoading(false);
        return;
      }

      const results = await Promise.all(accessKeys.map(async (accessKey, index) => {
        const credentials = {
          accessKeyId: accessKey,
          secretAccessKey: secretKeys[index],
          region: region[index]
        };

        const result = await connectToAWS([credentials]); // Connect to AWS for this account
        const accountDetails = {
          AccessKey: accessKey,
          SecretKey: secretKeys[index],
          Region: region[index],
          UserID: UserID,
          CloudProvider: 'AWS',
          // Add other account details here if needed
        };
        console.log("Printing account details: ",accountDetails);

        await addCloudAccount({ variables: { cloudAccountInput: accountDetails } });
        

        return result;
      }));

      setMessages('Connected to AWS successfully');

      return results;
    } catch (error) {
      setMessages('Failed to connect to AWS');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <Text className="title" size="lg">Connect to AWS Cloud Provider</Text>
      {accessKeys.map((accessKey, index) => (
        <div key={index} className="inputGroup">
          <label className="label">Access Key</label>
          <TextInput
            className="input"
            placeholder="Enter your AWS access key"
            value={accessKey}
            onChange={(event) => {
              const newAccessKeys = [...accessKeys];
              newAccessKeys[index] = event.target.value;
              setAccessKeys(newAccessKeys);
            }}
            disabled={loading}
          />
          <label className="label">Secret Key</label>
          <TextInput
            className="input"
            placeholder="Enter your AWS secret key"
            value={secretKeys[index]}
            onChange={(event) => {
              const newSecretKeys = [...secretKeys];
              newSecretKeys[index] = event.target.value;
              setSecretKeys(newSecretKeys);
            }}
            disabled={loading}
          />
          <label className="label">Region</label>
          <TextInput
            className="input"
            placeholder="Enter your AWS region"
            value={region[index]}
            onChange={(event) => {
              const newRegion = [...region];
              newRegion[index] = event.target.value;
              setRegion(newRegion); 
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

export default AWSCloudConnection;
