const { selectTeamName, selectTeamManagerName } = require('../repositories');
const teamByTeamId = require('./teamByTeamId');
const getAuctionResolver = require('./getAuction');
const allAuctionsResolver = require('./allAuctions');

const resolvers = ({ horizonRemoteSchema }) => {
  return {
    Team: {
      name: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return selectTeamName({ teamId: team.teamId }).then((result) => {
            return result && result.team_name ? result.team_name : team.name;
          });
        },
      },
      managerName: {
        selectionSet: `{ teamId }`,
        resolve(team, args, context, info) {
          return selectTeamManagerName({ teamId: team.teamId }).then((result) => {
            return result && result.team_manager_name ? result.team_manager_name : team.managerName;
          });
        },
      },
    },
    Bid: {
      teamByTeamId: {
        fragment: `... on Bid { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    Player: {
      teamByTeamId: {
        fragment: `... on Player { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    Query: {
      allAuctions: allAuctionsResolver,
      getAuction: getAuctionResolver,
    },
  };
};

module.exports = resolvers;
