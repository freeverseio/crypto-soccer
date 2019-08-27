const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = makeExtendSchemaPlugin(build => {
  // Get any helpers we need from `build`
  const { pgSql: sql, inflection } = build;

  return {
    typeDefs: gql`
    extend type Mutation {
        _: Boolean
    }`,
    resolvers: {
      Mutation: {
          _: () => true,
      }
    },
  };
});

module.exports = MyPlugin;