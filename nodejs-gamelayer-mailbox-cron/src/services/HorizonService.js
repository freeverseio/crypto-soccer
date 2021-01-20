const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');
const logger = require('../logger');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async getOfferersByPlayerId({ playerId }) {
    const query = gql`
      {
        allOffers(condition: { state: STARTED, playerId: "${playerId}" }) {
          nodes {
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

    return result && result.allOffers && result.allOffers.nodes
      ? result.allOffers.nodes
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
      result.allAuctionsHistories &&
      result.allAuctionsHistories.nodes
      ? result.allAuctionsHistories.nodes
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
            paymentDeadline
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allBids && result.allBids.nodes
      ? result.allBids.nodes
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
            auctionId
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allBids && result.allBids.nodes
      ? result.allBids.nodes
      : [];
  }

  async getTeamIdsFromOwner({ owner }) {
    const query = gql`
      {
        allTeams(condition: { owner: "${owner}" }) {
          nodes {
            teamId
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allTeams && result.allTeams.nodes
      ? result.allTeams.nodes
      : [];
  }

  async getPlayerHistoriesLast30BlockNumberTeams({ playerId }) {
    const query = gql`
      {
        allPlayersHistories(
          first: 30
          condition: { playerId: "${playerId}" }
          orderBy: BLOCK_NUMBER_DESC
        ) {
          nodes {
            teamId
            blockNumber
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result &&
      result.allPlayersHistories &&
      result.allPlayersHistories.nodes
      ? result.allPlayersHistories.nodes
      : [];
  }

  async getInfoFromPlayerId({ playerId }) {
    const query = gql`
      {
        allPlayers(condition: { playerId: "${playerId}" }) {
          nodes {
            teamId
            name
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allPlayers && result.allPlayers.nodes
      ? result.allPlayers.nodes[0]
      : [];
  }

  async getPaidBidByAuctionId({ auctionId }) {
    const query = gql`
      {
        allBids(
          condition: {
            auctionId: "${auctionId}",
            state: PAID
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

    return result && result.allBids && result.allBids.nodes
      ? result.allBids.nodes[0]
      : [];
  }
  async getAuction({ auctionId }) {
    const query = gql`
      {
        allAuctions(
          condition: {
            id: "${auctionId}"
          }
        ) {
          nodes {
            playerId
            price
            state
            playerByPlayerId {
              name
            }
            bidsByAuctionId{
              nodes {
                teamId
                extraPrice
              }
            }
            seller
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allAuctions && result.allAuctions.nodes
      ? result.allAuctions.nodes[0]
      : [];
  }

  async getLastAcceptedBidsHistories({ lastChecked }) {
    const query = gql`
      {
        allBidsHistories(
          condition: { state: ACCEPTED },
          filter: {
            insertedAt: { greaterThan: "${lastChecked}" }
          },
          orderBy: INSERTED_AT_ASC
        ) {
          nodes {
            insertedAt
            rnd
            signature
            state
            stateExtra
            auctionId
            extraPrice
            teamId
						paymentUrl
            paymentId
            paymentDeadline
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allBidsHistories && result.allBidsHistories.nodes
      ? result.allBidsHistories.nodes
      : [];
  }

  async getLastUnpayments() {
    const query = gql`
      {
        allUnpayments(condition: { notified: false }) {
          nodes {
            owner
            timeOfUnpayment
            notified
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allUnpayments && result.allUnpayments.nodes
      ? result.allUnpayments.nodes
      : [];
  }

  async getUnpaymentsByOwner({ owner }) {
    const query = gql`
      {
        allUnpayments(condition: { owner: "${owner}" }) {
          nodes {
            id
            owner
            timeOfUnpayment
            notified
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allUnpayments && result.allUnpayments.nodes
      ? result.allUnpayments.nodes
      : [];
  }

  async setUnpaymentNotified({ unpayment }) {
    const query = gql`
    mutation {
        setUnpaymentNotified(input: { id: "${unpayment.id}"})
    }
    `;
    logger.debug(
      `Sending setUnpaymentNotified(id:${unpayment.id}) mutation...`
    );
    await request(this.endpoint, query);
  }
}

module.exports = new HorizonService();
