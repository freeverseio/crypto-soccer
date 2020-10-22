const GamelayerService = require('../../GamelayerService');
const HorizonService = require('../../HorizonService');
const processStartedOffers = require('../processStartedOffers');
const notifyNewHigherOffer = require('../notifyNewHigherOffer');

jest.mock('../../HorizonService.js', () => ({
  getInfoFromTeamId: jest.fn().mockReturnValue({
    teamId: '2748779069857',
    name: 'Magicians Plus',
    managerName: 'asdas',
  }),
  getInfoFromPlayerId: jest.fn().mockReturnValue({
    teamId: '2748779069626',
    name: 'joreg',
  }),
  getOfferersByPlayerId: jest.fn().mockReturnValue([
    {
      insertedAt: '2020-09-29T10:12:51.070996+00:00',
      auctionId:
        'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
      playerId: '2748779076705',
      currencyId: 1,
      price: '50',
      rnd: '375503914',
      validUntil: '1601374670',
      signature:
        '6beb027ad523f3aa2a6629ab6153aa2a6a94f89597d31fbeef4e0170e7f0e0f36ee7eda8b2debcfcf90362ce6f6fa88b2c7a1850f336dc1911797e8e3a6270b61c',
      state: 'STARTED',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      buyer: '0xaC347a9Fa330c6c23136F1460086D436ed55a3f8',
      buyerTeamId: '2748779069845',
    },
    {
      insertedAt: '2020-09-29T10:12:51.070996+00:00',
      auctionId:
        'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
      playerId: '2748779076705',
      currencyId: 1,
      price: '50',
      rnd: '375503914',
      validUntil: '1601374670',
      signature:
        '6beb027ad523f3aa2a6629ab6153aa2a6a94f89597d31fbeef4e0170e7f0e0f36ee7eda8b2debcfcf90362ce6f6fa88b2c7a1850f336dc1911797e8e3a6270b61c',
      state: 'STARTED',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      buyer: '0xaC347a9Fa330c6c23136F1460086D436ed55a3f9',
      buyerTeamId: '2748779069846',
    },
  ]),
}));

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
}));

jest.mock('../notifyNewHigherOffer', () => jest.fn());

afterEach(() => {
  jest.clearAllMocks();

  HorizonService.getOfferersByPlayerId.mockReturnValueOnce([
    {
      insertedAt: '2020-09-29T10:12:51.070996+00:00',
      auctionId:
        'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
      playerId: '2748779076705',
      currencyId: 1,
      price: '50',
      rnd: '375503914',
      validUntil: '1601374670',
      signature:
        '6beb027ad523f3aa2a6629ab6153aa2a6a94f89597d31fbeef4e0170e7f0e0f36ee7eda8b2debcfcf90362ce6f6fa88b2c7a1850f336dc1911797e8e3a6270b61c',
      state: 'STARTED',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      buyer: '0xaC347a9Fa330c6c23136F1460086D436ed55a3f9',
      buyerTeamId: '2748779069846',
    },
  ]);
});

const offerHistory = {
  insertedAt: '2020-09-29T10:12:51.070996+00:00',
  playerId: '2748779076705',
  currencyId: 1,
  price: '50',
  rnd: '375503914',
  validUntil: '1601374670',
  signature:
    '6beb027ad523f3aa2a6629ab6153aa2a6a94f89597d31fbeef4e0170e7f0e0f36ee7eda8b2debcfcf90362ce6f6fa88b2c7a1850f336dc1911797e8e3a6270b61c',
  state: 'STARTED',
  stateExtra: '',
  seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
  buyer: '0xaC347a9Fa330c6c23136F1460086D436ed55a3f8',
  auctionId: 'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
};

test('processStartedOffers works correctly more than 1 offerer', async () => {
  await processStartedOffers({ offerHistory });

  expect(HorizonService.getOfferersByPlayerId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromPlayerId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromTeamId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(0);
  expect(notifyNewHigherOffer).toHaveBeenCalledTimes(3);
});

test('processStartedOffers works correctly with 1 offerer', async () => {
  await processStartedOffers({ offerHistory });

  expect(HorizonService.getOfferersByPlayerId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromPlayerId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromTeamId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(1);
  expect(notifyNewHigherOffer).toHaveBeenCalledTimes(1);
});
