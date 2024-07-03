"use client";
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import { listBuckets } from "../../../components/connection/AWS/aws_s3";
import useAllConnectedAccounts from '../../getCloudAccount'; 

import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';

const AllBUCKETPage = () => {
    const [allBuckets, setAllBuckets] = useState([]);
    const { data, loading, error } = useAllConnectedAccounts();

    useEffect(() => {
        async function fetchAllBuckets() {
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");
                    const credentials = {
                        "accountID": 2,
                        "region": "ap-south-1",
                    }
                    const s3BucketsPromises = awsAccounts.map(async (account) => {
                        const buckets = await listBuckets(credentials);
                        return buckets.map(bucket => ({
                            ...bucket,
                            accountName: account.name,
                        }));
                    });

                    const allS3Buckets = await Promise.all(s3BucketsPromises);
                    const flattenedS3Buckets = allS3Buckets.flat();
                    
                    setAllBuckets(flattenedS3Buckets);
                }
            } catch (error) {
                console.error("Error fetching all S3 buckets:", error);
            }
        }
        
        fetchAllBuckets();
    }, [data, loading, error]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
    if (!data || !data.cloudAccount) return <div>No data available</div>;

    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div>
                        <h2>All S3 Buckets</h2>
                        <table>
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Region</th>
                                    <th>Account Name</th>
                                </tr>
                            </thead>
                            <tbody>
                                {allBuckets.map((bucket, index) => (
                                    <tr key={index}>
                                        <td>{bucket.Name}</td>
                                        <td>{bucket.Region}</td>
                                        <td>{bucket.accountName}</td>
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

export default AllBUCKETPage;
