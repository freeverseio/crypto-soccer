const Web3 = require('web3');
const rewire = require('rewire');
const createOffer = rewire('../createOfferResolver');

test('Compute auction id', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  auctionId = createOffer.__get__('computeAuctionId')({
    currencyId: 1,
    price: 41234,
    rnd: 42321,
    playerId: '25723578238440869144533393071649442553899076447028039543423578',
    validUntil: '2000000000',
    web3,
  });
  expect(auctionId).toBe('a51ac22be9bcc29236b88373d494c30c5dff2108059aa14c00e4c2a1d4061eba');
});
