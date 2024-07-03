"use client";

import { Button, Modal } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import LandingPage from "@/components/index/LandingPage";
import { useRouter, usePathname, useSearchParams } from "next/navigation";
import AuthComp from "@/components/auth/AuthComp";



export default function Home() {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const [opened, { open, close }] = useDisclosure(false);

  const closeAuthModal = () => {
    router.push("/");
  };
  return (
    <main>
      <LandingPage/>

      <Modal
        opened={searchParams.get('authModal') === "login"}
        onClose={closeAuthModal}
        withCloseButton={false}
        size="md"
      >
        <AuthComp/>
      </Modal>
    </main>
  )
}
