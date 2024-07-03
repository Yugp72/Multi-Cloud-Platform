import React, { useState, useEffect } from "react";
import { TextInput } from "@mantine/core";
import { Text, Button } from "@mantine/core";
import { connectToAzure } from "./azure_connection";

const AzureCloudConnection = () => {
  const [clientId, setClientId] = useState('');
  const [clientSecret, setClientSecret] = useState('');
  const [tenantId, setTenantId] = useState('');
  const [subscriptionId, setSubsctiptionId] = useState('');
  const [loading, setLoading] = useState(false);
  const [messages, setMessages] = useState(['']);

  // UseEffect to handle sessionStorage
  useEffect(() => {
    // Load credentials from sessionStorage if available
    const storedClientId = sessionStorage.getItem('azureClientId');
    const storedClientSecret = sessionStorage.getItem('azureClientSecret');
    const storedTenantId = sessionStorage.getItem('azureTenantId');
    const storedSubscriptionId = sessionStorage.getItem('azureSubscriptionId');
    if (storedClientId && storedClientSecret && storedTenantId && storedSubscriptionId) {
      setClientId(storedClientId);
      setClientSecret(storedClientSecret);
      setTenantId(storedTenantId);
      setSubsctiptionId(storedSubscriptionId);
    }
  }, []);

  // Save credentials to sessionStorage whenever they change
  useEffect(() => {
    sessionStorage.setItem('azureClientId', clientId);
    sessionStorage.setItem('azureClientSecret', clientSecret);
    sessionStorage.setItem('azureTenantId', tenantId);
    sessionStorage.setItem('azureSubscriptionId', subscriptionId);
  }, [clientId, clientSecret, tenantId, subscriptionId]);

  const handleConnect = async () => {
    setLoading(true);
    setMessages(['']);

    // Validate inputs
    if (!clientId || !clientSecret || !tenantId || !subscriptionId) {
      setMessages(['Please provide all required Azure credentials.']);
      setLoading(false);
      return;
    }

    try {
      // Connect to Azure
      const results = await connectToAzure([{ clientId, clientSecret, tenantId, subscriptionId }]);
      // Update messages
      setMessages(results.success ? results.message : results);
      setMessages('Connected to Azure successfully');
      console.log("Azure result array:", results);
    } catch (error) {
      // Set error message
      setMessages('Failed to connect to Azure');
      console.error("Error connecting to Azure:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <Text className="title" size="lg">Connect to Azure Cloud Provider</Text>
      <div className="inputGroup">
        <label className="label">Client ID</label>
        <TextInput
          className="input"
          placeholder="Enter your Azure client ID"
          value={clientId}
          onChange={(event) => setClientId(event.target.value)}
          disabled={loading}
        />
        <label className="label">Client Secret</label>
        <TextInput
          className="input"
          placeholder="Enter your Azure client secret"
          value={clientSecret}
          onChange={(event) => setClientSecret(event.target.value)}
          disabled={loading}
        />
        <label className="label">Subscription ID</label>
        <TextInput
          className="input"
          placeholder="Enter your Azure client subscription Id"
          value={subscriptionId}
          onChange={(event) => setSubsctiptionId(event.target.value)}
          disabled={loading}
        />
        <label className="label">Tenant ID</label>
        <TextInput
          className="input"
          placeholder="Enter your Azure tenant ID"
          value={tenantId}
          onChange={(event) => setTenantId(event.target.value)}
          disabled={loading}
        />
      </div>
      <Button className="button" onClick={handleConnect} disabled={loading}>
        {loading ? 'Connecting...' : 'Connect'}
      </Button>
      <Text className="message" color={messages.includes('Failed') ? 'red' : 'green'}>
        {messages}
      </Text>
    </div>
  );
};

export default AzureCloudConnection;
