const { GraphQLError } = require("graphql");
const jwt = require("jsonwebtoken");
const bcrypt = require("bcrypt");
const models = require("../models");

// Function to generate JWT token
const generateToken = (user) => {
  const token = jwt.sign({ userID: user.UserID, username: user.Username }, process.env.JWT_SECRET, { expiresIn: '1h' });
  return token;
};

const resolvers = {
  Query: {
    user: async (_, { UserID }) => {
      try {
        const user = await models.user.findOne({ where: { UserID: UserID } });
        console.log(user);
        return user;
      } catch (error) {
        console.log(UserID);
        console.error("Error fetching user:", error);
        throw new GraphQLError("Error fetching user");
      }
    },
    cloudAccount: async (_, { UserID }) => {
      try {
        // Fetch the cloud accounts by UserID
        const cloudAccounts = await models.cloud.findAll({ where: { UserID: UserID } });
        
        // Filter out null values from each cloud account
        const filteredCloudAccounts = cloudAccounts.map(account => {
          // Create a new object with non-null values
          const filteredAccount = Object.fromEntries(
            Object.entries(account).filter(([key, value]) => value !== null)
          );
    
          return filteredAccount.dataValues;
        });
        console.log(filteredCloudAccounts);
    
        return filteredCloudAccounts;
      } catch (error) {
        console.error("Error fetching cloud accounts:", error);
        throw new GraphQLError("Error fetching cloud accounts");
      }
    }
    
    
  },
  Mutation: {
    signup: async (_, { signupInput: { Username, Email, Password, RolePermissionLevel } }) => {

      const existingUser = await models.user.findOne({ where: { Email } });
      if (existingUser) {
        throw new GraphQLError("Email is already in use");
      }
      const users = await models.user.findAll();
      const size = users.length;
      console.log(size);
      console.log(users);
      //clear the table data of user

      const hashedPassword = await bcrypt.hash(Password, 10);
      console.log('Password is hashed');
      console.log('Username: ' + Username);
      console.log('Email: ' + Email);
      console.log('Password: ' + hashedPassword);
      console.log('Role Permission Level: ' + RolePermissionLevel);

      const userid = size.toString() + 'xyz';
      console.log('UserID: ' + userid);
      const newUser = await models.user.create({
        //userid size into string
        UserID: userid,

        Username: Username,
        Email: Email,
        Password: hashedPassword,
        RolePermissionLevel: RolePermissionLevel
      });
      try {
        console.log(newUser); // Optionally log the created user
        return newUser; // Return the created user
      } catch (error) {
        // Handle validation errors
        console.error("Error creating user:", error);
        throw new GraphQLError("Error creating user");
      }


    },
    login: async (_, { loginInput: { Email, Password } }) => {
      const user = await models.user.findOne({ where: { Email }});
      console.log("user data" + user);
      if (!user) {
        throw new GraphQLError("User not found");
      }
      const isPasswordValid = await bcrypt.compare(Password, user.Password);

      if (!isPasswordValid) {
        throw new GraphQLError("Incorrect password");
      }
      const token = generateToken(user);
      return { token };
    },
    addCloudAccount: async (_, { cloudAccountInput }) => {
      try {
        const newCloudAccount = await models.cloud.create(cloudAccountInput);
        return newCloudAccount;
      } catch (error) {
        console.error("Error adding cloud account:", error);
        throw new GraphQLError("Error adding cloud account");
      }
    }
  }
};

module.exports = resolvers;
