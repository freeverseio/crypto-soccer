const repositories = require("../repositories")
const HorizonService = require("../services/HorizonService.js")
const BidValidation = require("./BidValidation.js")
const Web3 = require('web3')

jest.mock('../repositories', () => ({
  selectTeamMaximumBid: jest.fn().mockReturnValue(25),
  updateTeamMaximumBid: jest.fn(),
}));

jest.mock('../services/HorizonService.js', () => ({
  getBidsPayed: jest.fn().mockReturnValue([{ extraPrice: 1, auctionByAuctionId: { price: 5 }}]),
  getAuction: jest.fn().mockReturnValue({ price: 10 }),
}))

afterEach(() => {
  jest.clearAllMocks()
})

test('allowed to bid', async () => {
  const web3 = new Web3('')
  bidValidation = new BidValidation({ teamId: '234', rnd: 12345, auctionId: '555', extraPrice: 10, signature: '134ab', web3 });
  const result = await bidValidation.isAllowedToBid()

  expect(result).toBe(true)
  expect(repositories.selectTeamMaximumBid).toHaveBeenCalledWith({ teamId: '234' })
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' })
  expect(repositories.updateTeamMaximumBid).toHaveBeenCalledTimes(0)
  expect(HorizonService.getBidsPayed).toHaveBeenCalledTimes(0)

})

test('not allowed to bid because 0 is set', async () => {
  const web3 = new Web3('')
  repositories.selectTeamMaximumBid.mockReturnValue(0)

  bidValidation = new BidValidation({ teamId: '234', rnd: 12345, auctionId: '555', extraPrice: 10, signature: '134ab', web3 });
  const result = await bidValidation.isAllowedToBid()

  expect(result).toBe(false)
  expect(repositories.selectTeamMaximumBid).toHaveBeenCalledWith({ teamId: '234' })
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(0)
  expect(repositories.updateTeamMaximumBid).toHaveBeenCalledTimes(0)
  expect(HorizonService.getBidsPayed).toHaveBeenCalledTimes(0)

})

test('not allowed to bid because current maximum bid and computed maximum bid are lower than total price', async () => {
  const web3 = new Web3('')
  repositories.selectTeamMaximumBid.mockReturnValue(11)

  bidValidation = new BidValidation({ teamId: '234', rnd: 12345, auctionId: '555', extraPrice: 10, signature: '134ab', web3 });
  const result = await bidValidation.isAllowedToBid()

  expect(result).toBe(false)
  expect(repositories.selectTeamMaximumBid).toHaveBeenCalledWith({ teamId: '234' })
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' })
  expect(HorizonService.getBidsPayed).toHaveBeenCalledWith({ teamId: '234' })
  expect(repositories.updateTeamMaximumBid).toHaveBeenCalledWith({ teamId: '234', teamMaximumBid: 9 })
})

test('allowed to bid because current maximum bid is lower than total price but computed is greater', async () => {
  const web3 = new Web3('')
  repositories.selectTeamMaximumBid.mockReturnValue(11)
  HorizonService.getBidsPayed.mockReturnValue([{ extraPrice: 1, auctionByAuctionId: { price: 50 }}])

  bidValidation = new BidValidation({ teamId: '234', rnd: 12345, auctionId: '555', extraPrice: 10, signature: '134ab', web3 });
  const result = await bidValidation.isAllowedToBid()

  expect(result).toBe(true)
  expect(repositories.selectTeamMaximumBid).toHaveBeenCalledWith({ teamId: '234' })
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' })
  expect(HorizonService.getBidsPayed).toHaveBeenCalledWith({ teamId: '234' })
  expect(repositories.updateTeamMaximumBid).toHaveBeenCalledWith({ teamId: '234', teamMaximumBid: 76.5 })
})