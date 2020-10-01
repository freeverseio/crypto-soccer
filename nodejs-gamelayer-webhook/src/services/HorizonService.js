const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async checkOrderState({ teamId }) {
    const query = gql`
    {
        checkOrderState(teamId: "${teamId}")
    }
    `;
  }
}

module.exports = new HorizonService();
