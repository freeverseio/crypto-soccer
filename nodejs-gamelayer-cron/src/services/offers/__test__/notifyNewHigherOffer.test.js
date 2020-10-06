const GamelayerService = require('../../GamelayerService');
const notifyNewHigherOffer = require('../notifyNewHigherOffer');

jest.mock('../../GamelayerService', () => ({
  setMessage: jest.fn(),
}));

afterEach(() => {
  jest.clearAllMocks();
});

test('notifyNewHigherOffer works correctly', async () => {
  await notifyNewHigherOffer({
    destinatary: '1234324',
    auctionId: 'ab652d3123',
  });

  expect(GamelayerService.setMessage).toHaveBeenCalledTimes(1);
});
