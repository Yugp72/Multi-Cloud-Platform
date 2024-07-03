'use client';

import { createTheme, Button, Text, Title, notifications } from '@mantine/core';

import { Manrope } from "next/font/google";
const manropeFont = Manrope({ subsets: ["latin"] });

export const theme = createTheme({
  /* Put your mantine theme override here */
  fontFamily: manropeFont.style.fontFamily,
  headings: {
    fontFamily: manropeFont.style.fontFamily,
    sizes: {
      h1: { fontSize: "var(--h1)" },
      h2: { fontSize: "var(--h2)" },
      h3: { fontSize: "var(--h3)" },
      h4: { fontSize: "var(--h4)" },
      h5: { fontSize: "var(--h5)" },
      h6: { fontSize: "var(--h6)" },
    },
  },
  fontSizes: {
    xs: "0.625rem",
    sm: "0.875rem",
    md: "1rem",
    lg: "1.125rem",
    xl: "1.25rem",
    xxl: "1.5rem",
  },
  components: {
    Button: Button.extend({
      defaultProps: {
        color: '#283B4B',
        // variant: 'filled',
        radius: 'xl',
      },
    }),
    Text : Text.extend({
      defaultProps: {
        // color: '#0051A2',
        fw: "500",
      },
    }),
    // notifications: {
    //   position: 'fixed',
    //   bottom: 16,
    //   right: 16,
    // },
  },
});