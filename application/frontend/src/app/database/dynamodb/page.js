// AllDynamoDBTablesPage.jsx
"use client";
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import { listDynamoDBItems, listDynamoDBTables, deleteDynamoDBItem, deleteDynamoDBTable , createDynamoDBTable} from "../../../components/connection/AWS/aws_dynamodb";
import useAllConnectedAccounts from '../../getCloudAccount';
import Navbar from '@/components/navbar/Navbar';
import Sidebar from '@/components/navbar/Sidebar';
import GenerateTable from "../../../components/table/table";
import PopupSection from '../../../components/popup/page'; 

import { Button, Modal } from 'react-bootstrap';


const AllDynamoDBTablesPage = () => {
    const [allTables, setAllTables] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [isPopupOpen, setIsPopupOpen] = useState(false); // State for controlling the popup

    const [tableItems, setTableItems] = useState([]);
    const [selectedTableName, setSelectedTableName] = useState('');
    const { data, loading, error } = useAllConnectedAccounts();

    useEffect(() => {
        async function fetchAllTablesAndItems() {
            try {
                if (!loading && !error && data && data.cloudAccount) {
                    const connectedAccounts = data.cloudAccount;
                    const awsAccounts = connectedAccounts.filter(account => account.CloudProvider === "AWS");
                    const allDynamoDBTables = [];
                    const allTableItems = []; // Changed variable name to avoid conflict

                    for (const account of awsAccounts) {
                        let tablesResponse = await listDynamoDBTables({
                            accountID: account.AccountID,
                            region: account.Region
                        });
                        console.log("tablesResponse: ", tablesResponse);

                        const tables = tablesResponse.map(table => ({
                            id: `${account.AccountID}-${table.TableName}`,
                            Name: table.TableName,
                            Region: account.Region,
                            AccountName: account.name,
                            CreationDateTime: table.CreationDateTime,
                            TableStatus: table.TableStatus,
                            Attribute: table.AttributeDefinitions,
                            AccountID: account.AccountID,
                            TableId: table.TableId,

                        }));
                        

                        const firstTableItems = await listDynamoDBItems({
                            accountID: account.AccountID,
                            region: account.Region,
                            tableName: tables[0].Name
                        });
                        


                        const tableItemsObj = {
                            tableName: tables[0].Name,
                            items: firstTableItems
                        };
                        console.log("tableItemObj: ", tableItemsObj);

                        allDynamoDBTables.push(...tables);
                        allTableItems.push(tableItemsObj); // Push to allTableItems
                    }

                    setAllTables(allDynamoDBTables);
                    setTableItems(allTableItems); // Set tableItems
                }
            } catch (error) {
                console.error("Error fetching all DynamoDB tables and items:", error);
            }
        }

        fetchAllTablesAndItems();
    }, [data, loading, error]);

    const handleCloseModal = () => {
        setShowModal(false);
        setTableItems([]);
        setSelectedTableName('');
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
    if (!data || !data.cloudAccount) return <div>No data available</div>;
    const ItemDetails = ({ item, index }) => {
        console.log("item: ",item);
        const renderItemDetails = (item) => {
            return (
                <ul>
                    {Object.entries(item).map(([key, value]) => (
                        <li key={key}>
                            <strong>{key}: </strong>
                            {value === null ? 'null' : JSON.stringify(value)}
                        </li>
                    ))}
                </ul>
            );
        };
    
        return (
            <div>
                <p>Item {index + 1}:</p>
                {renderItemDetails(item)}
            </div>
        );
    };

    const deleteTable = (row) => { 
        console.log("Deleting Table:", row);
        
        if (!row.AccountID || !row.Name) {
            console.error('Error: Account ID or Instance ID not provided.');
            return;
        }

        deleteDynamoDBTable(row.AccountID, row.Name, row.Region) // Pass instance_id to terminateEC2Instance
            .then((response) => {
                console.log('Table terminated successfully:', response);
            })
            .catch((error) => {
                console.error('Error terminating Table:', error);
            });

        console.log("Terminating Table:", row);
    };
    const tableHeaders = [
        { id: 'Name', label: 'Name' },
        { id: 'Region', label: 'Region' },
        { id: 'CreationDateTime', label: 'Creation Date' },
        { id: 'TableStatus', label: 'Status' },
        { id: 'Actions', label: 'Actions' },
        { id: 'TableId', label: 'Table ID' }
    ];
    const getDetails = (tableItem) => (
        <>
            <p>Table Details:</p>
            <p>Name: {tableItem.tableName}</p>
            <p>Creation Date: {tableItem.creationDate}</p>
            {/* Add more details as needed */}
        </>
    );
    const createTableButton = (
        <Button variant="success" onClick={() => setIsPopupOpen(true)}>Create Instance</Button>
    );

    const buttons = [
        ({ row }) => (
            <Button variant="success" onClick={() => console.log('Edit row:', row)}>Edit</Button>
        ),
        ({ row }) => (
            <Button variant="danger" onClick={() => deleteTable(row)}>Delete Table</Button>
        ),
    ];
    const handleSubmitEC2Instance = async (params) => {
        console.log('Creating EC2 instance with params:', params);
        await createDynamoDBTable(params);
        setIsPopupOpen(false);
    };
    const inputFields = [
        { name: 'tableName', label: 'Table Name', type: 'text'},
        { name: 'region', label: 'Region', type: 'text' },
        { name: 'accountID', label: 'Account ID', type: 'Integer'},
        { name: 'attributeName', label: 'Attribute Name', type: 'text' },
        { name: 'attributeType', label: 'Attribute Type', type: 'text' },
        { name: 'keySchema', label: 'Key Schema', type: 'text' },
        { name: 'provisionedThroughput', label: 'Provisioned Throughput', type: 'text' }, 
    ]

    return (
        <div className={styles.pageContainer}>
            <Sidebar />
            <div className={styles.mainContent}>
                <Navbar />
                <div className={styles.content}>
                    <div>
                        <h2>All DynamoDB Tables</h2>
                        {createTableButton}
                        <GenerateTable
                            headers={tableHeaders}
                            data={allTables}
                            buttons={buttons}
                            onRowClick={(row) => {
                                console.log(row);
                            }}
                            getDetails={getDetails}
                        />
                        <PopupSection
                        isOpen={isPopupOpen}
                        onRequestClose={() => setIsPopupOpen(false)}
                        title="Create Table Successfully"
                        inputFields={inputFields}
                        onSubmit={handleSubmitEC2Instance}
                    />
                    </div>
                </div>
            </div>
            <div style={{marginLeft: '40%'}}className={styles.content}>
            <Modal show={showModal} onHide={handleCloseModal}>
                <Modal.Body>
                    <ul>
                        {tableItems
                            .filter(tableItem => tableItem.tableName === selectedTableName) // Filter tableItems
                            .map((tableItem, index) => (
                                <li key={index}>
                                    <p>Table: {tableItem.tableName}</p>
                                    {tableItem.items.map((item, itemIndex) => (
                                        <ItemDetails key={itemIndex} item={item} index={itemIndex} />
                                    ))}
                                </li>
                            ))}
                    </ul>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={handleCloseModal}>Close</Button>
                </Modal.Footer>
            </Modal>
            </div>

        </div>
    );
};

export default AllDynamoDBTablesPage;
