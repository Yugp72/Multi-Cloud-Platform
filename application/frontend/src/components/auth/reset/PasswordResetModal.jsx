import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { jwtDecode } from "jwt-decode";
import { createHmac } from 'crypto';
const secret = 'abcdefg';
import {
  Title,
  Text,
  TextInput,
  Button,
  LoadingOverlay,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { gql, useMutation } from '@apollo/client';

import styles from './PasswordResetModal.module.css';

import {
  showSuccessNotification,
  showErrorNotification,
} from '@/utils/notifications.helper';

const ChangePwd = gql(`
  mutation ChangePwd($changePwd: changePwd) {
    changePwd(changePwd: $changePwd) {
      token
    }
  }
`);

function PasswordResetComp({ onClose, showLoadingOverlay }) {
  const router = useRouter();
  const form = useForm({
    initialValues: {
      username: '',
      oldPassword:'',
      newPassword: '',
      confirmPassword: '',
    },
    validate: { 
      username: (value) => value.trim().length > 0,
      oldPassword: (value) => value.trim().length > 0,
      newPassword: (value) => value.trim().length > 0,
      confirmPassword: (value) => value === form.values.newPassword,
    },
  });

  const [changepwd] = useMutation(ChangePwd);
  const [resetSuccess, setResetSuccess] = useState(false);

  const handleResetPassword = async () => {
    console.log(form.values);
    console.log(JSON.stringify(createHmac('sha256', secret).update(form.values.oldPassword).digest('hex')).slice(1,-1));
    console.log(JSON.stringify(createHmac('sha256', secret).update(form.values.newPassword).digest('hex')).slice(1,-1))
    try {
      const { data } = await changepwd({
        variables: { 
          changePwd: {
            username: form.values.username,
            oldPassword: JSON.stringify(createHmac('sha256', secret).update(form.values.oldPassword).digest('hex')).slice(1,-1),
            newPassword: JSON.stringify(createHmac('sha256', secret).update(form.values.newPassword).digest('hex')).slice(1,-1),
          },
        },
      });

      console.log('RESET PASSWORD -- SUCCESS', data);
      showErrorNotification('Password Changed',"");
      const info = jwtDecode(data.changePwd.token);
      console.log(data.changePwd.token);
            console.log(info);
            localStorage.setItem("token", data.changePwd.token);
            if(info.username==null){
              setCreatePassword(true);
            }
            else{
              if(info.position==="Student"){
                router.push("\student");
              }
              else if(info.position==="Mess"){
                router.push("\mess");
              }
              else if(info.position==="House keeping"){
                router.push("\housekeeping");
              }
              else{
                router.push("\admin");
              }
            }
    } catch (error) {
      console.log('RESET PASSWORD -- ERROR', error);
      showErrorNotification('Failed to reset password', error?.message);
    }
  };

  return (
    <div>
      {showLoadingOverlay && (
        <LoadingOverlay
          loaderProps={{
            variant: 'bars',
          }}
          visible={resetPassword.loading}
          overlayProps={{ radius: 'sm', blur: 2 }}
        />
      )}

      {resetSuccess ? (
        <>
          <Title order={4}>Password Reset Successful</Title>
          <Text>Your password has been reset successfully.</Text>
        </>
      ) : (
        <>
          <Title order={4}>Reset Password</Title>
          <div className={styles.form}>
            <div className={styles.methodContainer}>
              <div className={styles.inputs}>
                <TextInput
                  placeholder="Username"
                  {...form.getInputProps('username')}
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
                <TextInput
                  placeholder="Old Password"
                  {...form.getInputProps("oldPassword")}
                  type="password"
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
                <TextInput
                  placeholder="New Password"
                  {...form.getInputProps('newPassword')}
                  type="password"
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
                <TextInput
                  placeholder="Confirm New Password"
                  {...form.getInputProps('confirmPassword')}
                  type="password"
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
              </div>
              <Button
                fullWidth
                size="md"
                classNames={{
                  root: styles.defaultRadius,
                }}
                onClick={handleResetPassword}
              >
                Reset Password
              </Button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

export default PasswordResetComp;