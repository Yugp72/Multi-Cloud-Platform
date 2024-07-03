"use client";
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import GenerateTable from "../../components/table/table"; // Import the GenerateTable component
import useAllConnectedAccounts from '../getCloudAccount'; 
import { createGCSBucket, uploadGCSObject, deleteGCSBucket, listGCSBuckets } from '../../components/connection/GCP/gcp_bucket'; // Import GCS bucket functions
import { listS3Buckets, uploadS3Object, createS3Bucket, deleteS3Bucket } from "../../components/connection/AWS/aws_s3"; // Import AWS S3 functions
import { listAzureStorageAccounts, createAzureStorageAccount, deleteAzureStorageAccount } from "../../components/connection/AZURE/azure_bucket";

import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';
import { Button } from 'react-bootstrap'; // Import the Button component
import PopupSection from '../../components/popup/page'; // Import the PopupSection component
import Cookies from 'js-cookie';


const AllBUCKETPage = () => {
    const [allS3Buckets, setAllS3Buckets] = useState([]);
    const [allGCSBuckets, setAllGCSBuckets] = useState([]);
    const [allAzureStorageAccounts, setAllAzureStorageAccounts] = useState([]);
    const { data, loading, error } = useAllConnectedAccounts();
    const [isAWSCreateObjectPopupOpen, setIsAWSCreateObjectPopupOpen] = useState(false);
    const [isGCPCreateObjectPopupOpen, setIsGCPCreateObjectPopupOpen] = useState(false);
    const [isAzureCreateObjectPopupOpen, setIsAzureCreateObjectPopupOpen] = useState(false);
    const [isUploadObjectPopupOpen, setIsUploadObjectPopupOpen] = useState(false);
    const AccessToken = "ya29.a0AXooCgsDJUN2e99FlPd8JvBezrVpL6tuLkVNqL5HC6gPFvGI37uAe0FSfGfnz_lN9Be0ie3PUh6uiYbytFChx7Gfc33JwTJV07s2HWjH_z6L2t9R35p933m6M3mnpuZs5-Iu50rdEx-g5mOKhAubGoIKcZzNUDFfYBvdaCgYKAdkSARISFQHGX2MildCQUJbWD2Oy69WUmmF1kg0171";
    const AccessTokenAZURE = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IkwxS2ZLRklfam5YYndXYzIyeFp4dzFzVUhIMCIsImtpZCI6IkwxS2ZLRklfam5YYndXYzIyeFp4dzFzVUhIMCJ9.eyJhdWQiOiJodHRwczovL21hbmFnZW1lbnQuYXp1cmUuY29tIiwiaXNzIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvODA1OWMzZDUtYTk2Mi00Mzk0LThiNjItZWY3YzkyMTE0MjJhLyIsImlhdCI6MTcxNTAyNDI0MywibmJmIjoxNzE1MDI0MjQzLCJleHAiOjE3MTUwMjkwNjQsImFjciI6IjEiLCJhaW8iOiJBWVFBZS84V0FBQUFmT01OdlBlcUFMamV2eXB0TDN2OE1JL1ZhSHRPWkxjQm9tQnY3NmlCNVRPUGxqZjJGVUFtM0xzZkkrTmUza2RtQXk0Nk5XQXJVL1Z4RVNRbk5XY3NyMmRTZHdFNVRYNE56YnlEeXRXeE5LK1RpYXRiaVYyQjY3ZU9oMlhEK1NncXB4YXlMTWVsODZycjM3MTJ3MDhrZzM0M0xKQWppOGFxVHRrS3pINU1RN0k9IiwiYWx0c2VjaWQiOiIxOmxpdmUuY29tOjAwMDMwMDAwMzQyNTJFRTYiLCJhbXIiOlsicHdkIiwibWZhIl0sImFwcGlkIjoiMzllNjYxNmItNTFiMi00MTkzLWFiYmEtYmRjYzgzOGEwMTU1IiwiYXBwaWRhY3IiOiIxIiwiZW1haWwiOiJ5dWdwNTYyMjlAZ21haWwuY29tIiwiZmFtaWx5X25hbWUiOiJQYXRlbCIsImdpdmVuX25hbWUiOiJZdWciLCJncm91cHMiOlsiNmFmYWRlYzMtYTI1MS00MTc3LWExNDQtYjUwMjIwZmY1OWE4IiwiNDRmZTgxZTctMzU1NS00NjAyLWIzOGYtZDcyNjc5NDJkODM0Il0sImlkcCI6ImxpdmUuY29tIiwiaWR0eXAiOiJ1c2VyIiwiaXBhZGRyIjoiMjQwNToyMDE6MjAxNTo0MDhlOjE4MTc6ZmVmZDpkODQ1OjQ4YWMiLCJuYW1lIjoiWXVnIFBhdGVsIiwib2lkIjoiMzRjMTQ2NzgtOTM2OS00MTlmLTg4MWItNTE2ODAzZjJmYWY0IiwicHVpZCI6IjEwMDMyMDAzNzFEQTlCQzYiLCJyaCI6IjAuQWNZQTFjTlpnR0twbEVPTFl1OThraEZDS2taSWYza0F1dGRQdWtQYXdmajJNQlBHQUFBLiIsInNjcCI6InVzZXJfaW1wZXJzb25hdGlvbiIsInN1YiI6IlY3Wkl4MExaZjV1b2gtMWh5d0xXU05iMUF3ZXFOVTdjbHdZeS12QmU1bTQiLCJ0aWQiOiI4MDU5YzNkNS1hOTYyLTQzOTQtOGI2Mi1lZjdjOTIxMTQyMmEiLCJ1bmlxdWVfbmFtZSI6ImxpdmUuY29tI3l1Z3A1NjIyOUBnbWFpbC5jb20iLCJ1dGkiOiJTRXN6T3V1ZUZFQ1dSS2JrS0JZV0FBIiwidmVyIjoiMS4wIiwid2lkcyI6WyI2MmU5MDM5NC02OWY1LTQyMzctOTE5MC0wMTIxNzcxNDVlMTAiLCIxNThjMDQ3YS1jOTA3LTQ1NTYtYjdlZi00NDY1NTFhNmI1ZjciLCJiNzlmYmY0ZC0zZWY5LTQ2ODktODE0My03NmIxOTRlODU1MDkiXSwieG1zX2Vkb3YiOnRydWUsInhtc190Y2R0IjoxNzEzMTY0OTQxfQ.Jk1IAoFgae8HnjyonGp073R2sUx4x0aeKKgL8tzzqP9PvxbTcW43PxXHRzNpoxblwF59eX8cawwaqh895uV9D956zSlbzRJpkyvrwhVM_56cGBVB2Pw9bkNGpTXMsmdbUEiDV3ay6XcK3QCRNS6i6LpGw0miBbqs9GPJNfV4i3zb2DzCMk3pvyyFXg6qT3kWM0Rja_WGeqKiuPmKIqF5Iprekro8-T-KSczks1Xggb6NHiACniHKunmdlZ-v4gunfNf0j0tuCgzLkH3g5tPHBd6MDxZPX5MiNbMWoYSVpC_xBARmXNqw0ec6NJGnZWEY-NHUNf37H-p7nbHa6GlCUg";    useEffect(() => {
        async function fetchAllBuckets() {
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");
                    const gcpAccounts = connectedAccounts.filter(account => account.CloudProvider === "GCP");
                    const azureAccounts = connectedAccounts.filter(account => account.CloudProvider === "Azure");

                    // Fetch and set AWS S3 buckets
                    const s3BucketsPromises = awsAccounts.map(async (account) => {
                        const buckets = await listS3Buckets(account.AccountID, account.Region);
                        buckets.forEach((bucket, index) => {
                            bucket.id = `${account.AccountID}-${index}`;
                            bucket.accountID = account.AccountID;
                            bucket.region = account.Region;
                            bucket.Provider = "AWS";
                        });
                        return buckets;
                    });
                    const allS3Buckets = await Promise.all(s3BucketsPromises);
                    const flattenedS3Buckets = allS3Buckets.flat();
                    setAllS3Buckets(flattenedS3Buckets);

                    // Fetch and set GCS buckets
                    const gcsBucketsPromises = gcpAccounts.map(async (account) => {
                        const buckets = await listGCSBuckets(account.AccountID, AccessToken);
                        buckets.forEach((bucket, index) => {
                            bucket.id = `${account.AccountID}-${index}`;
                            bucket.accountID = account.AccountID;
                            bucket.region = account.Region;
                            bucket.Provider = "GCP";
                            bucket.date = bucket.created;
                        });
                        return buckets;
                    });
                    const allGCSBuckets = await Promise.all(gcsBucketsPromises);
                    const flattenedGCSBuckets = allGCSBuckets.flat();
                    setAllGCSBuckets(flattenedGCSBuckets);

                    const azureStorageAccountsPromises = azureAccounts.map(async (account) => {
                        const storageAccounts = await listAzureStorageAccounts(account.AccountID, AccessTokenAZURE);
                        storageAccounts.forEach((storageAccount, index) => {
                            storageAccount.id = `${account.AccountID}-${index}`;
                            storageAccount.accountID = account.AccountID;
                            storageAccount.region = account.Region;
                            storageAccount.Provider = "Azure";
                        });
                        return storageAccounts;
                    });
                    console.log("azureStorageAccountsPromises", azureStorageAccountsPromises);
                    const allAzureStorageAccounts = await Promise.all(azureStorageAccountsPromises);
                    const flattenedAzureStorageAccounts = allAzureStorageAccounts.flat();
                    setAllAzureStorageAccounts(flattenedAzureStorageAccounts);
                }
            } catch (error) {
                console.error("Error fetching all buckets:", error);
            }
        }
        fetchAllBuckets();
    }, [data, loading, error])

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
    if (!data || !data.cloudAccount) return <div>No data available</div>;

    const tableHeaders = [
        { id: 'name', label: 'Name' },
        { id: 'date', label: 'Time of Creation' },
        { id: 'Provider', label: 'Cloud Provider' },
    ];

    const getBucketDetails = (bucket) => (
        <>
            <p>Bucket Details:</p>
            <p>Name: {bucket.name}</p>
            <p>Creation Date: {bucket.date}</p>
        </>
    );

    // Dummy functions for button actions
    const uploadObject = (row) => {
        console.log('Upload object:', row);
    };

    const handleAWSDeleteBucket = async (row) => {
        // Determine if it's an AWS S3 or GCS bucket and call the appropriate delete function
            await deleteS3Bucket(row.accountID, row.name,row.region);
        
    };

    const handleGCPDeleteBucket = async (row) => {
        // Determine if it's an AWS S3 or GCS bucket and call the appropriate delete function
         
            await deleteGCSBucket(row.accountID, row.name, AccessToken);
        
    };


    const uploadObjectButton = ({ row }) => (
        <Button variant="success" onClick={() => setAWSIsCreateObjectPopupOpen(true)}>Upload Object</Button>
    );

    const deleteObjectButton = ({ row }) => (
        console.log("Row:",row),
        row.Provider === "AWS" ?
        <Button variant="danger" onClick={() => handleAWSDeleteBucket(row)}>Delete Bucket</Button>
        :
        <Button variant="danger" onClick={() => handleGCPDeleteBucket(row)}>Delete Bucket</Button>
    );

    const buttons = [ uploadObjectButton, deleteObjectButton];

    const uploadObjectPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Region', name: 'region', type: 'text' },
        { label: 'Bucket Name', name: 'bucketName', type: 'text' },
        { label: 'Object Key', name: 'objectKey', type: 'text' },
        { label: 'Content', name: 'content', type: 'file' },
    ];

    const handleUploadObjectSubmit = async (params) => {
        console.log('Upload object params:', params);
        // Determine if it's an AWS S3 or GCS bucket and call the appropriate upload function
        if (params.CloudProvider === "AWS") {
            await uploadS3Object(params);
        } else if (params.CloudProvider === "GCP") {
            await uploadGCSObject(params);
        }
        setIsUploadObjectPopupOpen(false);
    };

    const createAWSBucketPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Region', name: 'region', type: 'text' },
        { label: 'Bucket Name', name: 'bucketName', type: 'text' },
        { label: 'Cloud Provider', name: 'CloudProvider', type: 'text', value: 'AWS', readOnly: true },
    ];

    const createGCSBucketPopupInputFields = [
        { label: 'Account ID', name: 'accountID', type: 'text' },
        { label: 'Bucket Name', name: 'bucketName', type: 'text' },
        { label: 'Cloud Provider', name: 'CloudProvider', type: 'text', value: 'GCP', readOnly: true },
    ];

    const createAZUREStorageAccountPopupInputFields = [
        { label: 'Account Name', name: 'accountName', type: 'text' },
        { label: 'Resource Group', name: 'resourceGroup', type: 'text' },
        { label: 'Subscription ID', name: 'subscriptionID', type: 'text' },
        { label: 'Token', name: 'token', type: 'text' },
        { label: 'Location', name: 'location', type: 'text' },
        { label: 'Storage Type', name: 'storageType', type: 'text' },
        { label: 'Access Tier', name: 'accessTier', type: 'text' },
    ];
    

    const handelCreateBucket = async (params) => {
        console.log('Create bucket params:', params);
        // Determine if it's an AWS S3 or GCS bucket and call the appropriate create function
        if (params.CloudProvider === "AWS") {
            await createS3Bucket(params);
        } else if (params.CloudProvider === "GCP") {
            //add access token
            params.accessToken = AccessToken;
            await createGCSBucket(params);
        }
        else if(params.CloudProvider ==="AZURE"){
            await createAzureStorageAccount(params);
        }
        
        setIsAWSCreateObjectPopupOpen(false);
        setIsGCPCreateObjectPopupOpen(false);
    };


    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div>
                        <h2>All Buckets</h2>
                        <Button variant="success" onClick={() => setIsAWSCreateObjectPopupOpen(true)}>Create AWS Bucket</Button>
                        <Button variant="success" onClick={() => setIsGCPCreateObjectPopupOpen(true)}>Create GCP Bucket</Button>

                        <GenerateTable
                            headers={tableHeaders}
                            data={[...allS3Buckets, ...allGCSBuckets]} // Merge AWS S3 and GCS bucket data
                            onRowClick={(row) => setIsUploadObjectPopupOpen(true)}
                            getDetails={getBucketDetails}
                            buttons={buttons}
                        />
                    </div>
                    <PopupSection
                        isOpen={isUploadObjectPopupOpen}
                        onRequestClose={() => setIsUploadObjectPopupOpen(false)}
                        title="Upload Object"
                        inputFields={uploadObjectPopupInputFields}
                        onSubmit={handleUploadObjectSubmit}
                    />
                    <PopupSection
                        isOpen={isAWSCreateObjectPopupOpen}
                        onRequestClose={() => setIsAWSCreateObjectPopupOpen(false)}
                        title="Create AWS Bucket"
                        inputFields={createAWSBucketPopupInputFields}
                        onSubmit={handelCreateBucket}
                    />
                    <PopupSection
                        isOpen={isGCPCreateObjectPopupOpen}
                        onRequestClose={() => setIsGCPCreateObjectPopupOpen(false)}
                        title="Create GCP Bucket"
                        inputFields={createGCSBucketPopupInputFields}
                        onSubmit={handelCreateBucket}
                    />
                    <PopupSection
                        isOpen={isAzureCreateObjectPopupOpen}
                        onRequestClose={() => setIsAzureCreateObjectPopupOpen(false)}
                        title="Create Azure Storage Account"
                        inputFields={createAZUREStorageAccountPopupInputFields}
                        onSubmit={handelCreateBucket}
                    />
                </div>
            </div>
        </div>
    );
};

export default AllBUCKETPage;
