const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async getOfferersByPlayerId({ playerId }) {
    const query = gql`
      {
        allOffersHistories(condition: { state: STARTED, playerId: ${playerId} }) {
          nodes {
            insertedAt
            auctionId
            playerId
            currencyId
            price
            rnd
            validUntil
            signature
            state
            stateExtra
            seller
            buyer
            buyerTeamId
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result &&
      result.allOffersHistories &&
      result.allOffersHistories.nodes
      ? result.allOffersHistories.nodes
      : [];
  }

  async getLastOfferHistories({ lastChecked }) {
    const query = gql`
      {
        allOffersHistories(
          filter: {
            insertedAt: { greaterThan: "${lastChecked}" }
          },
          orderBy: INSERTED_AT_ASC
        ) {
          nodes {
            insertedAt
            auctionId
            playerId
            currencyId
            price
            rnd
            validUntil
            signature
            state
            stateExtra
            seller
            buyer
            buyerTeamId
          }
        }
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

module.exports = new HorizonService();
