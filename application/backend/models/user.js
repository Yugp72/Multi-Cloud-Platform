const Sequelize = require('sequelize');
const { validator } = require('sequelize/lib/utils/validator-extras');

module.exports = function(sequelize, DataTypes) {
  return sequelize.define('User', {
    UserID: {
      autoIncrement: true,
      type: DataTypes.INTEGER,
      allowNull: false,
      primaryKey: true
    },
    Username: {
      type: DataTypes.STRING(255),
      allowNull: false,
      
    },
    Email: {
      type: DataTypes.STRING(255),
      allowNull: false,
      unique : true,
      validator : {
        unique:{
          msg: 'Please enter unique email'
        }
      }

    },

    Password: {
      type: DataTypes.STRING(255),
      allowNull: false

    },

    RolePermissionLevel: {
      type: DataTypes.STRING(50),
      allowNull: false
    }
  }, {
    sequelize,
    tableName: 'User',
    timestamps: false,
    indexes: [
      {
        name: "PRIMARY",
        unique: true,
        using: "BTREE",
        fields: [
          { name: "id" },
        ]
      },
    ]
  });
};
