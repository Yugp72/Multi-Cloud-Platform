import { useState } from "react";
import { useRouter } from "next/navigation";
import { Text } from "@mantine/core";
import styles from "./Sidebar.module.css";

export default function Sidebar() {
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [showDatabaseSubNav, setShowDatabaseSubNav] = useState(false);

  const handleNavigation = (path) => {
    router.push(path);
    setIsOpen(false);
  };

  return (
    <div className={
      isOpen ? `${styles.sidebar} ${styles.open}` : `${styles.sidebar} ${styles.close}`}> 
        
      <button onClick={() => setIsOpen(!isOpen)} className={styles.toggle}>
        {isOpen ? "Close" : "Open"}
      </button>

      <div className={styles.links}>
        <Text onClick={() => handleNavigation("/vm")} className={styles.link}>
          Virtual Machines
        </Text>
        <Text 
          onClick={() => {
            handleNavigation("/database");
            setShowDatabaseSubNav(false);
          }} 
          className={styles.link}
          onMouseEnter={() => setShowDatabaseSubNav(true)}
          onMouseLeave={() => setShowDatabaseSubNav(false)}
        >
          Database
          {showDatabaseSubNav && (
            <div className={styles.subNav}>
              <Text onClick={() => handleNavigation("/database/dynamodb")} className={styles.subNavLink}>
                AWS Dynamodb
              </Text>
              <Text onClick={() => handleNavigation("/database/firebase")} className={styles.subNavLink}>
                GCP Firebase
              </Text>
              <Text onClick={() => handleNavigation("/database/cosmosdb")} className={styles.subNavLink}>
                Azure Cosmosdb
              </Text>
            </div>
          )}
        </Text>
        <Text onClick={() => handleNavigation("/bucket")} className={styles.link}>
          Storage
        </Text>
        <Text onClick={() => handleNavigation("/serverless")} className={styles.link}>
          Serverless
        </Text>
        <Text onClick={() => handleNavigation("/network")} className={styles.link}>
          Network
        </Text>
        <Text onClick={() => handleNavigation("/profile")} className={styles.link}>
          Profile
        </Text>
        <Text onClick={() => handleNavigation("/settings")} className={styles.link}>
          Settings
        </Text>
      </div>
    </div>
  );
}
