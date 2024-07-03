// AllVPCPage.jsx
"use client";
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import GenerateTable from "../../components/table/table";
import { listVPCs, createVPC, deleteVPC, listSubnets, createSubnet, deleteSubnet, listInternetGateways, createInternetGateway, deleteInternetGateway, detachInternetGateway, attachInternetGateway } from "../../components/connection/AWS/aws_network";
import useAllConnectedAccounts from '../getCloudAccount';
import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';
import { Button } from 'react-bootstrap';
import PopupSection from '../../components/popup/page';

const AllVPCPage = () => {
    const [allVPCs, setAllVPCs] = useState([]);

    const { data, loading, error } = useAllConnectedAccounts();
    const [isCreateVPCPopupOpen, setIsCreateVPCPopupOpen] = useState(false);
    const [allSubnets, setAllSubnets] = useState([]);
    const [isCreateSubnetPopupOpen, setIsCreateSubnetPopupOpen] = useState(false);
    const [isCreateInternetGatewayPopupOpen, setIsCreateInternetGatewayPopupOpen] = useState(false);
    const [allInternetGateways, setAllInternetGateways] = useState([]);
    const [isAttachInternetGatewayPopupOpen, setIsAttachInternetGatewayPopupOpen] = useState(false);


    useEffect(() => {
        async function fetchAllVPCs() {
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");

                    const vpcsPromises = awsAccounts.map(async (account) => {
                        const vpcsResponse = await listVPCs(account.AccountID, account.Region);
                        const vpcs = vpcsResponse.Vpcs.map((vpc, index) => ({
                            ...vpc,
                            id: `${account.AccountID}-${index}`, // Add an id property
                            accountID: account.AccountID,
                            region: account.Region,
                            vpcId: vpc.VpcId,
                            state: vpc.State,
                            cidrBlock: vpc.CidrBlock,
                            dhcpOptionsId: vpc.DhcpOptionsId,
                        }));
                        return vpcs;
                    });

                    const allVpcs = await Promise.all(vpcsPromises);
                    const flattenedVpcs = allVpcs.flat();

                    setAllVPCs(flattenedVpcs);
                }
            } catch (error) {
                console.error("Error fetching all VPCs:", error);
            }
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");

                    const subnetsPromises = awsAccounts.map(async (account) => {
                        const subnetsResponse = await listSubnets(account.AccountID, account.Region);
                        const subnets = subnetsResponse.subnets.map((subnet) => ({
                            id: subnet.id,
                            name: subnet.name,
                            ipv4CIDRBlock: subnet.ipv4CIDRBlock,
                            availabilityZone: subnet.availabilityZone,
                            availabilityZoneID: subnet.availabilityZoneID,
                            availableIPAddressCount: subnet.availableIPAddressCount,
                            state: subnet.state,
                            subnetArn: subnet.subnetArn,
                            vpcID: subnet.vpcID,
                            accountID: account.AccountID,
                        }));
                        return subnets;
                    });


                    const allSubnets = await Promise.all(subnetsPromises);
                    const flattenedSubnets = allSubnets.flat();

                    setAllSubnets(flattenedSubnets);
                }
            } catch (error) {
                console.error("Error fetching all Subnets:", error);
            }
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");

                    const internetGatewaysPromises = awsAccounts.map(async (account) => {
                        const internetGatewaysResponse = await listInternetGateways(account.AccountID, account.Region);
                        const internetGateways = internetGatewaysResponse.map((internetGateway) => ({
                            internetGatewayId: internetGateway.internetGatewayId,
                            id: internetGateway.internetGatewayId,
                            ownerId: internetGateway.ownerId,

                            attachments: internetGateway.attachments ? internetGateway.attachments.map(attachment => ({
                                state: attachment.state,
                                vpcId: attachment.vpcId,
                            })) : [], // Handle case where attachments is null
                            region: account.Region,
                            accountID: account.AccountID,
                            vpcId: (internetGateway.attachments && internetGateway.attachments.length > 0) ? internetGateway.attachments[0].vpcId : null,

                        }));
                        return internetGateways;
                    });

                    const allInternetGateways = await Promise.all(internetGatewaysPromises);
                    const flattenedInternetGateways = allInternetGateways.flat();
                    console.log("flattenedInternetGateways", flattenedInternetGateways);

                    setAllInternetGateways(flattenedInternetGateways);
                }
            } catch (error) {
                console.error("Error fetching all Internet Gateways:", error);
            }

        }

        fetchAllVPCs();
    }, [data, loading, error]);

    const tableHeaders = [
        { id: 'vpcId', label: 'VpcID' },
        { id: 'cidrBlock', label: 'CIDR Block' },
        { id: 'dhcpOptionsId', label: 'DhcpOptionsId' },
        { id: 'region', label: 'Region' },
    ];

    const getVPCDetails = (vpc) => (
        <>
            <p>Details:</p>
            <p>Name: {vpc.name}</p>
            <p>CIDR Block: {vpc.cidrBlock}</p>
            <p>Account ID: {vpc.accountID}</p>
            <p>Region: {vpc.region}</p>
        </>
    );

    const handleCreateVPC = async (params) => {
        try {
            console.log('Create VPC params:', params);
            await createVPC(params);
            setIsCreateVPCPopupOpen(false);
        } catch (error) {
            console.error("Error creating VPC:", error);
        }
    };

    const handleDeleteVPC = async (vpc) => {
        try {
            console.log('Delete VPC:', vpc);
            await deleteVPC(vpc.accountID, vpc.vpcId, vpc.region);
        } catch (error) {
            console.error("Error deleting VPC:", error);
        }
    };

    const createVPCPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Region', name: 'region', type: 'text' },
        { label: 'Name', name: 'name', type: 'text' },
        { label: 'CIDR Block', name: 'cidrBlock', type: 'text' },
    ];

    const deleteVPCButton = ({ row }) => (
        <Button variant="danger" onClick={() => handleDeleteVPC(row)}>Delete VPC</Button>
    );



    const tableHeaders1 = [
        { id: 'id', label: 'Subnet ID' },
        { id: 'ipv4CIDRBlock', label: 'CIDR Block' },
        { id: 'availabilityZone', label: 'Availability Zone' },
        { id: 'state', label: 'Region' },
        { id: 'vpcID', label: 'VPC ID' },
        { id: 'subnetArn', label: 'Subnet Arn' },

    ];

    const getSubnetDetails = (subnet) => (
        <>
            <p>Details:</p>
            <p>Subnet ID: {subnet.id}</p>
            <p>CIDR Block: {subnet.ipv4CIDRBlock}</p>
            <p>State: {subnet.state}</p>
            <p>VpcID: {subnet.vpcID}</p>
        </>
    );

    const handleCreateSubnet = async (params) => {
        try {
            console.log('Create Subnet params:', params);
            await createSubnet(params);
            setIsCreateSubnetPopupOpen(false);
        } catch (error) {
            console.error("Error creating Subnet:", error);
        }
    };

    const handleDeleteSubnet = async (subnet) => {
        try {
            console.log('Delete Subnet:', subnet);
            await deleteSubnet(subnet.accountID, subnet.availabilityZone, subnet.id);
        } catch (error) {
            console.error("Error deleting Subnet:", error);
        }
    };

    const createSubnetPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Region', name: 'region', type: 'text' },
        { label: 'CIDR Block', name: 'cidrBlock', type: 'text' },
        { label: 'VPC ID', name: 'vpcId', type: 'text' },
    ];

    const deleteSubnetButton = ({ row }) => (
        <Button variant="danger" onClick={() => handleDeleteSubnet(row)}>Delete Subnet</Button>
    );

    const handleCreateInternetGateway = async (params) => {
        try {
            console.log('Create Internet Gateway params:', params);
            await createInternetGateway(params);
            setIsCreateInternetGatewayPopupOpen(false);
        } catch (error) {
            console.error("Error creating Internet Gateway:", error);
        }
    };

    const handleDeleteInternetGateway = async (internetGateway) => {
        try {
            console.log('Delete Internet Gateway:', internetGateway);
            await deleteInternetGateway(internetGateway.accountID, internetGateway.id);
        } catch (error) {
            console.error("Error deleting Internet Gateway:", error);
        }
    };

    const internetGatewayTableHeaders = [
        { id: 'internetGatewayId', label: 'Internet Gateway ID' },
        { id: 'vpcId', label: 'VPC ID' },
        { id: 'region', label: 'Region' },
    ];

    const internetGatewayDetails = (internetGateway) => (
        <>
            <p>Details:</p>
            <p>Internet Gateway ID: {internetGateway.id}</p>
            <p>VPC ID: {internetGateway.vpcID}</p>
            <p>Region: {internetGateway.region}</p>
        </>
    );

    const deleteInternetGatewayButton = ({ row }) => (
        <Button variant="danger" onClick={() => handleDeleteInternetGateway(row)}>Delete Internet Gateway</Button>
    );
    const createInternetGatewayPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Region', name: 'region', type: 'text' },
        { label: 'internetGatewayId', name: 'InternetGateway Id', type: 'text' },
    ];

    const attachInternetGatewayButton = ({ row }) => (
        <Button variant="primary" onClick={() => handleAttachInternetGateway(row)}>Attach Internet Gateway</Button>
    );

    const handleAttachInternetGateway = async (params) => {
        console.log("params fomr handleattach: ", params)
        try {
            console.log('Attach Internet Gateway params:', params);
            await attachInternetGateway(params.accountID, params.region, params.vpcId, params.internetGatewayId);
            setIsAttachInternetGatewayPopupOpen(false); s
        } catch (error) {
            console.error("Error attaching Internet Gateway:", error);
        }
    }

    const detachInternetGatewayButton = ({ row }) => (
        <Button variant="primary" onClick={() => handleDetachInternetGateway(row)}>Detach Internet Gateway</Button>
    );

    const handleDetachInternetGateway = async (params) => {
        try {
            console.log('Detach Internet Gateway params:', params);
            await detachInternetGateway(params.accountID, params.region, params.vpcId, params.internetGatewayId);
        } catch (error) {
            console.error("Error detaching Internet Gateway:", error);
        }
    }

    const uploadInternetGatewayPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Region', name: 'region', type: 'text' },
        { label: 'InternetGateway Id', name: 'internetGatewayId', type: 'text' },
        { label: 'Vpc Id', name: 'vpcId', type: 'text' },
    ];



    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div>
                        <h2>All VPCs</h2>
                        <Button variant="success" onClick={() => setIsCreateVPCPopupOpen(true)}>Create VPC</Button>
                        <GenerateTable
                            headers={tableHeaders}
                            data={allVPCs}
                            onRowClick={(row, accountID) => { /* Handle row click */ }}
                            getDetails={getVPCDetails}
                            buttons={[deleteVPCButton]}
                        />
                    </div>
                    <PopupSection
                        isOpen={isCreateVPCPopupOpen}
                        onRequestClose={() => setIsCreateVPCPopupOpen(false)}
                        title="Create VPC"
                        inputFields={createVPCPopupInputFields}
                        onSubmit={handleCreateVPC}
                    />
                </div>
            </div>
            <div className={styles.mainContent}>
                <div className={styles.content}>
                    <div>
                        <h2>All Subnets</h2>
                        <Button variant="success" onClick={() => setIsCreateSubnetPopupOpen(true)}>Create Subnet</Button>
                        <GenerateTable
                            headers={tableHeaders1}
                            data={allSubnets}
                            onRowClick={(row, accountID) => { /* Handle row click */ }}
                            getDetails={getSubnetDetails}
                            buttons={[deleteSubnetButton]}
                        />
                    </div>
                    <PopupSection
                        isOpen={isCreateSubnetPopupOpen}
                        onRequestClose={() => setIsCreateSubnetPopupOpen(false)}
                        title="Create Subnet"
                        inputFields={createSubnetPopupInputFields}
                        onSubmit={handleCreateSubnet}
                    />
                </div>
            </div>
            <div className={styles.mainContent}>
                <div className={styles.content}>
                    <div>
                        <h2>All Internet Gateways</h2>
                        <Button variant="success" onClick={() => setIsCreateInternetGatewayPopupOpen(true)}>Create Internet Gateway</Button>
                        <GenerateTable
                            headers={internetGatewayTableHeaders}
                            data={allInternetGateways}
                            getDetails={internetGatewayDetails}
                            buttons={[deleteInternetGatewayButton, attachInternetGatewayButton, detachInternetGatewayButton]}
                            onRowClick={(set) => {
                                console.log("set", set);
                                setIsAttachInternetGatewayPopupOpen(true);
                            }}
                        />
                    </div>
                    <PopupSection
                        isOpen={isCreateInternetGatewayPopupOpen}
                        onRequestClose={() => setIsCreateInternetGatewayPopupOpen(false)}
                        title="Create Internet Gateway"
                        inputFields={createInternetGatewayPopupInputFields}
                        onSubmit={handleCreateInternetGateway}
                    />
                    <PopupSection
                        isOpen={isAttachInternetGatewayPopupOpen} // Change isCreateInternetGatewayPopupOpen to isAttachInternetGatewayPopupOpen
                        onRequestClose={() => setIsAttachInternetGatewayPopupOpen(false)}
                        title="Attach Internet Gateway"
                        inputFields={uploadInternetGatewayPopupInputFields}
                        onSubmit={handleAttachInternetGateway}
                    />

                </div>
            </div>
        </div>




    );
};


export default AllVPCPage;
