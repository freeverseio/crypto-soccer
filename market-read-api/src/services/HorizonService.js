const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async getAllOffers() {
    const query = gql`
      {
        allOffers {
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
            teamByBuyerTeamId {
              teamId
              name
              managerName
            }
            playerByPlayerId {
              name
              playerId
              teamId
              defence
              speed
              pass
              shoot
              endurance
              shirtNumber
              preferredPosition
              potential
              dayOfBirth
              countryOfBirth
              race
            }
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allOffers && result.allOffers.nodes ? result.allOffers.nodes : [];
  }

  async getAllBids() {
    const query = gql`
      {
        allBids {
          nodes {
            auctionId
            extraPrice
            rnd
            teamId
            signature
            state
            stateExtra
            teamByTeamId {
              teamId
              name
              managerName
            }
            auctionByAuctionId {
              id
              playerId
              currencyId
              price
              rnd
              validUntil
              signature
              state
              stateExtra
              seller
              offerValidUntil
            }
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allBids && result.allBids.nodes ? result.allBids.nodes : [];
  }

  async getAllAuctions() {
    const query = gql`
    {
      allAuctions 
      (	filter: {
          and : [
            { 
              state: { 
                equalTo: STARTED
              } 
            },
            { 
              state: { 
                equalTo: ASSET_FROZEN
              } 
            }
          ]
        }
      )
      {
        nodes {
          id
          playerId
          price
          validUntil
          state
          seller
          bidsByAuctionId(orderBy: EXTRA_PRICE_DESC) {
            nodes {
                teamId
                extraPrice
                state
                teamByTeamId {
                    name
                }
                paymentDeadline
            }
          }
          playerByPlayerId {
            playerId
            name
            teamId
            defence
            speed
            pass
            shoot
            endurance
            preferredPosition
            potential
            dayOfBirth
            shirtNumber
            countryOfBirth
            teamByTeamId {
              name
              managerName
            }
          }
        }
      }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allAuctions && result.allAuctions.nodes ? result.allAuctions.nodes : [];
  }

  async getAuction({ id }) {
    const query = gql`
    {
      allAuctions(condition: {id: "${id}"})
      {
        nodes {
          id
          playerId
          price
          validUntil
          state
          seller
          bidsByAuctionId(orderBy: EXTRA_PRICE_DESC) {
            nodes {
                teamId
                extraPrice
                state
                teamByTeamId {
                    name
                }
                paymentUrl
                paymentDeadline
            }
          }
          playerByPlayerId {
            playerId
            name
            teamId
            defence
            speed
            pass
            shoot
            endurance
            preferredPosition
            potential
            dayOfBirth
            countryOfBirth
            teamByTeamId {
              name
              managerName
            }
          }
        }
      }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allAuctions && result.allAuctions.nodes ? result.allAuctions.nodes : [];
  }
}

module.exports = new HorizonService();
