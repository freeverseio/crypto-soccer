const schema = `
type Auction {
    id: String!
    playerId: String!
    currencyId: Int!
    price: BigInt!
    rnd: BigInt!
    validUntil: BigInt!
    signature: String!
    state: AuctionState!
    stateExtra: String!
    paymentUrl: String!
    seller: String!
    offerValidUntil: BigInt!
  }
  
  enum AuctionState {
    STARTED
    FAILED
    CANCELLED
    ENDED
    ASSET_FROZEN
    PAYING
    WITHADRABLE_BY_SELLER
    WITHADRABLE_BY_BUYER
    VALIDATION
  }
  
  type Bid {
    auctionId: String!
    extraPrice: Int!
    rnd: Int!
    teamId: String!
    signature: String!
    state: BidState!
    stateExtra: String!
    paymentId: String!
    paymentUrl: String!
    paymentDeadline: String!
    auctionByAuctionId: Auction
    teamByTeamId: Team
  }
  
  enum BidState {
    ACCEPTED
    PAYING
    PAID
    FAILED
  }
  
  enum OfferState {
    STARTED
    FAILED
    CANCELLED
    ENDED
    ACCEPTED
  }
  
  # A signed eight-byte integer. The upper big integer values are greater than the
  # max value for a JavaScript number. Therefore all big integers will be output as
  # strings and not numbers.
  scalar BigInt
  
  type Team {
    teamId: String!
    name: String!
    managerName: String!
  }
  
  type Player {
    name: String!
    playerId: String!
    teamId: String!
    defence: Int!
    speed: Int!
    pass: Int!
    shoot: Int!
    endurance: Int!
    shirtNumber: Int!
    preferredPosition: String!
    potential: Int!
    dayOfBirth: Int!
    countryOfBirth: String!
    race: String!
  }
  
  type Offer {
    auctionId: String!
    playerId: String!
    currencyId: Int!
    price: BigInt!
    rnd: BigInt!
    validUntil: BigInt!
    signature: String!
    state: OfferState!
    stateExtra: String!
    seller: String!
    buyer: String!
    buyerTeamId: String!
    teamByBuyerTeamId: Team
    playerByPlayerId: Player
  }
  
  type Query {
    allAuctions: [Auctions]
    allBids: [Bids]
    allOffers: [Offers]
  }  
`;

module.exports = schema;