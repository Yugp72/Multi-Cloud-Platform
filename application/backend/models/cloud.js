const { Sequelize, DataTypes } = require('sequelize');

module.exports = function(sequelize) {
  return sequelize.define('CloudAccount', {
    AccountID: {
      autoIncrement: true,
      type: DataTypes.INTEGER,
      allowNull: false,
      primaryKey: true
    },
    UserID: {
      type: DataTypes.INTEGER,
      allowNull: false,
      references: {
        model: 'User',
        key: 'UserID'
      }
    },
    CloudProvider: {
      type: DataTypes.STRING(50),
      allowNull: false
    },
    AccessKey: {
      type: DataTypes.STRING(255),
      allowNull: true 
    },
    SecretKey: {
      type: DataTypes.STRING(255),
      allowNull: true 
    },
    SubscriptionID: {
      type: DataTypes.STRING(255),
      allowNull: true 
    },
    TenantID: {
      type: DataTypes.STRING(255),
      allowNull: true 
    },
    ClientID: {
      type: DataTypes.STRING(255),
      allowNull: true 
    },
    ClientSecret: {
      type: DataTypes.STRING(255),
      allowNull: true 
    },
    Region: {
      type: DataTypes.STRING(50),
      allowNull: true
    },
    AdditionalInformation: {
      type: DataTypes.TEXT,
      allowNull: true
    },
    ClientEmail: {
      type: DataTypes.STRING(255),
      allowNull: true
    },
    PrivateKey: {
      type: DataTypes.STRING(255),
      allowNull: true
    },
    ProjectID: {
      type: DataTypes.STRING(255),
      allowNull: true
    },
    KeyFile: {
      type: DataTypes.BLOB,
      allowNull: true
    }
  }, {
    sequelize,
    tableName: 'CloudAccount',
    timestamps: false,
    indexes: [
      {
        name: "PRIMARY",
        unique: true,
        using: "BTREE",
        fields: [
          { name: "AccountID" },
        ]
      },
    ]
  });
};
