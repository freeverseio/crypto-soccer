const express = require('express');
const bodyParser = require('body-parser');
const { ApolloServer } = require('apollo-server-express');
const { HttpLink } = require('apollo-link-http');
const {
  introspectSchema,
  makeRemoteExecutableSchema,
  mergeSchemas,
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

    extend type Mutation {
      setTeamName(input: SetTeamNameInput!): ID!
      setTeamManagerName(input: SetTeamManagerNameInput!): ID!

    }
  `;

  let schemas = [];
  horizonRemoteSchema && schemas.push(horizonRemoteSchema);
  schemas.push(linkTypeDefs);

  const schema = mergeSchemas({
    schemas,
    resolvers,
  });

  const server = new ApolloServer({ schema });
  const app = express();
  app.use(bodyParser.json());
  app.post('/order/started', (req, res) => {
    console.log(`Request in /order/started:`);
    console.log('main -> req', req);
    console.log('Req body: ', JSON.stringify(req.body));
    console.log('Req headers: ', JSON.stringify(req.headers));
    res.send('OK');
  });
  server.applyMiddleware({ app });
  app.listen({ port: 4000 }, () => {
    console.log(
      `🚀 Server ready at http://localhost:4000${server.graphqlPath}`
    );
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
