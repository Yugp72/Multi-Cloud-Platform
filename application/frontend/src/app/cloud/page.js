"use client";
import React, { useState, useEffect } from 'react';
import Collapsible from 'react-collapsible';
import CustomButton from '@/assets/button/expandButton';
import AWSCloudConnection from '@/components/connection/AWS/Aws_connection/aws';
import AzureCloudConnection from '@/components/connection/AZURE/azure';
import GCPConnection from '@/components/connection/GCP/gcp';
import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';
import styles from './page.module.css'; 
import useAllConnectedAccounts from '../getCloudAccount';

const CloudConnectionPage = () => {
    const [isAWSCollapsed, setIsAWSCollapsed] = useState(true);
    const [isAzureCollapsed, setIsAzureCollapsed] = useState(true);
    const [isGCPCollapsed, setIsGCPCollapsed] = useState(true);
    const [connectedAccounts, setConnectedAccounts] = useState([]);
    const { data, loading, error } = useAllConnectedAccounts();

    useEffect(() => {
        if (!loading && !error && data) {
            setConnectedAccounts(data.cloudAccount);
            if (data.cloudAccount) {
                data.cloudAccount.forEach(account => {
                    if (account.CloudProvider === "GCP") {
                        handleLoginAndGetToken(account.AccountID);
                        console.log("if gcp: ",account.AccountID);
                    }
                    else if(account.CloudProvider === "AZURE"){
                        handleLoginAndGetAZUREToken(account.AccountID);
                        console.log("if azure: ",account.AccountID);
                    }
                });
            }
            
        }
       

    }, [data, error, loading]);

    const toggleAWS = () => {
        setIsAWSCollapsed(!isAWSCollapsed);
    };

    const toggleAzure = () => {
        setIsAzureCollapsed(!isAzureCollapsed);
    };

    const toggleGCP = (data) => {
        setIsGCPCollapsed(!isGCPCollapsed); 
    };
    

    const handleLoginAndGetToken = async (accountID) => {
        try {
            const requestBody = {
                accountID: accountID
            };
    
            const requestOptions = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestBody)
            };
    
            const loginResponse = await fetch('http://localhost:8080/auth/google/login', requestOptions);
            console.log("Jay responce: ",loginResponse);
            if (!loginResponse.ok) {
                throw new Error('Failed to initiate login');
            }
    
            // Redirect the user to Google for authentication
            window.location.href = loginResponse.url;
        } catch (error) {
            console.error('Error initiating Google login:', error);
        }
    };

    const handleLoginAndGetAZUREToken = async (accountID) => {
        try {
            const requestBody = {
                accountID: accountID
            };
    
            const requestOptions = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestBody)
            };
    
            const loginResponse = await fetch('http://localhost:8080/auth/azure/login', requestOptions);
            console.log("Jay responce: ",loginResponse);
            if (!loginResponse.ok) {
                throw new Error('Failed to initiate login');
            }
    
            // Redirect the user to Google for authentication
            window.location.href = loginResponse.url;
        } catch (error) {
            console.error('Error initiating Google login:', error);
        }
    }

    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div className={styles.leftContainer}>
                        <h1>Cloud Connection</h1>
                        <div className={styles.collapsibleContainer}>
                            {/* AWS section */}
                            <CustomButton
                                name="AWS"
                                providerSymbol=""
                                expandSymbol=""
                                onClick={toggleAWS}
                            />
                            <Collapsible open={!isAWSCollapsed} transitionTime={150} trigger="">
                                <AWSCloudConnection />
                            </Collapsible>

                            {/* Azure section */}
                            <CustomButton
                                name="Azure"
                                providerSymbol=""
                                expandSymbol=""
                                onClick={toggleAzure}
                            />
                            <Collapsible open={!isAzureCollapsed} transitionTime={150} trigger="">
                                <AzureCloudConnection />
                            </Collapsible>

                            {/* GCP section */}
                            <CustomButton
                                name="GCP"
                                providerSymbol=""
                                expandSymbol=""
                                onClick={toggleGCP}
                            />
                            <Collapsible open={!isGCPCollapsed} transitionTime={150} trigger="">
                                <GCPConnection />
                            </Collapsible>
                        </div>
                    </div>
                    <div className={styles.rightContainer}>
                        <h2>Connected Accounts</h2>
                        <table>
                            <thead>
                                <tr>
                                    <th>Provider</th>
                                    <th>Account ID</th>
                                    <th>Account Region / Other Properties</th>
                                </tr>
                            </thead>
                            <tbody>
                                {/* Render connected accounts */}
                                {connectedAccounts.map((account, index) => (
                                    <tr key={index}>
                                        <td>{account.CloudProvider}</td>
                                        <td>{account.AccountID}</td>
                                        <td>{account.Region || account.SubscriptionID || account.ProjectID}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default CloudConnectionPage;
