const { ApolloServer } = require("apollo-server");
const { HttpLink } = require("apollo-link-http");
const { introspectSchema, makeRemoteExecutableSchema, mergeSchemas } = require("graphql-tools");
const fetch = require("node-fetch");

const createRemoteSchema = async uri => {
  const link = new HttpLink({ uri, fetch });
  const schema = await introspectSchema(link);
  const executableSchema = makeRemoteExecutableSchema({
    schema,
    link
  });
  return executableSchema;
};

const main = async () => {
  const universeRemoteSchema = await createRemoteSchema("http://165.22.66.118:4000/graphql");
  const marketRemoteSchema = await createRemoteSchema("http://165.22.66.118:4001/graphql");

  const schema = mergeSchemas({
    schemas: [
      universeRemoteSchema,
      marketRemoteSchema
    ]
  });

  const server = new ApolloServer({ schema });

  // This `listen` method launches a web-server.  Existing apps
  // can utilize middleware options, which we'll discuss later.
  server.listen().then(({ url }) => {
    console.log(`🚀  Server ready at ${url}`);
  });
};

try {
  main();
} catch (e) {
  console.log(e, e.message, e.stack);
}
