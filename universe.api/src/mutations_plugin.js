const { makeExtendSchemaPlugin, gql } = require("graphile-utils");

const MyPlugin = () => {
  return makeExtendSchemaPlugin(build => {
    // Get any helpers we need from `build`
    const { pgSql: sql, inflection } = build;

    return {
      typeDefs: gql`
      extend type Mutation {
        createSpecialPlayer(
          playerId: String!,
          name: String!,
          teamId: String!,
          defence: Int!,
          speed: Int!,
          pass: Int!,
          shoot: Int!,
          endurance: Int!,
          preferredPosition: String!,
          potential: Int!,
          dayOfBirth: Int!
        ): Boolean,
        deleteSpecialPlayer(
          playerId: String!
        ): Boolean
      }`,
      resolvers: {
        Mutation: {
          createSpecialPlayer: (_, params) => {
            return true;
          },
        },
      },
    }
  });
};

module.exports = MyPlugin;