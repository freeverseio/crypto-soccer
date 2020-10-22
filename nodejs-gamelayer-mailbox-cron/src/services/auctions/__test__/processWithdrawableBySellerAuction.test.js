const HorizonService = require('../../HorizonService');
const GamelayerService = require('../../GamelayerService');
const processWithdrawableBySellerAuction = require('../processWithdrawableBySellerAuction');
const getTeamIdFromAuctionSeller = require('../getTeamIdFromAuctionSeller');

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
}));

jest.mock('../../HorizonService', () => ({
  getInfoFromTeamId: jest.fn().mockReturnValue({
    teamId: '2748779069857',
    name: 'Magicians Plus',
    managerName: 'asdas',
  }),
  getInfoFromPlayerId: jest.fn().mockReturnValue({
    teamId: '2748779069626',
  }),
  getPaidBidByAuctionId: jest.fn().mockReturnValue({
    extraPrice: 0,
    rnd: 0,
    teamId: '2748779069875',
    signature:
      '531a6e38203a3ba21a5a845b9287dedb985ca638a19c21f2f9e083a7ae5ac3210a0eaa6df4c573e1000b2c842b4dbff5606633ae1994f986fd533a6aea43daed1c',
    state: 'PAID',
    stateExtra: '',
  }),
}));

const auctionHistory = {
  insertedAt: '2020-09-29T09:58:09.431379+00:00',
  playerId: '2748779076711',
  currencyId: 1,
  price: '50',
  rnd: '321960049',
  validUntil: '1601373789',
  signature:
    '22795665491fb888d3ccb0be12f4c27357e7cd6a200875a608721bcd580b89ff0acb98b7d11f5424540da9573d5c64191a270ad4ba9243639810b9a261f8c0451b',
  state: 'STARTED',
  stateExtra: '',
  seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
  id: '591e7084cd10ffd282783744baa90a72a6ac949ecaffe2cd15f2120135750441',
  offerValidUntil: '0',
};

jest.mock('../getTeamIdFromAuctionSeller', () =>
  jest.fn().mockReturnValue('234324234')
);

afterEach(() => {
  jest.clearAllMocks();
});

test('processWithdrawableBySellerAuction works correctly', async () => {
  await processWithdrawableBySellerAuction({ auctionHistory });
  expect(HorizonService.getPaidBidByAuctionId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromPlayerId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromTeamId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(1);
});
