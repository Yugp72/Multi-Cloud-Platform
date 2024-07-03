import { Manrope } from "next/font/google";
import './globals.css'
import '@mantine/core/styles.css';
import { theme } from '../../theme';
import {
  ApolloProvider,
} from "@apollo/client";

import { MantineProvider, ColorSchemeScript } from '@mantine/core';
import { Notifications, notifications } from '@mantine/notifications';

const manropeFont = Manrope({ subsets: ["latin"] });

import Apollo  from "@/app/Apollo";

export const metadata = {
  title: 'FLH',
  description: '',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
      <ColorSchemeScript />
      <link rel="icon" href="/favicon.ico" sizes="any" />
      </head>
      <body className={manropeFont.className}>
          <MantineProvider theme={theme}>
              <Apollo children={children} />
              <Notifications position="top-right" styles={{
                notification: {
                  width: 300,
                  maxHeight: 100,
                  position: 'fixed',
                  bottom: 16,
                  right: 16,
                  zIndex: 9999,
                },
              }}/>
          </MantineProvider>
      </body>
    </html>
  )
}

// const crypto = require('crypto');
// const express = require('express');
// const axios = require('axios');
// const querystring = require('querystring');
// const cookieParser = require('cookie-parser');
// import useAllConnectedAccounts from './getCloudAccount';

// const app = express();
// app.use(cookieParser());

// const REDIRECT_URI = 'http://localhost:8081/auth/azure/callback';
// const AUTH_URL = 'https://login.microsoftonline.com/common/oauth2/authorize';
// const TOKEN_URL = 'https://login.microsoftonline.com/common/oauth2/token';
// const SCOPE = 'https://management.azure.com/.default';

// function base64URLEncode(buffer) {
//   return buffer.toString('base64')
//       .replace(/\+/g, '-')
//       .replace(/\//g, '_')
//       .replace(/=/g, '');
// }

// function generateCodeVerifier() {
//   return base64URLEncode(crypto.randomBytes(32));
// }

// function sha256(buffer) {
//   return crypto.createHash('sha256').update(buffer).digest();
// }

// function generateCodeChallenge(verifier) {
//   return base64URLEncode(sha256(verifier));
// }

// const generateRandomString = (length) => {
//   return crypto.randomBytes(length).toString('hex');
// };

// app.get('/auth/azure/login', async (req, res) => {
//   try {
//     // Fetch Azure account data
//     console.log("Fetching Azure account data");
//     // const { data, loading, error } = await useAllConnectedAccounts();
//     // console.log("Data fetched: ",data);

//     // Assuming the data structure has an array of accounts
//     // const azureAccounts = data.cloudAccount.filter(account => account.CloudProvider === 'Azure');

//     // Assuming we're using the first Azure account found
//     // const azureClientId = azureAccounts[0]?.ClientId;
//     const azureClientId = '39e6616b-51b2-4193-abba-bdcc838a0155'

//     if (!azureClientId) {
//       throw new Error('Azure client ID not found');
//     }

//     const codeVerifier = generateCodeVerifier();
//     const codeChallenge = generateCodeChallenge(codeVerifier);
//     const state = generateRandomString(16);

//     // Store the code verifier in the client's browser local storage
//     res.cookie('codeVerifier', codeVerifier, { httpOnly: true });
   
//     const queryParams = querystring.stringify({
//       client_id: azureClientId, // Use the Azure client ID
//       response_type: 'code',
//       redirect_uri: REDIRECT_URI,
//       scope: SCOPE,
//       code_challenge: "6WoFF5hEW7FFF2uGmRtLzvs3BxbhgzVlQhMI2DsISmk",
//       code_challenge_method: 'S256',
//     });

//     const authUrl = `${AUTH_URL}?${queryParams}`;
//     res.redirect(authUrl);
//   } catch (error) {
//     console.error('Error fetching Azure account data:', error);
//     res.status(500).send('Failed to fetch Azure account data');
//   }
// });
// app.get('/auth/azure/callback', async (req, res) => {
//   const code = req.query.code;
//   const codeVerifier = req.query.codeVerifier || req.cookies.codeVerifier; 

//   console.log('Code:', code);
//   console.log('Code Verifier:', codeVerifier);
//   const azureClientId = '39e6616b-51b2-4193-abba-bdcc838a0155'

//   const tokenParams = {
//     client_id: azureClientId, // Use the Azure client ID
//     code: code,
//     redirect_uri: REDIRECT_URI,
//     grant_type: 'authorization_code',
//     code_verifier: 'tYPEcO1sjzXFae8lZJRGsxQsQBOI8nkDX3-ses_u_E0ZVh9vETTKZXdKGY-ZWj_FBgv4ExEwriri9QjKOuNMWDEfWGj3FJPP2sA-E8qTusRTTaG2N4O-8a4Gwq_MqTw3',
//     scope: SCOPE
//   };
//   console.log(tokenParams)

//   try {
//     const response = await axios({
//       method: 'get',
//       url: TOKEN_URL,
//       headers: {  
//         'Content-Type': 'application/x-www-form-urlencoded'
//       },
//       data: qs.stringify(tokenParams)
//     });
//     const token = response.data.access_token;
//     console.log('Access Token:', response);
//     res.send(`Access Token: ${token}`);
//     res.redirect('http://localhost:3000');
//   } catch (error) {
//     console.error('Failed to exchange code for token:');
//     res.status(500).send('Failed to exchange authorization code for token');
//   }
// });




// app.listen(8081, () => {
//   console.log('Express server started on port 8081');
// });


// const crypto = require('crypto');
// const base64url = require('base64url');
// const express = require('express');
// const { default: axios } = require('axios');
// const querystring = require('querystring');
// const cookieParser = require('cookie-parser');

// const app = express();
// app.use(cookieParser());

// const CLIENT_ID = "10a62c80-0d4f-4cc1-89c8-667cfb4c3831";
// const CLIENT_SECRET = "Nv58Q~F.ZBbYt9YnUnLBTmgE.U5IPu4.0cYplaJB";
// const REDIRECT_URI = 'http://localhost:8081/auth/azure/callback';
// const AUTH_URL = 'https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/authorize';
// const TOKEN_URL = 'https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/token';

// const SCOPE = 'https://management.azure.com/.default';

// function base64URLEncode(buffer) {
//   return buffer.toString('base64')
//       .replace(/\+/g, '-')
//       .replace(/\//g, '_')
//       .replace(/=/g, '');
// }

// function generateCodeVerifier() {
//   return base64URLEncode(crypto.randomBytes(32));
// }

// const generateRandomString = (length) => {
//   return crypto.randomBytes(length).toString('hex');
// };

// function generateCodeChallenge(verifier) {
//   return base64URLEncode(sha256(verifier));
// }

// function sha256(buffer) {
//   return crypto.createHash('sha256').update(buffer).digest();
// }

// app.get('/auth/azure/login', (req, res) => {
//   const codeVerifier = generateCodeVerifier();
//   const codeChallenge = generateCodeChallenge(codeVerifier);
//   const state = generateRandomString(16);

//   // Store the code verifier in the client's browser local storage
//   res.cookie('codeVerifier', codeVerifier, { httpOnly: true });
 
//   const queryParams = querystring.stringify({
//     client_id: CLIENT_ID,
//     response_type: 'code',
//     redirect_uri: REDIRECT_URI,
//     scope: SCOPE,
//     code_challenge: codeChallenge,
//     code_challenge_method: 'S256',
//     state: state,
//   });

//   const authUrl = `${AUTH_URL}?${queryParams}`;
//   res.redirect(authUrl);
// });

// app.get('/auth/azure/callback', async (req, res) => {
//   const code = req.query.code;
//   const state = req.query.state; // You may want to validate the state parameter here
  
//   const codeVerifier = req.query.codeVerifier || req.cookies.codeVerifier; 

//   console.log('Code:', code);
//   console.log('State:', state);
//   console.log('Code Verifier:', codeVerifier);

//   const tokenParams = {
//     client_id: CLIENT_ID,
//     code: code,
//     redirect_uri: REDIRECT_URI,
//     grant_type: 'authorization_code',
//     code_verifier: codeVerifier,

//   };

//   try {
//     const response = await axios.post(TOKEN_URL, querystring.stringify(tokenParams));
//     const token = response.data.access_token;
//     console.log('Access Token:', token);
//     res.send(`Access Token: ${token}`);
//   } catch (error) {
//     console.error('Failed to exchange code for token:', error.response.data);
//     res.status(500).send('Failed to exchange authorization code for token');
//   }
// });

// app.listen(8081, () => {
//   console.log('Express server started on port 8081');
// });


// const jwt = require("jsonwebtoken")
// const fs = require("fs")
// const crypto = require('crypto')


// const safeBase64EncodedThumbprint = (thumbprint) => {
//     var numCharIn128BitHexString = 128/8*2
//     var numCharIn160BitHexString = 160/8*2
//     var thumbprintSizes  = {}
//     thumbprintSizes[numCharIn128BitHexString] = true
//     thumbprintSizes[numCharIn160BitHexString] = true
//     var thumbprintRegExp = /^[a-f\d]*$/
  
//     var hexString = thumbprint.toLowerCase().replace(/:/g, '').replace(/ /g, '')
  
//     if (!thumbprintSizes[hexString.length] || !thumbprintRegExp.test(hexString)) {
//       throw 'The thumbprint does not match a known format'
//     }
    
//     var base64 = (Buffer.from(hexString, 'hex')).toString('base64')
//     return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
//   }

// async function getAccessTokenWithCert() {
//       const clientId =  "10a62c80-0d4f-4cc1-89c8-667cfb4c3831" 
//       const tenantId =  "8059c3d5-a962-4394-8b62-ef7c9211422a"  
//       const resource = "https://veling.sharepoint.com"
      
//       const certificateFile = fs.readFileSync("/Users/velingeorgiev/Documents/ocr/orcpdfsharepoint/certs/newcertificate.pem", "utf-8")
//       const X509Certificate = crypto.X509Certificate
//       const x509 = new X509Certificate(certificateFile)
//       const thumbprint = x509.fingerprint

//       const options = {
//         header: {
//           alg: "RS256",
//           typ: "JWT",
//           x5t: safeBase64EncodedThumbprint(thumbprint),
//           kid: thumbprint.toUpperCase().replace(/:/g, '').replace(/ /g, '')
//         },
//         expiresIn: "1h",
//         audience: `https://login.microsoftonline.com/${tenantId}/oauth2/v2.0/token`,
//         issuer:  clientId,
//         subject: clientId,
//         notBefore: Math.floor(Date.now() / 1000) + 60,
//         jwtid: '22b3bb26-e046-42df-9c96-65dbd72c1c81'
//       }

//       const privateKey = fs.readFileSync("/Users/velingeorgiev/Documents/ocr/orcpdfsharepoint/certs/privatekey.pem", "utf-8")

//       const token = jwt.sign({}, privateKey, options)

//       const response = await fetch(
//         `https://login.microsoftonline.com/${tenantId}/oauth2/v2.0/token`,
//         {
//           method: "POST",
//           headers: {
//             "Content-Type": "application/x-www-form-urlencoded"
//           },
//           body: `grant_type=client_credentials&client_id=${clientId}&client_assertion_type=urn%3Aietf%3Aparams%3Aoauth%3Aclient-assertion-type%3Ajwt-bearer&client_assertion=${token}&scope=${`${resource}/.default`}` //&resource=${resource}
//         }
//       );
    
//       return response.json().then(data => {
//         console.log(data)
//         return data.access_token
//       })
//     }


// const crypto = require('crypto');
// const express = require('express');
// const axios = require('axios');
// const querystring = require('querystring');
// const cookieParser = require('cookie-parser');

// const app = express();
// app.use(cookieParser());

// const CLIENT_ID = "10a62c80-0d4f-4cc1-89c8-667cfb4c3831";
// const REDIRECT_URI = 'http://localhost:8081/auth/azure/callback';
// const AUTH_URL = 'https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/authorize';
// const TOKEN_URL = 'https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/token';
// const SCOPE = 'https://management.azure.com/.default';

// function base64URLEncode(buffer) {
//   return buffer.toString('base64')
//       .replace(/\+/g, '-')
//       .replace(/\//g, '_')
//       .replace(/=/g, '');
// }

// function generateCodeVerifier() {
//   return base64URLEncode(crypto.randomBytes(32));
// }

// function sha256(buffer) {
//   return crypto.createHash('sha256').update(buffer).digest();
// }

// function generateCodeChallenge(verifier) {
//   return base64URLEncode(sha256(verifier));
// }

// const generateRandomString = (length) => {
//   return crypto.randomBytes(length).toString('hex');
// };

// app.get('/auth/azure/login', (req, res) => {
//   const codeVerifier = generateCodeVerifier();
//   const codeChallenge = generateCodeChallenge(codeVerifier);
//   const state = generateRandomString(16);

//   // Store the code verifier in the client's browser local storage
//   res.cookie('codeVerifier', codeVerifier, { httpOnly: true });
 
//   const queryParams = querystring.stringify({
//     client_id: CLIENT_ID,
//     response_type: 'code',
//     redirect_uri: REDIRECT_URI,
//     scope: SCOPE,
//     code_challenge: codeChallenge,
//     code_challenge_method: 'S256',
//     state: state,
//   });

//   const authUrl = `${AUTH_URL}?${queryParams}`;
//   res.redirect(authUrl);
// });

// app.get('/auth/azure/callback', async (req, res) => {
//   const code = req.query.code;
//   const state = req.query.state; // You may want to validate the state parameter here
//   const codeVerifier = req.query.codeVerifier || req.cookies.codeVerifier; 

//   console.log('Code:', code);
//   console.log('State:', state);
//   console.log('Code Verifier:', codeVerifier);

//   const tokenParams = {
//     client_id: CLIENT_ID,
//     code: code,
//     redirect_uri: REDIRECT_URI,
//     grant_type: 'authorization_code',
//     code_verifier: codeVerifier,
//     scope: SCOPE,
//   };

//   try {
//     const response = await axios.post(TOKEN_URL, querystring.stringify(tokenParams));
//     const token = response.data.access_token;
//     console.log('Access Token:', token);
//     res.send(`Access Token: ${token}`);
//   } catch (error) {
//     console.error('Failed to exchange code for token:', error.response.data);
//     res.status(500).send('Failed to exchange authorization code for token');
//   }
// });

// app.listen(8081, () => {
//   console.log('Express server started on port 8081');
// });
