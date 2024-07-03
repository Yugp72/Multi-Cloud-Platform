const util = require('util');
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const typeDefs = require("./typeDefs");
const resolvers = require("./resolvers");
const context = require("../middleware/auth")
require('dotenv').config();

const models = require("../models");

const server = new ApolloServer({
  typeDefs,
  resolvers,
  cors: {
    origin: "http://localhost:3000/",
    credentials: true,
  },
});

models.sequelize.authenticate().then(() => {
  console.log('Connection has been established successfully.');
  return startStandaloneServer(server, {
      listen: { port: 8000},
      context: context,
  });
}).then((server) => {
  console.log(`🚀  Server ready at: ${server.url}`);
})
.catch((error) => {
  console.error("Error connecting to MySQL:", error);
});