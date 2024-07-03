// AllServerlessPage.jsx
"use client";
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import GenerateTable from "../../components/table/table";
import { listECS, createECS, deleteECS, listEKS, createEKS, deleteEKS, listLambda, createLambda, deleteLambda } from "../../components/connection/AWS/aws_serveless";
import useAllConnectedAccounts from '../getCloudAccount';
import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';
import { Button } from 'react-bootstrap';
import PopupSection from '../../components/popup/page';

const AllServerlessPage = () => {
    const [allECS, setAllECS] = useState([]);
    const [allEKS, setAllEKS] = useState([]);
    const [allLambda, setAllLambda] = useState([]);

    const { data, loading, error } = useAllConnectedAccounts();
    const [isCreateECSPopupOpen, setIsCreateECSPopupOpen] = useState(false);
    const [isCreateEKSPopupOpen, setIsCreateEKSPopupOpen] = useState(false);
    const [isCreateLambdaPopupOpen, setIsCreateLambdaPopupOpen] = useState(false);

    useEffect(() => {
        async function fetchAllServerlessResources() {
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");

                    // Fetch and set ECS clusters
                    const ecsPromises = awsAccounts.map(async (account) => {
                        const ecsResponse = await listECS(account.AccountID, account.Region);
                        return ecsResponse.map((ecs) => ({
                            ...ecs,
                            accountId: account.AccountID,
                            region: account.Region,
                        }));
                    });

                    const allECS = await Promise.all(ecsPromises);
                    const flattenedECS = allECS.flat();
                    setAllECS(flattenedECS);

                    // Fetch and set EKS clusters
                    const eksPromises = awsAccounts.map(async (account) => {
                        const eksResponse = await listEKS();
                        return eksResponse.map((eks) => ({
                            ...eks,
                            accountId: account.AccountID,
                            region: account.Region,
                        }));
                    });

                    const allEKS = await Promise.all(eksPromises);
                    const flattenedEKS = allEKS.flat();
                    setAllEKS(flattenedEKS);

                    // Fetch and set Lambda functions
                    const lambdaPromises = awsAccounts.map(async (account) => {
                        const lambdaResponse = await listLambda();
                        return lambdaResponse.map((lambda) => ({
                            ...lambda,
                            accountId: account.AccountID,
                            region: account.Region,
                        }));
                    });

                    const allLambda = await Promise.all(lambdaPromises);
                    const flattenedLambda = allLambda.flat();
                    setAllLambda(flattenedLambda);
                }
            } catch (error) {
                console.error("Error fetching serverless resources:", error);
            }
        }

        fetchAllServerlessResources();
    }, [data, loading, error]);

    const ecsTableHeaders = [
        { id: 'clusterName', label: 'Cluster Name' },
        { id: 'region', label: 'Region' },
        // Add more headers if needed
    ];

    const eksTableHeaders = [

    ];

    const lambdaTableHeaders = [
        // Define headers for Lambda table
    ];

    const handleCreateECS = async (params) => {
        try {
            console.log('Create ECS params:', params);
            await createECS(params.clusterName, params.region);
            setIsCreateECSPopupOpen(false);
        } catch (error) {
            console.error("Error creating ECS:", error);
        }
    };

    const handleCreateEKS = async (params) => {

    };

    const handleCreateLambda = async (params) => {
        // Implement similar to handleCreateECS for Lambda creation
    };

    const handleDeleteECS = async (ecsId) => {
        try {
            console.log('Delete ECS with ID:', ecsId);
            await deleteECS(ecsId);
        } catch (error) {
            console.error("Error deleting ECS:", error);
        }
    };

    const handleDeleteEKS = async (eksId) => {
        // Implement similar to handleDeleteECS for EKS deletion
    };

    const handleDeleteLambda = async (lambdaId) => {
        // Implement similar to handleDeleteECS for Lambda deletion
    };

    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div>
                        <h2>All ECS Clusters</h2>
                        <Button variant="success" onClick={() => setIsCreateECSPopupOpen(true)}>Create ECS Cluster</Button>
                        <GenerateTable
                            headers={ecsTableHeaders}
                            data={allECS}
                            // Add necessary props similar to AllVPCPage
                            onDelete={handleDeleteECS}
                        />
                    </div>
                    <PopupSection
                        isOpen={isCreateECSPopupOpen}
                        onRequestClose={() => setIsCreateECSPopupOpen(false)}
                        title="Create ECS Cluster"
                        // Add necessary props similar to AllVPCPage
                        onSubmit={handleCreateECS}
                        inputFields={[
                            { label: 'Cluster Name', name: 'clusterName', type: 'text' },
                            { label: 'Region', name: 'region', type: 'text' },
                        ]}
                    />
                    {/* Define similar sections for EKS and Lambda */}
                </div>
            </div>
        </div>
    );
};

export default AllServerlessPage;
