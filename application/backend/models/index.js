require('dotenv').config();

const {Sequelize, DataTypes} = require("sequelize");

var initModels = require("./init-models");

const sequelize = new Sequelize('multicloud', 'newuser', 'password' ,{
  host: '127.0.0.1',
  port: '3307',
  dialect: 'mysql',
  define: {
      freezeTableName: true,
  },
});

const database = sequelize.config.database;
const username = sequelize.config.username;
const password = sequelize.config.password;
const host = sequelize.config.host;
const dialect = sequelize.config.dialect;

console.log('Database:', database);
console.log('Username:', username);
console.log('Password:', password);
console.log('Host:', host);
console.log('Dialect:', dialect);


var models = initModels(sequelize);

models.sequelize = sequelize;

module.exports = models;
