const HorizonService = require('../../HorizonService');
const GamelayerService = require('../../GamelayerService');
const processPendingPayingBids = require('../processPendingPayingBids');

jest.mock('../../HorizonService.js', () => {
  const dayjs = require('dayjs');

  return {
    getPayingBids: jest.fn().mockReturnValue([
      {
        extraPrice: 0,
        rnd: 0,
        teamId: '2748779069852',
        signature:
          '962b739c7a1534c58b85891454d14ef8c3ec1d87cb95b2daac98ff8bad15f0ef4c0ad45f3689774339723ed11040cc73ca9e21f850f8479c979f5a48c35e5ffa1b',
        state: 'PAYING',
        stateExtra: 'expired',
        paymentDeadline: dayjs().add(2, 'day').unix(),
      },
      {
        extraPrice: 0,
        rnd: 576423761,
        teamId: '2748779069570',
        signature:
          '9211cdc55106b3d45d852d9a21a888ce94f7df39291e6328d7ea884f78e6d6bd78b0300617d2ac690600cce4d8ca235a700b06d247631d3a0c7ae18a0fde26c71b',
        state: 'PAYING',
        stateExtra: 'expired',
        paymentDeadline: dayjs().add(2, 'hour').unix(),
      },
    ]),
    getInfoFromTeamId: jest.fn().mockReturnValue({
      teamId: '2748779069857',
      name: 'Magicians Plus',
      managerName: 'asdas',
    }),
    getInfoFromPlayerId: jest.fn().mockReturnValue({
      teamId: '2748779069626',
    }),
    getAuction: jest.fn().mockReturnValue({
      playerId: '2748779077296',
      price: '50',
      playerByPlayerId: {
        name: 'Simón Njoroge',
      },
    }),
  };
});

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
}));

afterEach(() => {
  jest.clearAllMocks();
});

test('processPendingPayingBids works correctly', async () => {
  await processPendingPayingBids();

  expect(HorizonService.getPayingBids).toHaveBeenCalledTimes(1);
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(1);
  expect(HorizonService.getInfoFromTeamId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(1);
});
