import { useState } from 'react';
import { useRouter } from 'next/navigation';
import PasswordResetComp from './reset/PasswordResetModal';
import RegComp from './RegComp';
import Link from 'next/link';
import {
  Title,
  Radio,
  Text,
  TextInput,
  Button,
  LoadingOverlay,
} from '@mantine/core';
import { gql, useMutation } from '@apollo/client';
import { useForm } from '@mantine/form';
import { jwtDecode } from 'jwt-decode';

import Google from '@/assets/general/Google.svg';
import { notifications } from '@mantine/notifications';

import styles from './AuthComp.module.css';
import Image from 'next/image';
import {
  showSuccessNotification,
  showErrorNotification,
} from '@/utils/notifications.helper';

import { createHmac } from 'crypto';
const secret = 'abcdefg';

const Login = gql(`
  mutation Login($loginInput: LoginInput) {
    login(loginInput: $loginInput) {
      token
    }
  }
`);

function AuthComp() {
  const router = useRouter();
  const form = useForm({
    initialValues: {
      Email: '',
      Password: '',
    },

    validate: {
      Email: (value) => value.trim().length > 0,
      Password: (value) => value.trim().length > 0,
    },
  });

  const [login] = useMutation(Login);
  const [showPasswordReset, setShowPasswordReset] = useState(false);
  const [showRegComp, setShowRegComp] = useState(false);

  const handleRegistrationClose = () => {
    setShowRegComp(false);
  };

  const handleRegistrationReset = () => {
    };

  const eventLogin = async () => {
      try {
      console.log('Login:', form.values);

      const { data } = await login({
        variables: {
          loginInput: {
            Email: form.values.Email,
            Password: form.values.Password,
          },
        },
      });
      
      // Handle login success
      const info = jwtDecode(data.login.token);
      localStorage.setItem('token', data.login.token);
      router.push(`/cloud/`);
    } catch (error) {
      // Handle login error
      console.error('Login Error:', error);
      showErrorNotification('Login Failed', error?.message);
      form.reset();
    }
  };


  return (
    <div className={styles.container}>
      {!showRegComp && 
      <div> 
        <Title order={4}>Login</Title>
      <div className={styles.form}>
        <div className={styles.methodSelections}>
          <div className={styles.methodContainer}>
            <div className={styles.inputs}>
              <TextInput
                placeholder="Username"
                {...form.getInputProps('Email')}
                classNames={{
                  input: styles.defaultRadius,
                }}
                size="md"
              />
              <TextInput
                placeholder="Password"
                {...form.getInputProps('Password')}
                type="Password"
                classNames={{
                  input: styles.defaultRadius,
                }}
                size="md"
              />
            </div>
    
            <span className={styles.forgotPasswordButton}>
                   <a href="#" onClick={() => setShowPasswordReset(true)}>Set Password</a>
                </span>
            
            <Button
              fullWidth
              size="md"
              classNames={{
                root: styles.defaultRadius,
              }}
              onClick={eventLogin}
            >
              Login
            </Button>

            <span
              className={`${styles.defaultRadius} ${styles.registerButton}`}>
              <a href="##" onClick={() => setShowRegComp(true)}> <Text ta="center">Register</Text></a>
            </span>

          </div>
          <div className={styles.methodContainer}>
            <Text ta="center">or</Text>
            <Button
              disabled
              fullWidth
              size="md"
              classNames={{
                root: styles.defaultRadius,
              }}
              variant="outline"
            >
              <Image src={Google} alt="Google" /> &nbsp; Continue with Google
            </Button>
          </div>
        </div>
      </div>
      
      </div>
      }
      {showPasswordReset && <PasswordResetComp onClose={() => setShowPasswordReset(false)} />}
      {showRegComp && <RegComp onClose={handleRegistrationClose} onReset={handleRegistrationReset} />}
    </div>
  );
}

export default AuthComp;
