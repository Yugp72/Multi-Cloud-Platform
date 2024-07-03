import { gql } from '@apollo/client';
import { jwtDecode } from 'jwt-decode';
import { useQuery } from '@apollo/client';

const GET_CLOUD_ACCOUNTS = gql`
    query cloudAccount($userId: Int) {
        cloudAccount(UserID: $userId) {
            AccessKey
            AccountID
            AdditionalInformation
            ClientID
            ClientSecret
            CloudProvider
            Region
            SecretKey
            SubscriptionID
            TenantID
            UserID
        }
    }`;

const useAllConnectedAccounts = () => {
    const temptoken = String(localStorage.getItem('token'));
    const decodedToken = jwtDecode(temptoken);
    const userID = parseInt(decodedToken.userID);
    const token = localStorage.getItem('token');

    const { data, loading, error } = useQuery(GET_CLOUD_ACCOUNTS, {
        variables: { userId: userID },
        context: {
            headers: {
                Authorization: token || '',
            }
        }
    });

    console.log(data);
    // store 
    return { data, loading, error };
};

export default useAllConnectedAccounts;
