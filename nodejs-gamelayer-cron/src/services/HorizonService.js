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

  async getLastAuctionsHistories({ lastChecked }) {
    const query = gql`
      {
        allAuctionsHistories(
          filter: {
            insertedAt: { greaterThan: "${lastChecked}" }
          },
          orderBy: INSERTED_AT_ASC
        ) {
          nodes {
            insertedAt
            playerId
            currencyId
            price
            rnd
            validUntil
            signature
            state
            stateExtra
            seller
            id
            offerValidUntil
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

  async getBidsByAuctionId({ auctionId }) {
    const query = gql`
      {
        allBids(
          condition: {
            auctionId: "${auctionId}"
          }
        ) {
          nodes {
            extraPrice
            rnd
            teamId
            signature
            state
            stateExtra
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

  async getPayingBids() {
    const query = gql`
      {
        allBids(condition: { state: PAYING }) {
          nodes {
            extraPrice
            rnd
            teamId
            signature
            state
            stateExtra
            paymentDeadline
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
