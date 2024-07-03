import React, { useState } from "react";
import { Table, Button } from "react-bootstrap";
import styles from './table.module.css'; // Import the CSS file

const GenerateTable = ({ headers, data, onRowClick, buttons, getDetails }) => {
    const [expandedRowId, setExpandedRowId] = useState(null);

    const toggleRowExpansion = (rowId) => {
      console.log("rowID:", rowId);
        if (expandedRowId === rowId) {
            setExpandedRowId(null); // When currently expanded row is clicked again, collapse it
        } else {
            setExpandedRowId(rowId); // Else expand this row
        }
    };

    return (
        <div className={styles.tablecontainer}>
            <Table className={styles.table} striped bordered hover>
                <thead>
                    <tr>
                        {headers.map((header) => (
                            <th key={header.id}>{header.label}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {data.map((row) => (
                        <React.Fragment key={row.id}>
                            <tr
                                key={`row-${row.id}`}
                                onClick={() => {
                                    onRowClick(row, row.accountID);
                                    toggleRowExpansion(row.id);
                                }}
                                className={`${styles.cursorpointer} ${expandedRowId === row.id ? styles.expandedrow : ''}`}
                                style={{ marginBottom: '10px' }} // Add margin to the bottom
                            >
                                {headers.map((header) => (
                                    <td key={`${row.id}-${header.id}`}>{row[header.id]}</td>
                                ))}
                            </tr> 
                            {expandedRowId === row.id && (
                                <tr key={`expanded-${row.id}`}>
                                    <td colSpan={headers.length} style={{ padding: '10px' }}> {/* Add padding to td */}
                                        <div className={styles.expandedrowcontent} style={{ padding: '20px' }}> {/* Add padding to div */}
                                            <div>
                                                {getDetails(row)}
                                            </div>
                                            <div className={`${styles.expandedrowcontent} ${styles.expandedrowbuttons}`}> {/* Combine classes */}
                                                {/* Render the buttons */}
                                                {buttons.map((ButtonComponent, index) => (
                                                    <ButtonComponent key={index} row={row} />
                                                ))}
                                            </div>
                                        </div>
                                    </td>
                                </tr>
                            )}
                        </React.Fragment>
                    ))}
                </tbody>
            </Table>
        </div>
    );
};

export default GenerateTable;