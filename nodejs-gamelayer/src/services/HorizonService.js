const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async getTeamOwner({ teamId }) {
    const query = gql`
    {
        teamByTeamId(teamId: "${teamId}"){ owner }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.teamByTeamId && result.teamByTeamId.owner ? result.teamByTeamId.owner : '';
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
    `;
    const result = await request(this.endpoint, query);

    return result &&
      result.playerByPlayerId &&
      result.playerByPlayerId.teamByTeamId &&
      result.playerByPlayerId.teamByTeamId.owner
      ? result.playerByPlayerId.teamByTeamId.owner
      : '';
  }

  async getAllTeamIds() {
    const query = gql`
      {
        allTeams {
          nodes {
            teamId
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allTeams && result.allTeams.nodes ? result.allTeams.nodes : [];
  }

  async getAuction({ auctionId }) {
    const query = gql`
    {
        auctionById(id: "${auctionId}"){ 
            price
        }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.auctionById ? result.auctionById : {};
  }

  async getBidsPayed({ teamId }) {
    const query = gql`
    {
      allBids(condition: { teamId: "${teamId}", state: PAID }){
        nodes {
          extraPrice
          auctionByAuctionId{
            id
            state
            price
          }
        }
      }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allBids && result.allBids.nodes ? result.allBids.nodes : [];
  }
}

module.exports = new HorizonService();
