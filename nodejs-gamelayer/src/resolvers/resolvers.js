const Web3 = require('web3');
const {
  selectTeamName,
  selectTeamManagerName,
  updateTeamName,
  updateTeamManagerName,
  selectLeaderboard,
} = require('../repositories');
const { TeamValidation } = require('../validations');

const web3 = new Web3('');

const resolvers = {
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
          return result && result.team_manager_name
            ? result.team_manager_name
            : team.managerName;
        });
      },
    },
  },
  Mutation: {
    setTeamName: async (_, { input: { teamId, name, signature } }) => {
      const teamValidation = new TeamValidation({
        teamId,
        name,
        signature,
        web3,
      });
      const isSignerOwner = await teamValidation.isSignerOwner();

      if (isSignerOwner) {
        await updateTeamName({ teamId, teamName: name });
        return teamId;
      } else {
        return 'Signer is not the team owner';
      }
    },
    setTeamManagerName: async (_, { input: { teamId, name, signature } }) => {
      const teamValidation = new TeamValidation({
        teamId,
        name,
        signature,
        web3,
      });
      const isSignerOwner = await teamValidation.isSignerOwner();

      if (isSignerOwner) {
        await updateTeamManagerName({ teamId, teamManagerName: name });
        return teamId;
      } else {
        return 'Signer is not the team owner';
      }
    },
  },
  Query: {
    selectLeaderboard: async (
      _,
      { input: { timezoneIdx, countryIdx, leagueIdx } }
    ) => {
      return selectLeaderboard({ timezoneIdx, countryIdx, leagueIdx });
    },
  },
};

module.exports = resolvers;
