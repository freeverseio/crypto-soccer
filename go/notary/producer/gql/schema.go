//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package gql

const Schema = ` 
	input CreatePutPlayerForSaleInput {
  		signature: String!
  		playerId: String!
  		currencyId: Int!
  		price: Int!
  		validUntil: String!
  		rnd: Int!
	}

	input CancelAuctionInput {
  		signature: String!
		auctionId: ID!
	}

	input CreateBidInput {
  		signature: String!
		auctionId: ID!
  		extraPrice: Int!
  		rnd: Int!
  		teamId: String!
	}

	input GetWorldPlayersInput {
		signature: String!
		teamId: ID!
	}

	input SubmitPlayStorePlayerPurchaseInput {
		signature: String!
		receipt: String!
		playerId: ID!
		teamId: ID!
	}

	input DismissPlayerInput {
		signature: String!
		validUntil: String!
		playerId: ID!
	}

	input CompletePlayerTransitInput {
		playerId: ID!
	}

	input CreateOfferInput {
		signature: String!
		playerId: String!
		currencyId: Int!
		price: Int!
		validUntil: String!
		rnd: Int!
		buyerTeamId: String!
	}

	input CancelOfferInput {
		signature: String!
		offerId:   ID!
	}

	input AcceptOfferInput {
		signature: String!
		playerId: String!
		currencyId: Int!
		price: Int!
		validUntil: String!
		offerValidUntil: String!
		rnd: Int!
  }

	type WorldPlayer {
		playerId: ID!
		name: String!
		dayOfBirth: Int! 
		preferredPosition: String!
		defence: Int!
		speed: Int!
		pass: Int!
		shoot: Int!
		endurance: Int!
		potential: Int! 
		validUntil: String!
		countryOfBirth: String!
		race: String!
		productId: String!
	}

	type Query {
		getWorldPlayers(input: GetWorldPlayersInput!): [WorldPlayer]! 
	}

	type Mutation {
		createOffer(input: CreateOfferInput!): ID!
		createAuctionFromPutForSale(input: CreatePutPlayerForSaleInput!): ID!
		acceptOffer(input: AcceptOfferInput!): ID!
		createBid(input: CreateBidInput!): ID!
		cancelAuction(input: CancelAuctionInput!): ID!
		submitPlayStorePlayerPurchase(input: SubmitPlayStorePlayerPurchaseInput!): ID!
		dismissPlayer(input: DismissPlayerInput!): ID!
		completePlayerTransit(input: CompletePlayerTransitInput!): ID!
		cancelOffer(input: CancelOfferInput!): ID!
	}
`
