// AllVMsPage.jsx
"use client";
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import GenerateTable from "../../components/table/table";
import { listEC2Instances, terminateEC2Instance, createEC2Instance } from "../../components/connection/AWS/aws_vm";
import { listGCPInstances, terminateGCPInstance, createGCPInstance } from "../../components/connection/GCP/gcp_vm";
import useAllConnectedAccounts from '../getCloudAccount';
import PopupSection from '../../components/popup/page'; // Import the PopupSection component
import { Button } from 'react-bootstrap'; // Import the Button component

import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';
const AccessToken = "ya29.a0AXooCgsjA0BnW7P_ef3YhhTIu_prrizuAHaGqBEGDArwYo_elD25fA4qi69BmvuKt341Zp_kCMGQgvdWf4pLDvdMo5uI6oaOAr503C7GyLMGSCskTKn5UPAOXwPPeLGZLSGVmxUG11BggzluXmVe4Q7S632u8APOZa21aCgYKAXISARISFQHGX2Mi5L8tGcM7uI-7CKU-XMahxw0171";
console.log("Access token:",AccessToken);

const AllVMsPage = () => {
    const [allVMs, setAllVMs] = useState([]);
    const [isPopupOpen, setIsPopupOpen] = useState(false); // State for controlling the popup
    const { data, loading, error } = useAllConnectedAccounts();


    useEffect(() => {
        async function fetchAllVMs() {
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;

                    // Fetch and set AWS instances
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");
                    const awsVMsPromises = awsAccounts.map(async (account) => {
                        const instances = await listEC2Instances(account.AccountID, account.Region);
                        if (instances != null) {
                            instances.forEach((instance, index) => {
                                instance.accountID = account.AccountID;
                                instance.id = `${account.AccountID}-${index}`; // Add an id property
                                instance.cloudProvider = "AWS";
                            });
                        }
                        return instances;
                    });

                    const allAWSInstances = await Promise.all(awsVMsPromises);
                    const flattenedAWSInstances = allAWSInstances.flat();

                    // Fetch and set GCP instances
                    const gcpAccounts = connectedAccounts.filter(account => account.CloudProvider === "GCP");
                    const gcpVMsPromises = gcpAccounts.map(async (account) => {
                        const instances = await listGCPInstances(account.AccountID, AccessToken);
                        if (instances != null) {
                            instances.forEach((instance, index) => {
                                instance.accountID = account.AccountID;
                                instance.id = `${account.AccountID}-${index}`; // Add an id property
                                instance.cloudProvider = "GCP";
                            });
                        }
                        return instances;
                    });


                    const allGCPInstances = await Promise.all(gcpVMsPromises);
                    const flattenedGCPInstances = allGCPInstances.flat();

                    // Merge AWS and GCP instances
                    const allInstances = [...flattenedAWSInstances, ...flattenedGCPInstances];
                    setAllVMs(allInstances);
                }
            } catch (error) {
                console.error("Error fetching all VMs:", error);
            }
        }

        fetchAllVMs();
    }, [data, loading, error]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
    if (!data || !data.cloudAccount) return <div>No data available</div>;

    const tableHeaders = [
        { id: 'id', label: 'ID' },
        { id: 'instance_id', label: 'Instance ID' },
        { id: 'state_code', label: 'Status' },
        { id: 'instance_type', label: 'Instance Type' },
        { id: 'launch_time', label: 'Launch Time' },
        { id: 'availability_zone', label: 'Availability Zone' },
        { id: 'private_ip_address', label: 'Private IP Address' },
        { id: 'public_ip_address', label: 'Public IP Address' },
        { id: 'cloudProvider', label: 'Provider'}
    ];

    const connectInstance = (row) => {
        console.log("Connecting instance:", row);
    };
    const getInstanceDetails = (instance) => (
        <>
            <p>Instance Details:</p>
        </>
    );

    const terminateInstance = (row) => {
        console.log("Terminating instance:", row.accountID, row.instance_id);

        if (!row.accountID || !row.instance_id) {
            console.error('Error: Account ID or Instance ID not provided.');
            return;
        }

        if (row.cloudProvider === "AWS") {
            terminateEC2Instance(row.accountID, row.instance_id)
                .then((response) => {
                    console.log('AWS instance terminated successfully:', response);
                })
                .catch((error) => {
                    console.error('Error terminating AWS instance:', error);
                });
        } else if (row.cloudProvider === "GCP") {
            terminateGCPInstance(row.accountID, row.instance_id)
                .then((response) => {
                    console.log('GCP instance terminated successfully:', response);
                })
                .catch((error) => {
                    console.error('Error terminating GCP instance:', error);
                });
        }

        console.log("Terminating instance:", row);
    };

    // Define buttons for connecting and terminating instances
    const connectButton = ({ row }) => (
        <Button variant="primary" onClick={() => connectInstance(row)}>Connect</Button>
    );

    const terminateButton = ({ row }) => (
        <Button variant="danger" onClick={() => terminateInstance(row)}>Terminate</Button>
    );

    const createAWSInstanceButton = (
        <Button variant="success" onClick={() => setIsPopupOpen(true)}>Create AWS Instance</Button>
    );

    const createGCPInstanceButton = (
        <Button variant='success' onClick={()=> setIsPopupOpen(true)}>Create GCP Instance</Button>
    )

    // Combine the buttons into an array
    const buttons = [connectButton, terminateButton];

    const awsInputFields = [
        { name: 'amiID', label: 'AMI ID', type: 'text' },
        { name: 'instanceType', label: 'Instance Type', type: 'text' },
        { name: 'keyName', label: 'Key Name', type: 'text' },
        { name: 'securityGroupIDs', label: 'Security Group IDs', type: 'text' },
        { name: 'subnetID', label: 'Subnet ID', type: 'text' },
        { name: 'region', label: 'Region', type: 'text' },
        { name: 'accountID', label: 'Account ID', type: 'Integer' },
    ];

    const gcpInputFields = [
        { name: 'machineType', label: 'Machine Type', type: 'text' },
        { name: 'image', label: 'Image', type: 'text' },
        { name: 'zone', label: 'Zone', type: 'text' },
        { name: 'name', label: 'Name', type: 'text' },
        { name: 'accountID', label: 'Account ID', type: 'Integer' },
    ];

    const handleSubmitAWSInstance = async (params) => {
        console.log('Creating AWS instance with params:', params);
        await createEC2Instance(params);
        setIsPopupOpen(false);
    };

    const handleSubmitGCPInstance = async (params) => {
        console.log('Creating GCP instance with params:', params);
        params.token=AccessToken;
        await createGCPInstance(params);
        setIsPopupOpen(false);
    };

    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div className={styles.header}>
                        <h2>All Instances</h2>
                        {createAWSInstanceButton}
                        {createGCPInstanceButton}
                    </div>
                    <GenerateTable
                        headers={tableHeaders}
                        data={allVMs}
                        onRowClick={(row, accountID) => { /* Handle row click */ }}
                        getDetails={getInstanceDetails} // Pass the function to get instance details
                        buttons={buttons} // Pass the array of button components
                    />
                    <PopupSection
                        isOpen={isPopupOpen}
                        onRequestClose={() => setIsPopupOpen(false)}
                        title="Create Instance"
                        inputFields={awsInputFields} // Change to gcpInputFields for GCP instance creation
                        onSubmit={handleSubmitAWSInstance} // Change to handleSubmitGCPInstance for GCP instance creation
                    />

                    <PopupSection
                        isOpen={isPopupOpen}
                        onRequestClose={() => setIsPopupOpen(false)}
                        title="Create GCP Instance"
                        inputFields={gcpInputFields} // Change to gcpInputFields for GCP instance creation
                        onSubmit={handleSubmitGCPInstance} // Change to handleSubmitGCPInstance for GCP instance creation
                    />
                </div>
            </div>
        </div>
    );
};

export default AllVMsPage;
