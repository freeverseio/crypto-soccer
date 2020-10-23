const schema = `
type Auction {
    price: BigInt!
    validUntil: BigInt!
    state: AuctionState!
    seller: String!
    offerValidUntil: BigInt!
    bidsByAuctionId: BidsConnection
    playerByPlayerId: Player
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
    extraPrice: Int!
    state: BidState!
    paymentDeadline: String!
    teamByTeamId: Team
  }

  type BidsConnection {
    nodes: [Bid]
  }

  enum BidState {
    ACCEPTED
    PAYING
    PAID
    FAILED
  }
  
  
  # A signed eight-byte integer. The upper big integer values are greater than the
  # max value for a JavaScript number. Therefore all big integers will be output as
  # strings and not numbers.
  scalar BigInt
  
  type Team {
    name: String!
    managerName: String!
  }
  
  type Player {
    name: String!
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
    teamByTeamId: Team
  }
  
  
  type Query {
    allAuctions: [Auction]
    getAuction(id: String!): [Auction]
  }  
`;

module.exports = schema;
