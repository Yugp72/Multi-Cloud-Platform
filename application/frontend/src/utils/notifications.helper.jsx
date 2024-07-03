import { notifications } from "@mantine/notifications";
import { IconCheck, IconX } from "@tabler/icons-react";

export const showSuccessNotification = (title, message) => {
  notifications.show({
    icon: <IconCheck size="1.1rem" />,
    title: title || "Success",
    message: message,
    color: "teal",
  });
};

export const showErrorNotification = (title, message) => {
  notifications.show({
    icon: <IconX size="1.1rem" />,
    title: title || "Error",
    message: message,
    color: "red",
  });
};