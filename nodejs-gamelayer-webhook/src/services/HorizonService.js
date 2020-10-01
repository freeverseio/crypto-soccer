const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async processAuction({ auctionId }) {
    const query = gql`
    {
        processAuction(id: "${auctionId}")
    }
    `;
  }
}

module.exports = new HorizonService();
