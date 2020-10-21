const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');
const logger = require('../logger');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async processAuction({ auctionId }) {
    const query = gql`
    mutation {
        processAuction(input: { id: "${auctionId}"})
    }
    `;
    logger.debug(`Sending processAuction(id:${auctionId}) mutation...`);
    await request(this.endpoint, query);
  }
}

module.exports = new HorizonService();
