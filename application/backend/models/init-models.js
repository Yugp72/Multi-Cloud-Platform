var DataTypes = require("sequelize").DataTypes;
var _user = require("./user.js");
var _cloud = require("./cloud.js");

function initModels(sequelize) {
  var user = _user(sequelize, DataTypes);
  var cloud = _cloud(sequelize, DataTypes);
  var sequelize;

  return {
    user,
    sequelize,
    cloud,
  };
}
module.exports = initModels;
module.exports.initModels = initModels;
module.exports.default = initModels;