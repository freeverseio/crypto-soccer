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
        auctionId:
          'cb127f7252fb10da3599484d3a33a682e505793f9b590a6c4e2ed6bd36e6a9a9',
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
        auctionId:
          'cb127f7252fb10da3599484d3a33a682e505793f9b590a6c4e2ed6bd36e6a9b9',
      },
    ]),

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
  getInfoFromTeamId: jest.fn().mockReturnValue({
    teamId: '2748779069857',
    name: 'Magicians Plus',
    managerName: 'asdas',
  }),
  getMessages: jest.fn().mockReturnValue([
    {
      destinatary: '2748779069616',
      id: '23',
      category: 'auction',
      auctionId:
        'cb127f7252fb10da3599484d3a33a682e505793f9b590a6c4e2ed6bd36e6a9b9',
      title: 'undefined',
      text: 'auction_buyer_wins_auction',
      customImageUrl: '',
      metadata:
        '{"bidderTeamId":"2748779069616", "bidderTeamName":"RED BULL", "amount": "50", "playerId": "2748779076574", "playerName":"undefined", "paymentDeadline":"undefined"}',
      isRead: false,
      createdAt: '2020-10-26T15:37:06+00:00',
    },
    {
      destinatary: '2748779069616',
      id: '25',
      category: 'auction',
      auctionId:
        'cb127f7252fb10da3599484d3a33a682e505793f9b590a6c4e2ed6bd36e6a9b9',
      title: 'undefined',
      text: 'auction_buyer_wins_auction',
      customImageUrl: '',
      metadata:
        '{"bidderTeamId":"2748779069616", "bidderTeamName":"RED BULL", "amount": "50", "playerId": "2748779076574", "playerName":"undefined", "paymentDeadline":"undefined"}',
      isRead: false,
      createdAt: '2020-10-26T15:37:06+00:00',
    },
  ]),
}));

afterEach(() => {
  jest.clearAllMocks();
});

test('processPendingPayingBids works correctly', async () => {
  await processPendingPayingBids();

  expect(HorizonService.getPayingBids).toHaveBeenCalledTimes(1);
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(1);
  expect(GamelayerService.getMessages).toHaveBeenCalledTimes(1);
  expect(GamelayerService.getInfoFromTeamId).toHaveBeenCalledTimes(1);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(1);
});

test('processPendingPayingBids works correctly when there is already a pending payment message', async () => {
  GamelayerService.getMessages.mockReset();
  GamelayerService.getMessages.mockReturnValue([
    {
      destinatary: '2748779069616',
      id: '23',
      category: 'auction',
      auctionId:
        'cb127f7252fb10da3599484d3a33a682e505793f9b590a6c4e2ed6bd36e6a9b9',
      title: 'undefined',
      text: 'auction_buyer_wins_auction',
      customImageUrl: '',
      metadata:
        '{"bidderTeamId":"2748779069616", "bidderTeamName":"RED BULL", "amount": "50", "playerId": "2748779076574", "playerName":"undefined", "paymentDeadline":"undefined"}',
      isRead: false,
      createdAt: '2020-10-26T15:37:06+00:00',
    },
    {
      destinatary: '2748779069616',
      id: '1058',
      category: 'auction',
      auctionId:
        'cb127f7252fb10da3599484d3a33a682e505793f9b590a6c4e2ed6bd36e6a9b9',
      title: '',
      text: 'auction_buyer_pending_payment',
      customImageUrl: '',
      metadata:
        '{"amount":"50", "playerId":"2748779076574", "playerName":"Lucas Menendez", "bidderTeamId":"2748779069616", "bidderTeamName":"RED BULL", "paymentDeadline":"1603880321"}',
      isRead: false,
      createdAt: '2020-10-28T12:52:00+00:00',
    },
  ]);
  await processPendingPayingBids();

  expect(HorizonService.getPayingBids).toHaveBeenCalledTimes(1);
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(0);
  expect(GamelayerService.getMessages).toHaveBeenCalledTimes(1);
  expect(GamelayerService.getInfoFromTeamId).toHaveBeenCalledTimes(0);
  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(0);
});
