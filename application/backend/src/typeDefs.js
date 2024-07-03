const { gql } = require("graphql-tag");


module.exports = gql`
  scalar Upload
  
  type User {
    UserID: Int
    Username: String
    Email: String
    RolePermissionLevel: String
  }

  type Session {
    sessionID: Int
    UserID: Int
    sessionToken: String
    expiryDateTime: String
  }
  
  type CloudAccount {
    AccountID: Int
    UserID: Int
    CloudProvider: String
    AccessKey: String
    SecretKey: String
    SubscriptionID: String
    TenantID: String
    ClientID: String
    ClientSecret: String
    Region: String
    AdditionalInformation: String
    ClientEmail: String
    PrivateKey: String
    ProjectID: String
    KeyFile: Upload

  }

  type AccessControl {
    accessControlID: Int
    UserID: Int
    resourceFeatureName: String
    accessLevelPermission: String
  }

  type AuditLog {
    logID: Int
    UserID: Int
    actionPerformed: String
    timestamp: String
    ipAddress: String
  }

  input LoginInput {
    Email: String
    Password: String
  }

  input SignupInput {
    Username: String
    Password: String
    Email: String
    RolePermissionLevel: String
  }

  type Query {
    user(UserID: Int): User
    cloudAccount(UserID: Int): [CloudAccount]
  }

  type Mutation {
    login(loginInput: LoginInput): Token
    signup(signupInput: SignupInput): User
    addCloudAccount(cloudAccountInput: CloudAccountInput): CloudAccount
  }

  input CloudAccountInput {
    UserID: Int
    CloudProvider: String
    AccessKey: String
    SecretKey: String
    SubscriptionID: String
    TenantID: String
    ClientID: String
    ClientSecret: String
    Region: String
    AdditionalInformation: String
    ClientEmail: String
    PrivateKey: String
    ProjectID: String
    KeyFile: Upload
  }

  type Token {
    token: String
  }
`;
