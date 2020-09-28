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
  }) {
    const query = gql`
      mutation {
        setMessage(input: {
        destinatary: "${destinatary}",
        category: "${category}",
        auctionId: "${auctionId}",
        text: "${text}",
        customImageUrl: "${customImageUrl}",
        metadata: "${metadata}",
        })
      }
    `;
    const result = await request(this.endpoint, query);

    return result &&
      result.allOffersHistories &&
      result.allOffersHistories.nodes
      ? result.allOffersHistories.nodes
      : [];
  }
}

module.exports = new GamelayerService();
