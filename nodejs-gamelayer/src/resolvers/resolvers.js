const Web3 = require('web3')
const { selectPlayerName, selectTeamName, selectTeamManagerName, updatePlayerName, updateTeamName, updateTeamManagerName } = require("../repositories");
const { TeamValidation, PlayerValidation } = require("../validations");

const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")

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
      setGamePlayerName: async (_, { input: { playerId, name, signature } }) => {
          const playerValidation = new PlayerValidation({ playerId, name, signature, web3 })
          const isSignerOwner = await playerValidation.isSignerOwner()
          
          if(isSignerOwner) {
            await updatePlayerName({ playerId, playerName: name })
            return playerId 
          } else {
            return "Signer is not the player owner"
          }
        },
      setGameTeamName: async (_, { input: { teamId, name, signature } }) => {
        const teamValidation = new TeamValidation({ teamId, name, signature, web3 })
        const isSignerOwner = await teamValidation.isSignerOwner()
    
        if(isSignerOwner) {
          await updateTeamName({ teamId, teamName: name })
          return teamId 
        } else {
          return "Signer is not the team owner"
        }
      },
      setGameTeamManagerName: async (_, { input: { teamId, name, signature } }) => {
        const teamValidation = new TeamValidation({ teamId, name, signature, web3 })
        const isSignerOwner = await teamValidation.isSignerOwner()
    
        if(isSignerOwner) {
          await updateTeamManagerName({ teamId, teamManagerName: name })
          return teamId 
        } else {
          return "Signer is not the team owner"
        }
      },
    }
  };

  module.exports = resolvers;