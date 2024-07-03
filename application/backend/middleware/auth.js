const { GraphQLError } = require("graphql");
const jwt = require("jsonwebtoken");

const getUser = async (token) => {
  try {
    if (token) {
      const user = jwt.verify(token, '123');
      return user;
    }
    return null;
  } catch (error) {
    return null;
  }
};

const context = async ({ req }) => {
  const exemptOps = ["IntrospectionQuery", "AddCloudAccount","Login", "Signup", "ChangePwd"];

  if (exemptOps.includes(req.body.operationName)) {
    return {};
  }

  const token = req.headers.authorization || "";
  console.log("My token: " + token);
  const user = await getUser(token);

  if (!user) {
    throw new GraphQLError("User is not authenticated");
  }

  return { user };
};

module.exports = context;
