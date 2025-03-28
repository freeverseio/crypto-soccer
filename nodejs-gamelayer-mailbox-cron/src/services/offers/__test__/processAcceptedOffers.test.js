const HorizonService = require('../../HorizonService');
const GamelayerService = require('../../GamelayerService');
const processAcceptedOffers = require('../processAcceptedOffers');

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
  getInfoFromTeamId: jest.fn().mockReturnValue({
    teamId: '2748779069857',
    name: 'Magicians Plus',
    managerName: 'asdas',
  }),
}));

jest.mock('../../HorizonService.js', () => ({
  getInfoFromPlayerId: jest.fn().mockReturnValue({
    teamId: '2748779069626',
    name: 'joreg',
  }),
}));

afterEach(() => {
  jest.clearAllMocks();
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
  state: 'ACCEPTED',
  stateExtra: '',
  seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
  buyer: '0xaC347a9Fa330c6c23136F1460086D436ed55a3f8',
  auctionId: 'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
};

test('processAcceptedOffers works correctly', async () => {
  await processAcceptedOffers({ offerHistory });
  expect(HorizonService.getInfoFromPlayerId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.getInfoFromTeamId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(1);
});
