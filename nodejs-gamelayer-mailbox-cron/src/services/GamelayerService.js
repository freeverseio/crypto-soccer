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
}

module.exports = new GamelayerService();
