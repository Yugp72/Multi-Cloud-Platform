'use client'
import React from "react";

import Image from 'next/image'
import styles from './LandigPage.module.css'
import Link from 'next/link'
import { Button, Title, Text, Modal } from "@mantine/core";
import { useRouter } from "next/router";
import ImageStackComponent from "@/components/helper/ImageStackComponenet";

import Logo from '@/assets/general/Logo.svg'
import Hero from '@/assets/general/Hero.svg'
import Poster1 from '@/assets/general/Poster1.svg'
import Book from '@/assets/general/Book.svg'
import Colab from '@/assets/general/AHDUNIxFLH.png'
import Linkedin from "@/assets/footer/icons/Linkedin.svg";
import Facebook from "@/assets/footer/icons/Facebook.svg";
import Instagram from "@/assets/footer/icons/Instagram.svg";
import LinkedinH from "@/assets/footer/icons/LinkedinH.svg";
import FacebookH from "@/assets/footer/icons/FacebookH.svg";
import InstagramH from "@/assets/footer/icons/InstagramH.svg";

export const CONTACT_LINKS = {
  Facebook: "https://www.facebook.com/",
  Linkedin: "https://www.linkedin.com/",
  Instagram: "https://www.instagram.com/",
};

const images = {
  Linkedin,
  Facebook,
  Instagram,
  LinkedinH,
  FacebookH,
  InstagramH,
};

const Company = ["Who We Are", "Careers", "Team", "Report Fraud"];
const CompanyRoutes = ["", "", "", ""];
const Legal = [
  "Terms & Conditions",
  "Refund & Cancellation",
  "Privacy Policy",
  "Cookie Policy",
  "Offer Terms",
];
const LegalRoutes = ["", "", "", "", ""];
const Support = ["abc@gmail.com", "(+91) 1234567890", "(+91) 1234567891"];
const SupportHrefs = [
  "mailto:help@abc.xyz",
  "tel:+911234567890",
  "tel:+911234567891",
];
const SocialMediaIcons = ["Linkedin", "Facebook", "Instagram"];

const CompanyList = Company.map((General, index) => {
  return (
    <li key={index}>
      <Link href={`${CompanyRoutes[index]}`}>{General}</Link>
    </li>
  );
});
const LegalList = Legal.map((Browse, index) => {
  return (
    <li key={index}>
      <Link href={`${LegalRoutes[index]}`}>{Browse}</Link>
    </li>
  );
});
const SupportList = Support.map((Support, index) => {
  return (
    <li key={index}>
      <a href={`${SupportHrefs[index]}`}>{Support}</a>
    </li>
  );
});
const SocialMediaIconsList = SocialMediaIcons.map(
  (socialMediaIconName, index) => {
    return (
      <ImageStackComponent
        key={index}
        link={CONTACT_LINKS[socialMediaIconName]}
        normalDisplay={images[socialMediaIconName]}
        hoverDisplay={
          images[(socialMediaIconName + "H")]
        }
        iconsClass={styles.Icons}
        iconsWrapperClass={styles.IconsSubWrapper}
        StackImageStyle={{}}
      />
    );
  }
);


export default function LandingPage() {
  const footerCompanyLinksRef = React.useRef(12);
  const footerLegalLinksRef = React.useRef(123);
  const footerSupportLinksRef = React.useRef(1234);
  return (
    <main>
      <div>
        <div className={styles.container1}>
          <div className={styles.navbar}>
            <Link href="/">
              <Image src={Logo} height={40} width={40} alt="FLH"/>
            </Link>
            <div className={styles.loginDiv}>
                <Button className={styles.loginButton} 
                component={Link}
                href={{
                  pathname: "/",
                  query: { authModal: "login" },
                }}>
                  Login
                </Button>
            </div>
          </div>


          <div className={styles.heroContainer}>
            <div className={styles.heroContent}>
              <div className={styles.content}>
                <div className={styles.contentText}>
                <Title order={1} className={styles.heroTitle}>
                Lorem ipsum dolor sit amet consectetur.
                                </Title>
                <Text
                  className={styles.heroDesc}
                >
                  Lorem ipsum dolor sit amet consectetur. Nisl ac tortor diam
                  gravida diam. Sem egestas cras ultricies massa morbi erat in. Nec
                  habitasse a et ut duis dignissim. In dui viverra pulvinar magna
                  nunc urna sed egestas. In ut aliquam netus.
                </Text>
                </div>
                <div>
                  <Button>
                    See Photos
                  </Button>
                </div>
              </div>
            </div>
            <div className={styles.heroChar}>
              <Image src={Hero} alt="FLH" />
            </div>
          </div>
        </div>

        <div className={styles.mainContainer2}>
          <div className={styles.container2}>
            <div className={styles.part1}>
              <div className={styles.poster1}>
                <Image src={Poster1} alt="Poster1"/>
              </div>
              <div className={styles.poster1Text}>
                <Title size="h3" className={styles.text1Title}>
                  Comfy and Comprehensive Hostel Living
                </Title>
                <Text className={styles.text1Desc}>
                  Anything and Everything that makes a comfortable home at your disposal, to make your stay fun and fulfilling.
                </Text>
              </div>
            </div>
            <div className={styles.part2}>
              <div className={styles.poster2Text}>
                <Title size="h3" className={styles.text2Title}>
                  Engrossing Library Space
                </Title>
                <Text className={styles.text2Desc}>
                Flick through your favourite page turners and international bestsellers in a space mode for solitude and indulgence.
                </Text>
              </div>
              <div className={styles.book}>
                <Image src={Book} alt="Book"/>
              </div>
            </div>
          </div>
        </div>

        <div className={styles.container3}>
          <Title size="h4" c="white" className={styles.colabText}>
            WE ARE NOW THE OFFICIAL STUDENT HOUSING PARTNER OF AHMEDABAD UNIVERSITY.
          </Title>
          <Image src={Colab} alt="FLHxAHDUNI"/>
        </div>

        <footer className={styles.Wrapper}>
          <div className={styles.UpperContainer}>
            <div className={styles.Links}>
              <div ref={footerCompanyLinksRef}>
                <span>
                  Company{" "}
                  {/* <Image src={arrorwImg} className={styles.Arrow} alt="arrow" /> */}
                </span>
                <ul>{CompanyList}</ul>
              </div>
              <div ref={footerCompanyLinksRef}>
                <span>
                  Legal{" "}
                  {/* <Image src={arrorwImg} className={styles.Arrow} alt="arrow" /> */}
                </span>
                <ul>{LegalList}</ul>
              </div>
              <div ref={footerCompanyLinksRef}>
                <span>
                  Support{" "}
                  {/* <Image src={arrorwImg} className={styles.Arrow} alt="" /> */}
                </span>
                <ul>{SupportList}</ul>
              </div>
            </div>
          </div>

          <div className={styles.LowerContainer}>
            <div className={styles.CompanyLogoWrapper}>
              <Link target="_blank" rel="noopener noreferrer" href="/">
                <Image src={Logo} alt="Logo" />
              </Link>
              <div className={styles.Copyright}><Text>Copyright (c) 2023</Text></div>
            </div>
            <div className={styles.IconsWrapper}>
              <div
                style={{
                  display: "flex",
                }}
              >
                {SocialMediaIconsList}
              </div>
            </div>
          </div>
        </footer>
      </div>
    </main>
  )
}
