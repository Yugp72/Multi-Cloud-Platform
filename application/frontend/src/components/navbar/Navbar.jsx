"use client"
import { useState, useEffect } from "react";

import {useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import {Text} from "@mantine/core";
import { jwtDecode } from "jwt-decode";
import styles from "@/components/navbar/Navbar.module.css"

import Logo from "@/assets/general/Logo.svg"
import ProfilePhoto from "@/assets/general/ProfilePhoto.png"



export default function Navbar({userT}){
   
    const check = typeof window !== "undefined" && window.localStorage;
    const token = check ? localStorage.getItem("token") : "";
    const searchParams = useSearchParams();
    const [isDropdownOpen, setDropdownOpen] = useState(false);
    const router = useRouter();
    const handleLogout = () => {
        localStorage.clear();
          router.replace("/");
      };
    const toggleDropdown = () => {
      setDropdownOpen(!isDropdownOpen);
    };
    const type = searchParams.get('userType');
    const user = token ? jwtDecode(token) : null;
    var info;

    useEffect(() => {
      if(user === null){
        router.push("/");
      }
      else{
        info = jwtDecode(token);
        console.log(info);
      }
    },[]);

    // useEffect(() => {
    //     info = jwtDecode(token);
    //     console.log(info);
    // },[token]);


    
    
    return(
      <div className={styles.container}>
            <Image src={Logo} alt="FLH" height={100} width={100}/>
            {userT === "admin" && 
                <div className={styles.menu}>
                    <Link
                        href={{
                        pathname: "/admin",
                        query: { authModal: "registration" },
                        }}>
                        <Text fw={700}> User Registration </Text>
                    </Link>
                    <Link href={"/admin/userDetails"}>
                    <Text fw={700}> User Details </Text>
                    </Link>
                </div>
            }
            {userT === "student" && 
                <div className={styles.menu}>
                    <Link
                        href={{
                        pathname: "/student",
                        query: { authModal: "booking" },
                        }}>
                        <Text fw={700}> Your Booking </Text>
                    </Link>
                    
                </div>
            }

<div className={styles.profileDiv}>
        <div className={styles.profile}>
          <div className={styles.dropdownContainer}>
            <Text fw={700} onClick={toggleDropdown}>
              {user ? user.name : ""}
            </Text>
            {isDropdownOpen && (
              <div className={styles.dropdownContent}>
                
                  <Text fw={700} onClick={handleLogout}> Logout</Text>
              </div>
            )}
          </div>
          <div className={styles.profilePhoto}>
            <Image src={ProfilePhoto} alt="ProfilePhoto" width={40} height={40} />
          </div>
        </div>
      </div>
        </div>
    );
}