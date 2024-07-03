import { useState } from "react";
import { useRouter } from "next/navigation";
import { Text, Button, Space, UnstyledButton } from "@mantine/core";
import { gql, useQuery, useMutation } from "@apollo/client";
import { IconTrash } from "@tabler/icons";
import Navbar from "@/components/navbar/Navbar";
import Table from "@/components/table/Table";
import LoadingOverlay from "@/components/loading-overlay/LoadingOverlay";
import RegComp from "@/components/auth/RegComp";
import styles from "@/app/cloud/page.module.css";

//make a dahsboard where all data of different clouds stays like  billing, resources like storage, ec2, etc





