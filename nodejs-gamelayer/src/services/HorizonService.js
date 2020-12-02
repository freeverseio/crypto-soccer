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
        allTeams(filter: { owner: { notEqualTo: "0x0000000000000000000000000000000000000000" } }) {
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

  async getBidsPayedByTeamId({ teamId }) {
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

  async getBidsPayedByOwner({ owner }) {
    const teams = await this.getTeamsByOwner({ owner });
    const query = gql`
    {
      allBids(condition: { state: PAID }, filter: { teamId: { in: "${teams.map((t) => t.teamId)}"}}){
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

  async getTeamsByOwner({ owner }) {
    const query = gql`
    {
      allTeams(condition: { owner: "${owner}" }){
        nodes {
          teamId
        }
      }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allTeams && result.allTeams.nodes ? result.allTeams.nodes : [];
  }

  async getMatchesPlayedByTeamId({ teamId }) {
    const query = gql`
      {
        allMatchesHistories(
          condition: { matchDayIdx: 13, state: END }
          orderBy: BLOCK_NUMBER_DESC
          filter: {
            or: [{ homeTeamId: { equalTo: "${teamId}" } }, { visitorTeamId: { equalTo: "${teamId}" } }]
          }
        ) {
          nodes {
            blockNumber
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allMatchesHistories && result.allMatchesHistories.nodes
      ? result.allMatchesHistories.nodes
      : [];
  }

  async getPlayerTeamId({ playerId }) {
    const query = gql`
    {
        playerByPlayerId(playerId: "${playerId}"){ 
            teamId
              teamByTeamId {
              teamId
            }
        }
    }
    `;
    const result = await request(this.endpoint, query);

    return result &&
      result.playerByPlayerId &&
      result.playerByPlayerId.teamByTeamId &&
      result.playerByPlayerId.teamByTeamId.teamId
      ? result.playerByPlayerId.teamByTeamId.teamId
      : '';
  }

  async getEncodedSkillsByBlockNumberPlayerId({ playerId, blockNumber }) {
    const query = gql`
      {
        allPlayersHistories(condition: { playerId: "${playerId}", blockNumber: "${blockNumber}" }) {
          nodes {
            encodedSkills
          }
        }
      }
    `;

    const result = await request(this.endpoint, query);

    return result && result.allPlayersHistories && result.allPlayersHistories.nodes
      ? result.allPlayersHistories.nodes[0].encodedSkills
      : '';
  }

  async getPlayerHistory({ playerId, count }) {
    const query = gql`
      {
        playerByPlayerId(playerId: "${playerId}") {
          playersHistoriesByPlayerId(first: ${count}, orderBy: BLOCK_NUMBER_DESC) {
            nodes {
              encodedSkills
            }
          }
        }
      }
    `;

    const result = await request(this.endpoint, query);

    return result && result.playerByPlayerId && result.playerByPlayerId.playersHistoriesByPlayerId && result.playerByPlayerId.playersHistoriesByPlayerId.nodes
      ? result.playerByPlayerId.playersHistoriesByPlayerId
      : '';
  }
}

module.exports = new HorizonService();
