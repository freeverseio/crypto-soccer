const { request, gql } = require('graphql-request');
const { gamelayerConfig } = require('../config.js');

class GamelayerService {
  constructor() {
    this.endpoint = gamelayerConfig.url;
  }

  async setMessage({
    destinatary,
    category,
    auctionId,
    text,
    customImageUrl,
    metadata,
    title,
  }) {
    const query = gql`
      mutation {
        setMessage(input: {
        destinatary: "${destinatary}",
        category: "${category}",
        auctionId: "${auctionId}",
        title: "${title}",
        text: "${text}",
        customImageUrl: "${customImageUrl}",
        metadata: "${metadata}",
        })
      }
    `;
    const result = await request(this.endpoint, query);

    return result;
  }

  async getInfoFromTeamId({ teamId }) {
    const query = gql`
      {
        allTeams(condition: { teamId: "${teamId}" }) {
          nodes {
            teamId
            name
            managerName
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allTeams && result.allTeams.nodes
      ? result.allTeams.nodes[0]
      : [];
  }

  async getMessages({ teamId, auctionId }) {
    const query = gql`
      {
        getMessages(teamId: "${teamId}", auctionId: "${auctionId}") {
          nodes {
            text
            auctionId
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.getMessages && result.getMessages.nodes
      ? result.getMessages.nodes
      : [];
  }
}

module.exports = new GamelayerService();
