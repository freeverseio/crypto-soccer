const { request, gql } = require('graphql-request')
const { horizonConfig } = require('../config.js')

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url
  }

  async getTeamOwner({ teamId }) {
    const query = gql`
    {
        teamByTeamId(teamId: "${teamId}"){ owner }
    }
    `
    const result = await request(this.endpoint, query)
    
    return result && result.teamByTeamId && result.teamByTeamId.owner ? result.teamByTeamId.owner : ''
  }

  async getPlayerOwner({ playerId }) {
    const query = gql`
    {
        playerByPlayerId(playerId: "${playerId}"){ 
            teamId
              teamByTeamId {
              owner
            }
        }
    }
    `
    const result = await request(this.endpoint, query)
    
    return result && result.playerByPlayerId && result.playerByPlayerId.teamByTeamId && result.playerByPlayerId.teamByTeamId.owner ? result.playerByPlayerId.teamByTeamId.owner : ''
  }
}

module.exports = new HorizonService();
