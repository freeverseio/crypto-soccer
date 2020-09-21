const repositories = require('../repositories');
const HorizonService = require('../services/HorizonService.js');
const BidValidation = require('./BidValidation.js');
const Web3 = require('web3');

jest.mock('../services/HorizonService.js', () => ({
  getBidsPayedByOwner: jest
    .fn()
    .mockReturnValue([{ extraPrice: 1, auctionByAuctionId: { price: 5 } }]),
  getAuction: jest.fn().mockReturnValue({ price: 10 }),
  getTeamOwner: jest
    .fn()
    .mockReturnValue('0x7AAB315885FB74a292781e78c33E130be0e326c4'),
}));

afterEach(() => {
  jest.clearAllMocks();
});

test('Bid Prefixed Hash', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  bidValidation = new BidValidation({
    teamId: '274877906945',
    rnd: 1243523,
    auctionId:
      '55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a',
    extraPrice: 332,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });

  const hash = bidValidation.prefixedHash();

  expect(hash).toBe(
    '0xc0ad1683b9afe071d698763b7143e7cff7bcc661c7074497d870964dd58d9976'
  );
});

test('not allowed to bid because computed maximum bid is lower than total price', async () => {
  const web3 = new Web3('');

  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const result = await bidValidation.isAllowedToBid();

  expect(result).toBe(false);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledWith({
    owner: '0x7AAB315885FB74a292781e78c33E130be0e326c4',
  });
});

test('allowed to bid because computed is greater', async () => {
  const web3 = new Web3('');
  HorizonService.getBidsPayedByOwner.mockReturnValue([
    { extraPrice: 1, auctionByAuctionId: { price: 50 } },
  ]);

  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const result = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledWith({
    owner: '0x7AAB315885FB74a292781e78c33E130be0e326c4',
  });
});

test('not allowed to bid because signer is not owner', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x0000');
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const result = await bidValidation.isAllowedToBid();

  expect(result).toBe(false);
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(0);
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledTimes(0);
});
