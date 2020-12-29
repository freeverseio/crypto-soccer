const repositories = require('../repositories');
const HorizonService = require('../services/HorizonService.js');
const BidValidation = require('./BidValidation.js');
const Web3 = require('web3');
const dayjs = require('dayjs');
const Validation = require('./Validation');

jest.mock('../services/HorizonService.js', () => ({
  getBidsPayedByOwner: jest.fn().mockReturnValue([{ extraPrice: 1, auctionByAuctionId: { price: 5 } }]),
  getAuction: jest.fn().mockReturnValue({ price: 10 }),
  getTeamOwner: jest.fn().mockReturnValue('0x7AAB315885FB74a292781e78c33E130be0e326c4'),
  getTeamsByOwner: jest.fn().mockReturnValue([{ teamId: '123123123' }]),
  hasAuctionPass: jest.fn().mockReturnValue(false),
  hasSpentInWorldPlayers: jest.fn().mockReturnValue(false),
  getUnpaymentsByOwner: jest.fn().mockReturnValue([]),
}));

jest.mock('../repositories', () => ({
  selectOwnerMaxBidAllowed: jest.fn().mockReturnValue(undefined),
}));

afterEach(() => {
  jest.clearAllMocks();
});

test('isAllowedByUnpayments is true because there is no unpayment record', async () => {
  const web3 = new Web3('');
  HorizonService.getUnpaymentsByOwner.mockReturnValue({});
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const result = await bidValidation.validation.isAllowedByUnpayments({ owner: 'chico' });
  expect(result).toBe(true);

  HorizonService.getUnpaymentsByOwner.mockReturnValue({});
  const resultWhenIs0 = await bidValidation.validation.isAllowedByUnpayments({ owner: 'chico' });

  expect(resultWhenIs0).toBe(true);
});

test('isAllowedByUnpayments is false because there is 3 unpayments', async () => {
  const web3 = new Web3('');
  HorizonService.getUnpaymentsByOwner.mockReturnValue([
    { owner: 'chico', TimeOfUnpayment: '' },
    { owner: 'chico', TimeOfUnpayment: '' },
    { owner: 'chico', TimeOfUnpayment: '' },
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
  const result = await bidValidation.validation.isAllowedByUnpayments({ owner: 'chico' });
  expect(result).toBe(false);
});

test('isAllowedByUnpayments is false because there is 1 unpayments but from less than 7 days ago', async () => {
  const web3 = new Web3('');
  HorizonService.getUnpaymentsByOwner.mockReturnValue([
    {
      owner: 'chico',
      TimeOfUnpayment: dayjs().subtract(6, 'day'),
    },
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
  const result = await bidValidation.validation.isAllowedByUnpayments({ owner: 'chico' });
  expect(result).toBe(false);
});

test('isAllowedToBidByMaxBidAllowedByOwner is true because there is no max bid allowed set', async () => {
  const web3 = new Web3('');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue(undefined);
  const validation = new Validation();
  const result = await validation.isAllowedToBidByMaxBidAllowedByOwner({ owner: 'chico', totalPrice: 50 });
  expect(result).toBe(true);
});

test('isAllowedToBidByMaxBidAllowedByOwner is true because there is a row but is set to null in db', async () => {
  const web3 = new Web3('');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: null });
  const validation = new Validation();
  const result = await validation.isAllowedToBidByMaxBidAllowedByOwner({ owner: 'chico', totalPrice: 50 });
  expect(result).toBe(true);

  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: undefined });
  expect(result).toBe(true);

  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: 'asda' });
  expect(result).toBe(true);
});

test('isAllowedToBidByMaxBidAllowedByOwner is true because there is a row with an integer value and the price is lesser than max bid', async () => {
  const web3 = new Web3('');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: 1000 });
  const validation = new Validation();
  const result = await validation.isAllowedToBidByMaxBidAllowedByOwner({ owner: 'chico', totalPrice: 50 });
  expect(result).toBe(true);
});

test('isAllowedToBidByMaxBidAllowedByOwner is false because there is a row with an integer value and the price is greater than max bid', async () => {
  const web3 = new Web3('');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: 1000 });
  const validation = new Validation();
  const result = await validation.isAllowedToBidByMaxBidAllowedByOwner({ owner: 'chico', totalPrice: 5000 });
  expect(result).toBe(false);
});

test('true because has spent in worldplayers', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x65879739f8523163300586992eC5c2d0a367ecE7');
  HorizonService.hasSpentInWorldPlayers.mockReturnValue(true);
  const validation = new Validation();
  const result = await validation.isAllowedToBidBySpentInWorldplayers({
    owner: '0x65879739f8523163300586992eC5c2d0a367ecE7',
  });

  expect(result).toBe(true);
});

test('false because has not spent in worldplayers', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x65879739f8523163300586992eC5c2d0a367ecE7');
  HorizonService.hasSpentInWorldPlayers.mockReturnValue(false);
  const validation = new Validation();
  const result = await validation.isAllowedToBidBySpentInWorldplayers({
    owner: '0x65879739f8523163300586992eC5c2d0a367ecE7',
  });

  expect(result).toBe(false);
});
