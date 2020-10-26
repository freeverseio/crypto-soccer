const Web3 = require('web3');
const { selectTeamName, selectTeamManagerName, updateTeamName, updateTeamManagerName } = require('../repositories');
const { TeamValidation } = require('../validations');
const getMessagesResolver = require('./getMessagesResolver');
const setMessageReadResolver = require('./setMessageReadResolver');
const setMailboxStartResolver = require('./setMailboxStartResolver');
const setBroadcastMessageResolver = require('./setBroadcastMessageResolver');
const setMessageResolver = require('./setMessageResolver');
const teamByTeamId = require('./teamByTeamId');
const teamByHomeTeamId = require('./teamByHomeTeamId');
const teamByVisitorTeamId = require('./teamByVisitorTeamId');
const teamByBuyerTeamId = require('./teamByBuyerTeamId');
const getNumUnreadMessagesResolver = require('./getNumUnreadMessagesResolver');
const getLastTimeLoggedInResolver = require('./getLastTimeLoggedIn');
const setLastTimeLoggedInResolver = require('./setLastTimeLoggedIn');
const createBidResolver = require('./createBidResolver');

const web3 = new Web3('');

const resolvers = ({ horizonRemoteSchema }) => {
  return {
    TeamsHistory: {
      name: {
        fragment: `... on TeamsHistory { teamId }`,
        resolve(team, args, context, info) {
          return selectTeamName({ teamId: team.teamId }).then((result) => {
            return result && result.team_name ? result.team_name : team.name;
          });
        },
      },
    },
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
    MatchEvent: {
      teamByTeamId: {
        fragment: `... on MatchEvent { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    Match: {
      teamByHomeTeamId: {
        fragment: `... on Match { homeTeamId }`,
        resolve(match, args, context, info) {
          return teamByHomeTeamId(match, args, context, info, horizonRemoteSchema);
        },
      },
      teamByVisitorTeamId: {
        fragment: `... on Match { visitorTeamId }`,
        resolve(match, args, context, info) {
          return teamByVisitorTeamId(match, args, context, info, horizonRemoteSchema);
        },
      },
    },
    MatchesHistory: {
      teamByHomeTeamId: {
        fragment: `... on MatchesHistory { homeTeamId }`,
        resolve(match, args, context, info) {
          return teamByHomeTeamId(match, args, context, info, horizonRemoteSchema);
        },
      },
      teamByVisitorTeamId: {
        fragment: `... on MatchesHistory { visitorTeamId }`,
        resolve(match, args, context, info) {
          return teamByVisitorTeamId(match, args, context, info, horizonRemoteSchema);
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
    PlayersHistory: {
      teamByTeamId: {
        fragment: `... on PlayersHistory { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
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
    Offer: {
      teamByBuyerTeamId: {
        fragment: `... on Offer { buyerTeamId }`,
        resolve(offer, args, context, info) {
          return teamByBuyerTeamId(offer, args, context, info, horizonRemoteSchema);
        },
      },
    },
    Tactic: {
      teamByTeamId: {
        fragment: `... on Tactic { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    TacticsHistory: {
      teamByTeamId: {
        fragment: `... on TacticsHistory { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    Training: {
      teamByTeamId: {
        fragment: `... on Training { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info), horizonRemoteSchema;
        },
      },
    },
    TrainingsHistory: {
      teamByTeamId: {
        fragment: `... on TrainingsHistory { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    TeamsProp: {
      teamByTeamId: {
        fragment: `... on TeamsProp { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    TeamsPropsHistory: {
      teamByTeamId: {
        fragment: `... on TeamsPropsHistory { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    UpdateTacticPayload: {
      teamByTeamId: {
        fragment: `... on UpdateTacticPayload { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    UpdateTeamsPropPayload: {
      teamByTeamId: {
        fragment: `... on UpdateTeamsPropPayload { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
        },
      },
    },
    UpdateTrainingPayload: {
      teamByTeamId: {
        fragment: `... on UpdateTrainingPayload { teamId }`,
        resolve(parent, args, context, info) {
          return teamByTeamId(parent, args, context, info, horizonRemoteSchema);
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
      createBid: async (parent, args, context, info) => {
        return createBidResolver(parent, args, context, info, horizonRemoteSchema, web3);
      },
      setMessage: setMessageResolver,
      setBroadcastMessage: setBroadcastMessageResolver,
      setMailboxStart: setMailboxStartResolver,
      setMessageRead: setMessageReadResolver,
      setLastTimeLoggedIn: setLastTimeLoggedInResolver,
    },
    Query: {
      getMessages: getMessagesResolver,
      getNumUnreadMessages: getNumUnreadMessagesResolver,
      getLastTimeLoggedIn: getLastTimeLoggedInResolver,
    },
  };
};

module.exports = resolvers;
