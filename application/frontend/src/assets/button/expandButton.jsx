import React from 'react';
import styles from './expandButton.module.css'; // Import your CSS file for button styling

import PropTypes from 'prop-types';

const CustomButton = ({ name, providerSymbol, expandSymbol, onClick }) => {
  return (
    <button className={styles.custombutton} onClick={onClick}>
      <span className={styles.buttontext}>{name}</span>
      <span className={styles.symbol}>{providerSymbol}</span>
      <span className={styles.symbol}>{expandSymbol}</span>
    </button>
  );
};

CustomButton.propTypes = {
  name: PropTypes.string.isRequired,
  providerSymbol: PropTypes.symbol.isRequired,
  expandSymbol: PropTypes.symbol.isRequired,
  onClick: PropTypes.func.isRequired,
};

export default CustomButton;
