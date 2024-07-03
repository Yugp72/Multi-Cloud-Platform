import axios from 'axios';

// Define API functions for each endpoint
export const createVPC = async (vpcData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/network/createVPC', vpcData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const listVPCs = async (accountID,region) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/network/listVPCs',{
        accountID: accountID   ,
        region: region 
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteVPC = async (accountID, id, region) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/network/deleteVPC', {
        vpcId: id,
        region: region,
        accountID: accountID
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const createSubnet = async (subnetData) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/network/createSubnet', subnetData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const updateSubnet = async (subnetData) => {
  try {
    const response = await axios.put('http://localhost:8080/aws/network/updateSubnet', subnetData);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteSubnet = async (accountID, region, subnetId) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/network/deleteSubnet', {
        accountID: accountID,
        region: region,
        subnetId: subnetId
    
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};


export const listSubnets = async (accountID, region) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/network/listSubnet',{
            accountID: accountID,
            region: region
        });
        return response.data;
    } catch (error) {
        throw error;
    }
    }

//createRouteTable
export const createRouteTable = async (routeTableData) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/network/createRouteTable', routeTableData);
        return response.data;
    } catch (error) {
        throw error;
    }
    }

//deleteRouteTable
export const deleteRouteTable = async (routeTableData) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/network/deleteRouteTable', routeTableData);
        return response.data;
    } catch (error) {
        throw error;
    }
    }

//listRouteTables
export const listRouteTables = async () => {
    try {
        const response = await axios.get('http://localhost:8080/aws/network/listRouteTables');
        return response.data;
    } catch (error) {
        throw error;
    }
    }

export const createInternetGateway = async (internetGatewayData) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/network/createInternetGateway', internetGatewayData);
        return response.data;
    } catch (error) {
        throw error;
    }
    }

export const attachInternetGateway = async (accountID, region, vpcId, internetGatewayId) => {
   //conver to integer
    accountID = parseInt(accountID, 10);
    try {
        const response = await axios.post('http://localhost:8080/aws/network/attachInternetGateway', {
            accountID: accountID,
            region: region,
            vpcId: vpcId,
            internetGatewayId: internetGatewayId
        });
        return response.data;
    } catch (error) {
        throw error;
    }
    }

export const detachInternetGateway = async (accountID, region, vpcId, internetGatewayId) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/network/detachInternetGateway', {
            accountID: accountID,
            region: region,
            vpcId: vpcId,
            internetGatewayId: internetGatewayId
        });
        return response.data;
    } catch (error) {
        throw error;
    }
    }

export const deleteInternetGateway = async (internetGatewayData) => {
    try {
        const response = await axios.post('http://localhost:8080/aws/network/deleteInternetGateway', internetGatewayData);
        return response.data;
    } catch (error) {
        throw error;
    }
    }


export const listInternetGateways = async (accountID, region) => {
  try {
    const response = await axios.post('http://localhost:8080/aws/network/listInternetGateway',{
        accountID: accountID,
        region: region
    });
    console.log(response.data);
    return response.data;
  } catch (error) {
    throw error;
  }
};
