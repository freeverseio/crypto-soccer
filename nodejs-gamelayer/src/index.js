const { ApolloServer } = require("apollo-server");
const { HttpLink } = require("apollo-link-http");
const { introspectSchema, makeRemoteExecutableSchema, mergeSchemas } = require("graphql-tools");
const selectPlayerName = require("./repositories/selectPlayerName.js")
const fetch = require("node-fetch");
const program = require("commander");
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
    extend type Player {
      playerPropsByPlayerId: PlayerProp
      otraCosa: String
    }

    extend type Team {
      teamPropsByTeamId: TeamProp
    }
  `;

  const resolvers = {
    Player: {
      playerPropsByPlayerId: {
        selectionSet: `{ playerId }`,
        resolve(player, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: gameRemoteSchema,
            operation: 'query',
            fieldName: 'playerPropByPlayerId',
            args: {
              playerId: player.playerId,
              condition: {
                playerId: player.playerId
              }
            },
            context,
            info,
          })
        }
      },
      name: {
        selectionSet: `{ playerId }`,
        resolve(player, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: gameRemoteSchema,
            operation: 'query',
            fieldName: 'playernamebyplayerid',
            args: {
              playerid: player.playerId
            },
            context,
            info,
          }).then(result => {
            console.log("Name result: ", result)

            if(result) {
              return result
            }

            return player.name
          })
        }
      },
      otraCosa: {
        selectionSet: `{ playerId }`,
        resolve(player, args, context, info) {
          console.log("entro otracosa resolve")
          const playerName = selectPlayerName({ playerId: player.playerId })
          console.log("Thi player name", playerName)
          return playerName.then(result => { 
            console.log("resuuult", result.player_name)
            return result.player_name
          })
        }
      },
    },
    Team: {
      teamPropsByTeamId: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: gameRemoteSchema,
            operation: 'query',
            fieldName: 'teamPropByTeamId',
            args: {
              teamId: team.teamId,
              condition: {
                teamId: team.teamId
              }
            },
            context,
            info,
          })
        }
      },
      name: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: gameRemoteSchema,
            operation: 'query',
            fieldName: 'teamnamebyteamid',
            args: {
              teamid: team.teamId
            },
            context,
            info,
          }).then(result => {
            if(result) {
              return result
            }

            return team.name
          })
        }
      },
      managerName: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return info.mergeInfo.delegateToSchema({
            schema: gameRemoteSchema,
            operation: 'query',
            fieldName: 'teammanagernamebyteamid',
            args: {
              teamid: team.teamId
            },
            context,
            info,
          }).then(result => {
            if(result) {
              return result
            }

            return team.name
          })
        }
      },
    },

  };

  let schemas = [];
  horizonRemoteSchema && schemas.push(horizonRemoteSchema);
  gameRemoteSchema && schemas.push(gameRemoteSchema);
  schemas.push(linkTypeDefs);

  const schema = mergeSchemas({
    schemas,
    resolvers,
  });

  const server = new ApolloServer({ schema });

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
