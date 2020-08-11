const { ApolloServer } = require("apollo-server");
const { HttpLink } = require("apollo-link-http");
const { introspectSchema, makeRemoteExecutableSchema, mergeSchemas } = require("graphql-tools");
const selectPlayerName = require("./repositories/selectPlayerName.js")
const selectTeamName = require("./repositories/selectTeamName.js")
const selectTeamManagerName = require("./repositories/selectTeamManagerName.js")
const fetch = require("node-fetch");
const program = require("commander");
const updatePlayerName = require("./repositories/updatePlayerName.js");
const updateTeamName = require("./repositories/updateTeamName.js");
const updateTeamManagerName = require("./repositories/updateTeamManagerName.js");
const version = require("../package.json").version;

const createRemoteSchema = async uri => {
  const link = new HttpLink({ uri, fetch });
  const schema = await introspectSchema(link);
  const executableSchema = makeRemoteExecutableSchema({
    schema,
    link
  });
  return executableSchema;
};

program
  .version(version)
  .option("-h, --horizonUrl <url>", "graphql horizon url", "")
  .option("-g, --gameUrl <url>", "graphql game url", "")
  .parse(process.argv);

const { horizonUrl, gameUrl } = program;

console.log("--------------------------------------------------------");
console.log("horizonUrl       : ", horizonUrl);
console.log("gameUrl         : ", gameUrl);
console.log("--------------------------------------------------------");

const main = async () => {
  const horizonRemoteSchema = await createRemoteSchema(horizonUrl);
  
  const linkTypeDefs = `
    input SetPlayerNameInput2 {
      signature: String!
      playerId: ID!
      name: String!
    }

    input SetTeamNameInput2 {
      signature: String!
      teamId: ID!
      name: String!
    }
  
    input SetTeamManagerNameInput2 {
      signature: String!
      teamId: ID!
      name: String!
    }

    extend type Mutation {
      setPlayerName2(input: SetPlayerNameInput2!): ID!
      setTeamName2(input: SetTeamNameInput2!): ID!
      setTeamManagerName2(input: SetTeamManagerNameInput2!): ID!

    }
  `;

  const resolvers = {
    Player: {
      name: {
        selectionSet: `{ playerId }`,
        resolve(player, args, context, info) {
          return selectPlayerName({ playerId: player.playerId }).then(result => { 
            return result && result.player_name ? result.player_name : player.name
          })
        }
      },
    },
    Team: {
      name: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return selectTeamName({ teamId: team.teamId }).then(result => { 
            return result && result.team_name ? result.team_name : team.name
          })
        }
      },
      managerName: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return selectTeamManagerName({ teamId: team.teamId }).then(result => { 
            return result && result.team_manager_name ? result.team_manager_name : team.managerName
          })
        }
      },
    },
    Mutation: {
      setPlayerName2: async (_, { input: { playerId, name, signature } }) => {
          await updatePlayerName({ playerId, playerName: name })
          return playerId 
        },
      setTeamName2: async (_, { input: { teamId, name, signature } }) => {
        await updateTeamName({ teamId, teamName: name })
        return teamId 
      },
      setTeamManagerName2: async (_, { input: { teamId, name, signature } }) => {
        await updateTeamManagerName({ teamId, teamManagerName: name })
        return teamId 
      },
    }
  };

  let schemas = [];
  horizonRemoteSchema && schemas.push(horizonRemoteSchema);
  schemas.push(linkTypeDefs);

  const schema = mergeSchemas({
    schemas,
    resolvers,
  });
  
  const server = new ApolloServer({ schema, tracing: true });

  server.listen().then(({ url }) => {
    console.log(`🚀  Server ready at ${url}`);
  });
};

const run = () => {
  main()
  .catch(e => {
    console.error(e);
    console.log("wainting ......");
    setTimeout(run, 3000);
  })
};

run();
