const HorizonService = require('../../HorizonService');
const GamelayerService = require('../../GamelayerService');
const getTeamIdFromAuctionSeller = require('../../auctions/getTeamIdFromAuctionSeller');
const processAcceptedBids = require('../processAcceptedBids');

jest.mock('../../HorizonService.js', () => {
  const dayjs = require('dayjs');

  return {
    getLastAcceptedBidsHistories: jest.fn().mockReturnValue([
      {
        insertedAt: '2020-09-29T09:58:37.827772+00:00',
        rnd: 149275350,
        signature:
          '87ac4e518f4648961fea5b4762259e8bab272a256710a51561840138095d553d5ea3ca73b8df8989c89d8d74b6ddd973d5fbca0ad45cca78e322c8e1e15d14361b',
        state: 'ACCEPTED',
        stateExtra: '',
        auctionId:
          '591e7084cd10ffd282783744baa90a72a6ac949ecaffe2cd15f2120135750441',
        extraPrice: 0,
        teamId: '2748779069845',
        paymentUrl: '',
        paymentId: '',
        paymentDeadline: '0',
      },
      {
        insertedAt: '2020-09-29T10:13:47.226091+00:00',
        rnd: 0,
        signature:
          '6beb027ad523f3aa2a6629ab6153aa2a6a94f89597d31fbeef4e0170e7f0e0f36ee7eda8b2debcfcf90362ce6f6fa88b2c7a1850f336dc1911797e8e3a6270b61c',
        state: 'ACCEPTED',
        stateExtra: '',
        auctionId:
          'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
        extraPrice: 0,
        teamId: '2748779069845',
        paymentUrl: '',
        paymentId: '',
        paymentDeadline: '0',
      },
      {
        insertedAt: '2020-09-29T12:45:29.147433+00:00',
        rnd: 0,
        signature:
          'd175435f7dd55ee4edb8388b112df707cc7d9cc12da43486870b679a277dc6023a141ecc59b77ea753b35916e8e57d0c027ba4e84e7110dbe3a0c2ee39f3e8d81b',
        state: 'ACCEPTED',
        stateExtra: '',
        auctionId:
          '3829046fe1b9f3a74ad3ee8c28edf2c79761737e785dc22d7339e756449b280a',
        extraPrice: 0,
        teamId: '2748779069852',
        paymentUrl: '',
        paymentId: '',
        paymentDeadline: '0',
      },
    ]),

    getInfoFromPlayerId: jest.fn().mockReturnValue({
      teamId: '2748779069626',
    }),
    getAuction: jest.fn().mockReturnValue({
      playerId: '2748779076846',
      price: '100',
      state: 'ASSET_FROZEN',
      playerByPlayerId: {
        name: 'Hector Galindo',
      },
      bidsByAuctionId: {
        nodes: [
          {
            teamId: '2748779069852',
          },
          {
            teamId: '2748779069852',
          },
          {
            teamId: '2748779069852',
          },
          {
            teamId: '2748779069852',
          },
          {
            teamId: '2748779069852',
          },
        ],
      },
      seller: '0x9f46F66b079d469920e4e72a99ef42D8A3447C10',
    }),
  };
});

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
  getInfoFromTeamId: jest.fn().mockReturnValue({
    teamId: '2748779069857',
    name: 'Magicians Plus',
    managerName: 'asdas',
  }),
}));

jest.mock('../../auctions/getTeamIdFromAuctionSeller', () =>
  jest.fn().mockReturnValue('134324')
);

afterEach(() => {
  jest.clearAllMocks();
});

test('processAcceptedBids works correctly', async () => {
  await processAcceptedBids({ lastChecked: '2020-06-01T16:00:00.000Z' });

  expect(HorizonService.getLastAcceptedBidsHistories).toHaveBeenCalledTimes(1);
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(3);
  expect(getTeamIdFromAuctionSeller).toHaveBeenCalledTimes(3);
  expect(GamelayerService.getInfoFromTeamId).toHaveBeenCalledTimes(3);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(13);
});
