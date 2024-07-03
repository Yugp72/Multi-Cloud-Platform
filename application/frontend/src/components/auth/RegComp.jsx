import { useRouter } from 'next/navigation';
import {
  Title,
  Radio,
  Text,
  TextInput,
  Button,
  LoadingOverlay,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { gql, useMutation } from '@apollo/client';

import styles from './RegComp.module.css';

import {
  showErrorNotification,
} from '@/utils/notifications.helper';


const Signup = gql(`
    mutation Signup($signupInput: SignupInput) {
      signup(signupInput: $signupInput) {
        Email
        RolePermissionLevel
        UserID
        Username
      }
    }
  `);

function RegComp(onClose, showLoadingOverlay) {
  const router = useRouter();
  const form = useForm({
    initialValues: {
      Email: "",
      Password: "",
      RolePermissionLevel: "",
      Username: "",
    },

    validate: {
      Email: (value) => value.trim().length > 0,
      Password: (value) => value.trim().length > 0,
      RolePermissionLevel: (value) => value.trim().length > 0,
      Username: (value) => value.trim().length > 0,
    },
  });
  const [signup] = useMutation(Signup);

  const createUser = async () => {
    try {
      const {loading, data} = await signup({
        variables: { "signupInput": {
          "Email": form.values.Email,
          "Password": form.values.Password,
          "RolePermissionLevel": form.values.RolePermissionLevel,
          "Username": form.values.Username
        }},
      });
      console.log("AUTHENTICATE -- SUCCESS", data);
      router.push(`/dashboard/${info.role}`);
    } catch (error) {
      console.log("AUTHENTICATE -- ERROR", error);
      showErrorNotification("Failed to create user", error?.message);
    }
    form.reset();
  };

  return (
    <div className={styles.container}>
      <LoadingOverlay
        loaderProps={{
          variant: "bars",
        }}
        visible={signup.loading}
        overlayProps={{ radius: "sm", blur: 2 }}
      />

      <Title order={4}>
        Create Account
      </Title>
      <div className={styles.form}>
        <div className={styles.methodSelections}>
          <div className={styles.methodContainer}>
              <div className={styles.inputs}>
                <TextInput
                  placeholder="Email"
                  {...form.getInputProps("Email")}
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
                <TextInput 
                  placeholder="Password"
                  {...form.getInputProps("Password")}
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
                <TextInput 
                  placeholder="Username"
                  {...form.getInputProps("Username")}
                  classNames={{
                    input: styles.defaultRadius,
                  }}
                  size="md"
                />
                <TextInput 
                  placeholder="Role Permission Level"
                  {...form.getInputProps("RolePermissionLevel")}
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
              onClick={createUser}
            >
              Create Account
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default RegComp;
