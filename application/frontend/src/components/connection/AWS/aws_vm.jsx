// "use client";
// import { useState } from "react";
// import { TextInput } from "@mantine/core";
// import { Text, Button } from "@mantine/core";
// import { connectToAWS } from '../aws_connection';
// import styles from './aws.css';

// const AWSCloudConnection = () => {
//   const [accessKeys, setAccessKeys] = useState(['']);
//   const [secretKeys, setSecretKeys] = useState(['']);
//   const [loading, setLoading] = useState(false);
//   const [messages, setMessages] = useState(['']);

//   const handleConnect = async () => {
//     setLoading(true);
//     setMessages(['']);

//     // Validate access keys and secret keys
//     if (accessKeys.some(key => key === '') || secretKeys.some(key => key === '')) {
//       setMessages(['Please provide both access key and secret key for each AWS account.']);
//       setLoading(false);
//       return;
//     }

//     // Construct array of credentials objects
//     const credentialsArray = accessKeys.map((accessKey, index) => ({
//       accessKeyId: accessKey,
//       secretAccessKey: secretKeys[index],
//     }));

//     try {
//       const results = await connectToAWS(credentialsArray);

//       // Update messages for each connection attempt
//       setMessages(results.results.messages);
//       setMessages('Connected to AWS successfully');
//       console.log(results);
//         } catch (error) {
//       // Set error message for all connection attempts
//       setMessages('Failed to connect to AWS');
//       // setMessages(credentialsArray.map(() => `Failed to connect to AWS: ${error.message}`));
//     } finally {
//       setLoading(false);
//     }
//   };

//   return (
//     <div className="container">
//       <Text className="title" size="lg">Connect to AWS Cloud Provider</Text>
//       {accessKeys.map((accessKey, index) => (
//         <div key={index} className="inputGroup">
//           <label className="label">Access Key</label>
//           <TextInput
//             className="input"
//             placeholder="Enter your AWS access key"
//             value={accessKey}
//             onChange={(event) => {
//               const newAccessKeys = [...accessKeys];
//               newAccessKeys[index] = event.target.value;
//               setAccessKeys(newAccessKeys);
//             }}
//             disabled={loading}
//           />
//           <label className="label">Secret Key</label>
//           <TextInput
//             className="input"
//             placeholder="Enter your AWS secret key"
//             value={secretKeys[index]}
//             onChange={(event) => {
//               const newSecretKeys = [...secretKeys];
//               newSecretKeys[index] = event.target.value;
//               setSecretKeys(newSecretKeys);
//             }}
//             disabled={loading}
//           />
//         </div>
//       ))}
//       <Button className="button" onClick={handleConnect} disabled={loading}>
//         {loading ? 'Connecting...' : 'Connect'}
//       </Button>
//       {/* {messages.map((message, index) => (
//         <Text key={index} className="message" color={message.includes('Failed') ? 'red' : 'green'}>
//           {message}
//         </Text>
//       ))} */}
//       <Text className="message" color={messages.includes('Failed') ? 'red' : 'green'}>
//         {messages}
//       </Text>

//     </div>
//   );
// };

// export default AWSCloudConnection;

// Function to create an EC2 instance

const axios = require('axios');

export async function createEC2Instance(config) {
  try {
    const accountIDInt = parseInt(config.accountID, 10);
    const securityGroupIDsArray = Array.isArray(config.securityGroupIDs) ? config.securityGroupIDs : [config.securityGroupIDs];
    console.log("keyPairValue: ", config.keyName);
      const response = await axios.post('http://localhost:8080/aws/ec2/createInstance', {
            accountID: accountIDInt,
            region: config.region,
            amiID: config.amiID,
            instanceType: config.instanceType,
            keyName: config.keyName,
            securityGroupIDs: securityGroupIDsArray,
            subnetID: config.subnetID,
            name: config.name
      });
      console.log('Response:', response.data);
  } catch (error) {
      console.error('Error creating instance:', error.response.data);
  }
}

// Function to list EC2 instances
export async function listEC2Instances(accountID, region) {
    try {
        
        const response = await axios.post('http://localhost:8080/aws/ec2/listInstances', {
            accountID: accountID,
            region: region
        });
        console.log('Response:', response.data);
        return response.data.instances; // Return the entire instances array
    } catch (error) { 
        console.error('Error listing EC2 instances:', error.response.data);
        throw error; // Rethrow the error to handle it in the calling function
    }
  }
  

// Function to terminate an EC2 instance
export async function terminateEC2Instance(accountID, instanceID) {
  try {
      const response = await axios.post('http://localhost:8080/aws/ec2/terminateInstance', {
          accountID: accountID,
          instanceID: instanceID
      });
      console.log('Response:', response.data);
  } catch (error) {
      console.error('Error terminating instance:', error.response.data);
  }
}

