import React, { useState } from 'react';
import { Text, Button, Space } from '@mantine/core';

const Collapsible = ({ title, children }) => {
  const [collapsed, setCollapsed] = useState(true);

  const toggleCollapsed = () => {
    setCollapsed(!collapsed);
  };

  return (
    <div>
      <Button onClick={toggleCollapsed} variant="light">
        {title} {!collapsed ? '-' : '+'}
      </Button>
      {!collapsed && <div>{children}</div>}
    </div>
  );
};

export default Collapsible;
