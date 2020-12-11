const { ApolloServer } = require('apollo-server');
const { HttpLink } = require('apollo-link-http');
const {
  introspectSchema,
  makeRemoteExecutableSchema,
  mergeSchemas,
  transformSchema,
  FilterRootFields,
  FilterTypes,
} = require('graphql-tools');
const fetch = require('node-fetch');
const resolvers = require('./resolvers/resolvers.js');
const { horizonConfig } = require('./config.js');

const createRemoteSchema = async (uri) => {
  const link = new HttpLink({ uri, fetch });
  const schema = await introspectSchema(link);
  const executableSchema = makeRemoteExecutableSchema({
    schema,
    link,
  });
  return executableSchema;
};
const horizonUrl = horizonConfig.url;
console.log('--------------------------------------------------------');
console.log('horizonUrl       : ', horizonUrl);
console.log('--------------------------------------------------------');

const main = async () => {
  const horizonRemoteSchema = await createRemoteSchema(horizonUrl);

  const linkTypeDefs = `
    input SetTeamNameInput {
      signature: String!
      teamId: ID!
      name: String!
    }
  
    input SetTeamManagerNameInput {
      signature: String!
      teamId: ID!
      name: String!
    }
    
    input SetMessageInput {
      destinatary: String!
      category: String!
      auctionId: String
      title: String!
      text: String!
      customImageUrl: String
      metadata: String
    }

    input SetBroadcastMessageInput {
      category: String!
      auctionId: String
      title: String!
      text: String!
      customImageUrl: String
      metadata: String
    }

    type Message {
      id: String
      destinatary: String!
      category: String!
      auctionId: String
      title: String!
      text: String!
      customImageUrl: String
      metadata: String
      isRead: Boolean
      createdAt: String
    }

    type Messages {
      totalCount: Int!
      nodes: [Message]
    }

    type PlayerHistoryGraphEncodedSkills {
      encodedSkills: String
    }
    
    type PlayerHistoryGraph {
      nodes: [PlayerHistoryGraphEncodedSkills]
    }

    input CreateBidInput {
  		signature: String!
		  auctionId: ID!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
    }
    
    input SetGetSocialIdInput {
      signature: String!
      teamId: String!
      getSocialId: String!
    }

    extend type Mutation {
      setTeamName(input: SetTeamNameInput!): ID!
      setTeamManagerName(input: SetTeamManagerNameInput!): ID!
      setMessage(input: SetMessageInput!): ID!
      setBroadcastMessage(input: SetBroadcastMessageInput!): Boolean!
      setMailboxStart(teamId: ID!): Boolean
      setMessageRead(id: ID!): Boolean
      setLastTimeLoggedIn(teamId: ID!): Boolean
      createBid(input: CreateBidInput!): ID!
      setGetSocialId(input: SetGetSocialIdInput!): Boolean
    }
    
    extend type Query {
      getMessages(teamId: ID!, auctionId: ID, limit: Int, offset: Int): Messages!
      getNumUnreadMessages(teamId : ID!): Int!
      createBid(input: CreateBidInput!): ID!
      getLastTimeLoggedIn(teamId: ID!): String!
    }

    extend type Team {
      auctionPassByOwner: Boolean!
    }
    
    extend type Player {
      playerHistoryGraphByPlayerId(first: Int!): PlayerHistoryGraph
    }
  `;

  const schemas = [];
  const transformedHorizonRemoteSchema = transformSchema(horizonRemoteSchema, [
    new FilterRootFields((operation, fieldName, field) => {
      if (fieldName == 'processAuction') {
        return false;
      }
      return true;
    }),
    new FilterTypes((typeName, fieldName, field) => {
      if (fieldName == 'ProcessAuctionInput') {
        return false;
      }
      return true;
    }),
    new FilterRootFields((operation, fieldName, field) => {
      if (fieldName == 'createBid') {
        return false;
      }
      return true;
    }),
    new FilterTypes((typeName, fieldName, field) => {
      if (fieldName == 'CreateBidInput') {
        return false;
      }
      return true;
    }),
  ]);
  transformedHorizonRemoteSchema && schemas.push(transformedHorizonRemoteSchema);
  schemas.push(linkTypeDefs);

  const schema = mergeSchemas({
    schemas,
    resolvers: resolvers({ horizonRemoteSchema }),
  });

  const server = new ApolloServer({ schema });

  server.listen().then(({ url }) => {
    console.log(`🚀  Server ready at ${url}`);
  });
};

const run = () => {
  main().catch((e) => {
    console.error(e);
    console.log('wainting ......');
    setTimeout(run, 3000);
  });
};

run();
