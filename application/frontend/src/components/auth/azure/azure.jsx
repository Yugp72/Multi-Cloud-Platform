const crypto = require('crypto');
const express = require('express');
const axios = require('axios');
const querystring = require('querystring');
const cookieParser = require('cookie-parser');

const app = express();
app.use(cookieParser());

const REDIRECT_URI = 'http://localhost:8081/auth/azure/callback';
const AUTH_URL = 'https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/authorize';
const TOKEN_URL = 'https://login.microsoftonline.com/8059c3d5-a962-4394-8b62-ef7c9211422a/oauth2/v2.0/token';
const SCOPE = 'https://management.azure.com/.default';

function base64URLEncode(buffer) {
  return buffer.toString('base64')
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '');
}

function generateCodeVerifier() {
  return base64URLEncode(crypto.randomBytes(32));
}

function sha256(buffer) {
  return crypto.createHash('sha256').update(buffer).digest();
}

function generateCodeChallenge(verifier) {
  return base64URLEncode(sha256(verifier));
}

const generateRandomString = (length) => {
  return crypto.randomBytes(length).toString('hex');
};

app.get('/auth/azure/login', async (req, res) => {
  try {
    // Fetch Azure account data
    const { data, loading, error } = await useAllConnectedAccounts();

    // Assuming the data structure has an array of accounts
    const azureAccounts = data.cloudAccount.filter(account => account.CloudProvider === 'Azure');

    // Assuming we're using the first Azure account found
    const azureClientId = azureAccounts[0]?.ClientId;

    if (!azureClientId) {
      throw new Error('Azure client ID not found');
    }

    const codeVerifier = generateCodeVerifier();
    const codeChallenge = generateCodeChallenge(codeVerifier);
    const state = generateRandomString(16);

    // Store the code verifier in the client's browser local storage
    res.cookie('codeVerifier', codeVerifier, { httpOnly: true });
   
    const queryParams = querystring.stringify({
      client_id: azureClientId, // Use the Azure client ID
      response_type: 'code',
      redirect_uri: REDIRECT_URI,
      scope: SCOPE,
      code_challenge: codeChallenge,
      code_challenge_method: 'S256',
      state: state,
    });

    const authUrl = `${AUTH_URL}?${queryParams}`;
    res.redirect(authUrl);
  } catch (error) {
    console.error('Error fetching Azure account data:', error);
    res.status(500).send('Failed to fetch Azure account data');
  }
});

app.get('/auth/azure/callback', async (req, res) => {
  const code = req.query.code;
  const state = req.query.state; // You may want to validate the state parameter here
  const codeVerifier = req.query.codeVerifier || req.cookies.codeVerifier; 

  console.log('Code:', code);
  console.log('State:', state);
  console.log('Code Verifier:', codeVerifier);

  const tokenParams = {
    client_id: azureClientId, // Use the Azure client ID
    code: code,
    redirect_uri: REDIRECT_URI,
    grant_type: 'authorization_code',
    code_verifier: codeVerifier,
    scope: SCOPE,
  };

  try {
    const response = await axios.post(TOKEN_URL, querystring.stringify(tokenParams));
    const token = response.data.access_token;
    console.log('Access Token:', token);
    res.send(`Access Token: ${token}`);
  } catch (error) {
    console.error('Failed to exchange code for token:', error.response.data);
    res.status(500).send('Failed to exchange authorization code for token');
  }
});

app.listen(8081, () => {
  console.log('Express server started on port 8081');
});
