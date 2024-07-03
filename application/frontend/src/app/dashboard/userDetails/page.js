"use client"
import Navbar from "@/components/navbar/Navbar"
import {Text, Table, TableData, TextInput, Button, Space, Center} from "@mantine/core"
import { gql, useQuery } from "@apollo/client";
import { useState, useEffect } from "react";

import styles from "@/app/admin/userDetails/page.module.css"


const STUDENTDETAILS=gql`
  query AllStudents {
    allStudents {
      name
      roomNo
      username
    }
  }
`;

// const elements2 = [
//   { place: "UC Cafeteria", lunch: 120, tea: 50 },
//   { place: "Hostel", lunch: 50, tea: 120 },
//   { place: "UC Cafeteria", lunch: 120, tea: 50 },
//   { place: "Hostel", lunch: 50, tea: 120 },
// ];


  

export default function Home() {
  const [allStudents, setAllStudents] = useState([]);
  const check = typeof window !== "undefined" && window.localStorage;
  const token = check ? localStorage.getItem("token") : "";

  const rows = allStudents.map((student) => (
    <Table.Tr key={student.username}>
      <Table.Td>{student.name}</Table.Td>
      <Table.Td>{student.roomNo}</Table.Td>
      <Table.Td>{student.username}</Table.Td>
    </Table.Tr>
  ));

  const { loading: studentsLoading, error: studentsError, data: studentsData } = useQuery(
    STUDENTDETAILS,
    {
      context: {
        headers: {
          authorization: token || "",
        },
      },
    }
  );

  useEffect(() => {
    if (studentsData) {
      console.log(studentsData);
      setAllStudents(studentsData.allStudents);
    }
  }, [studentsData]);

  console.log(allStudents);

  return (
    <main>
      <div className={styles.mainContainer}>
        <Navbar userT="housekeeping" />

        <div className={styles.content}>
          <div className={styles.mContainer}>
            <div className={styles.heading}>
              <Text ta="center" className={styles.headingText}>
                User Details
              </Text>
            </div>
            <div className={styles.subContainer}>
              <div className={styles.table}>
                <Table
                  stickyHeader
                  stickyHeaderOffset={0}
                  horizontalSpacing="xl"
                  highlightOnHover
                  withTableBorder
                  withColumnBorders
                >
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th ta="center" >Name</Table.Th>
                      <Table.Th ta="center">Room No</Table.Th>
                      <Table.Th ta="center">Username</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  <Table.Tbody ta="center">{rows}</Table.Tbody>
                </Table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}