const HorizonService = require('../../HorizonService');
const GamelayerService = require('../../GamelayerService');
const processPayingAuction = require('../processPayingAuction');
const getTeamIdFromAuctionSeller = require('../getTeamIdFromAuctionSeller');
jest.mock('../../HorizonService.js', () => ({
  getBidsByAuctionId: jest.fn().mockReturnValue([
    {
      extraPrice: 0,
      rnd: 1688090333,
      teamId: '2748779069859',
      signature:
        '7aab2862668d5ee1195cd7cbedefea0a7f0158ee31ec09a2a81d094389afce1002778522f901a160feadd81be99c42570e109e46604fc5aa4d1b21b5952c20bf1c',
      state: 'PAYING',
      stateExtra: 'expired',
    },
    {
      extraPrice: 50,
      rnd: 1676548700,
      teamId: '2748779069860',
      signature:
        '799caa88c00a7b6d4c19db172e25d819c4b216cd4169124fb67fee0a7308f3ce57d18e9beb865fbf2d43946748dcf74011a6f6c9cca255f20eecf2d0e4174e0d1b',
      state: 'ACCEPTED',
      stateExtra: 'expired',
    },
    {
      extraPrice: 100,
      rnd: 1336658068,
      teamId: '2748779069859',
      signature:
        '9132d377118048e66c34e146b613b0fc4589edb79ef6e7a695c63e773d647cf74505df6e3defb17510b3d9d2a4699f892601a8f938897e75f5ea3ab7df8159ea1c',
      state: 'ACCEPTED',
      stateExtra: 'expired',
    },
  ]),
  getInfoFromPlayerId: jest.fn().mockReturnValue({
    teamId: '2748779069626',
    name: 'joreg',
  }),
}));

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
  getInfoFromTeamId: jest.fn().mockReturnValue({
    teamId: '2748779069857',
    name: 'Magicians Plus',
    managerName: 'asdas',
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

test('processPayingAuction works correctly', async () => {
  await processPayingAuction({ auctionHistory });

  expect(getTeamIdFromAuctionSeller).toHaveBeenCalledTimes(1);
  expect(HorizonService.getBidsByAuctionId).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromPlayerId).toHaveBeenCalledTimes(3);
  expect(GamelayerService.getInfoFromTeamId).toHaveBeenCalledTimes(4);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(4);
});
