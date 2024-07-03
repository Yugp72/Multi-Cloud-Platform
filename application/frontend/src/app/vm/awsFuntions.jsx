// import React from "react";
// import { Table, Dropdown } from "react-bootstrap";
// import { BsThreeDots } from "react-icons/bs";

// const generateTable = ({ headers, data, onRowClick, selectedRowId, onConnectInstance, onTerminateInstance }) => {
//   const cellStyle = {
//     fontSize: "16px", // Adjust the font size as needed
//   };

//   const handleConnect = (car) => {
//     onConnectInstance(car);
//   };

//   const handleTerminate = (car) => {
//     onTerminateInstance(car);
//   };

//   return (
//     <Table striped bordered hover>
//       <thead>
//         <tr>
//           {headers.map((header) => (
//             <th key={header.id} style={cellStyle}>{header.label}</th>
//           ))}
//           <th style={cellStyle}>Actions</th> {/* Add a column for the actions */}
//         </tr>
//       </thead>
//       <tbody>
//         {data.map((row) => (
//           <tr
//             key={row.id}
//             onClick={() => onRowClick(row)}
//             style={{
//               fontSize: "14px",
//               color: "#3a7a11", 
//             }}
//           >
//             {headers.map((header) => (
//               <td key={`${row.id}-${header.id}`}>{row[header.id]}</td>
//             ))}
//             <td>
//               <Dropdown>
//                 <Dropdown.Toggle style={{ cursor: "pointer" }}>
//                   <BsThreeDots />
//                 </Dropdown.Toggle>
//                 <Dropdown.Menu>
//                   <Dropdown.Item onClick={() => handleConnect(row)}>Connect Instance</Dropdown.Item>
//                   <Dropdown.Item onClick={() => handleTerminate(row)}>Terminate Instance</Dropdown.Item>
//                 </Dropdown.Menu>
//               </Dropdown>
//             </td>
//           </tr>
//         ))}
//       </tbody>
//     </Table>
//   );
// };

// export default generateTable;
