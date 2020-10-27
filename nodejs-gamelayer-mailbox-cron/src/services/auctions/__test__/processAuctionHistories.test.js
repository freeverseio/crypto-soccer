const processAuctionHistories = require('../processAuctionHistories');
const HorizonService = require('../../HorizonService');
const processPayingAuction = require('../processPayingAuction');
const processWithdrawableBySellerAuction = require('../processWithdrawableBySellerAuction');
const selectLastChecked = require('../../../repositories/selectLastChecked');
const updateLastChecked = require('../../../repositories/updateLastChecked');
const processAcceptedBids = require('../../bids/processAcceptedBids');

jest.mock('../../HorizonService.js', () => ({
  getLastAuctionsHistories: jest.fn().mockReturnValue([
    {
      insertedAt: '2020-09-29T09:58:09.431379+00:00',
      playerId: '2748779076711',
      currencyId: 1,
      price: '50',
      rnd: '321960049',
      validUntil: '1601373789',
      signature:
        '22795665491fb888d3ccb0be12f4c27357e7cd6a200875a608721bcd580b89ff0acb98b7d11f5424540da9573d5c64191a270ad4ba9243639810b9a261f8c0451b',
      state: 'WITHADRABLE_BY_SELLER',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      id: '591e7084cd10ffd282783744baa90a72a6ac949ecaffe2cd15f2120135750441',
      offerValidUntil: '0',
    },
    {
      insertedAt: '2020-09-29T09:58:55.710169+00:00',
      playerId: '2748779076711',
      currencyId: 1,
      price: '50',
      rnd: '321960049',
      validUntil: '1601373789',
      signature:
        '22795665491fb888d3ccb0be12f4c27357e7cd6a200875a608721bcd580b89ff0acb98b7d11f5424540da9573d5c64191a270ad4ba9243639810b9a261f8c0451b',
      state: 'PAYING',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      id: '591e7084cd10ffd282783744baa90a72a6ac949ecaffe2cd15f2120135750441',
      offerValidUntil: '0',
    },
    {
      insertedAt: '2020-09-29T09:59:55.710169+00:00',
      playerId: '2748779076711',
      currencyId: 1,
      price: '50',
      rnd: '321960049',
      validUntil: '1601373789',
      signature:
        '22795665491fb888d3ccb0be12f4c27357e7cd6a200875a608721bcd580b89ff0acb98b7d11f5424540da9573d5c64191a270ad4ba9243639810b9a261f8c0451b',
      state: 'PAYING',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      id: '591e7084cd10ffd282783744baa90a72a6ac949ecaffe2cd15f2120135750441',
      offerValidUntil: '0',
    },
  ]),
}));

jest.mock('../processPayingAuction', () => jest.fn());
jest.mock('../processWithdrawableBySellerAuction', () => jest.fn());
jest.mock('../../../repositories/selectLastChecked', () => jest.fn());
jest.mock('../../../repositories/updateLastChecked', () => jest.fn());
jest.mock('../../bids/processAcceptedBids', () => jest.fn());

afterEach(() => {
  jest.clearAllMocks();
});

test('processAuctionHistories works correctly', async () => {
  await processAuctionHistories();

  expect(selectLastChecked).toHaveBeenCalled();
  expect(updateLastChecked).toHaveBeenCalledWith({
    entity: 'auction',
    lastChecked: '2020-09-29T09:59:55.710169+00:00',
  });
  expect(HorizonService.getLastAuctionsHistories).toHaveBeenCalledTimes(1);
  expect(processPayingAuction).toHaveBeenCalledTimes(2);
  expect(processWithdrawableBySellerAuction).toHaveBeenCalledTimes(1);
  expect(processAcceptedBids).toHaveBeenCalledTimes(1);
});
