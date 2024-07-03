import axios from 'axios';

export const createNetwork = async (networkData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/network/createNetwork', networkData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listNetworks = async () => {
  try {
    const response = await axios.get('http://localhost:8080/gcp/network/listNetworks');
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteNetwork = async (networkData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/network/deleteNetwork', networkData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const createSubnet = async (subnetData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/network/createSubnet', subnetData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteSubnet = async (subnetData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/network/deleteSubnet', subnetData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listSubnets = async () => {
  try {
    const response = await axios.get('http://localhost:8080/gcp/network/listSubnets');
    return response.data;
  } catch (error) {
    throw error;
  }
};

//createRoute
export const createRoute = async (routeData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/router/createRoute', routeData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

//deleteRoute
export const deleteRoute = async (routeData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/router/deleteRoute', routeData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

//listRoutes
export const listRoutes = async () => {
  try {
    const response = await axios.get('http://localhost:8080/gcp/router/listRoutes');
    return response.data;
  } catch (error) {
    throw error;
  }
};

//createFirewallRule
export const createFirewallRule = async (firewallRuleData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/router/createFirewallRule', firewallRuleData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

//deleteFirewallRule
export const deleteFirewallRule = async (firewallRuleData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/router/deleteFirewallRule', firewallRuleData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

//listFirewallRules
export const listFirewallRules = async () => {
  try {
    const response = await axios.get('http://localhost:8080/gcp/router/listFirewallRules');
    return response.data;
  } catch (error) {
    throw error;
  }
};

//createCloudRouter
export const createCloudRouter = async (cloudRouterData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/router/createCloudRouter', cloudRouterData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

//deleteCloudRouter
export const deleteCloudRouter = async (cloudRouterData) => {
  try {
    const response = await axios.post('http://localhost:8080/gcp/router/deleteCloudRouter', cloudRouterData);
    return response.data;
  } catch (error) {
    throw error;
  }
};


export const listCloudRouters = async () => {
  try {
    const response = await axios.get('http://localhost:8080/gcp/router/listCloudRouters');
    return response.data;
  } catch (error) {
    throw error;
  }
};
