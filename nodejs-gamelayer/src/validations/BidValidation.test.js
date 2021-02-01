const repositories = require('../repositories');
const HorizonService = require('../services/HorizonService.js');
const BidValidation = require('./BidValidation.js');
const Web3 = require('web3');
const dayjs = require('dayjs');

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

test('Bid Prefixed Hash', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  bidValidation = new BidValidation({
    teamId: '274877906945',
    rnd: 1243523,
    auctionId: '55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a',
    extraPrice: 332,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });

  const hash = bidValidation.prefixedHash();

  expect(hash).toBe('0xe791281515bce955edbc5cef6af64fcc018a5a7b0384f7cc5357b9c40476983a');
});

test('Bid expected owner', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  bidValidation = new BidValidation({
    teamId: '274877906945',
    rnd: 1243523,
    auctionId: '0000000000000000000000000000000000000000000000000000000000032123',
    extraPrice: 332,
    signature:
      'fae9e592282290bc3dc0b650f539e88e0cff58df8459ba0ebf4311c6c96848dd589856be8288fec610b2e65ea36e49b48e1ffdc80eab994346a696110cdee6ae1c',
    web3,
  });

  const addressRecovered = await bidValidation.signerAddress();

  expect(addressRecovered).toBe('0x1760B51E59C7378607B59aA05aB315AeB4c8C39F');
});

test('allowed to bid because bid is lower than minimum default bid', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD');
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to bid because has auction passs', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x65879739f8523163300586992eC5c2d0a367ecE7');
  HorizonService.hasAuctionPass.mockReturnValue(true);
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 2000,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.hasAuctionPass).toHaveBeenCalledWith({ owner: '0x65879739f8523163300586992eC5c2d0a367ecE7' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalledWith();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to bid because has spent in worldplayers', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x65879739f8523163300586992eC5c2d0a367ecE7');
  HorizonService.hasAuctionPass.mockReturnValue(false);
  HorizonService.hasSpentInWorldPlayers.mockReturnValue(true);
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 2000,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.hasAuctionPass).toHaveBeenCalledWith({ owner: '0x65879739f8523163300586992eC5c2d0a367ecE7' });
  expect(HorizonService.hasSpentInWorldPlayers).toHaveBeenCalledWith({ teamId: '123123123' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalledWith();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to bid because has spent in bids', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x65879739f8523163300586992eC5c2d0a367ecE7');
  HorizonService.hasAuctionPass.mockReturnValue(false);
  HorizonService.hasSpentInWorldPlayers.mockReturnValue(false);
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 2000,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.hasAuctionPass).toHaveBeenCalledWith({ owner: '0x65879739f8523163300586992eC5c2d0a367ecE7' });
  expect(HorizonService.getTeamsByOwner).toHaveBeenCalledWith({ owner: '0x65879739f8523163300586992eC5c2d0a367ecE7' });
  expect(HorizonService.hasSpentInWorldPlayers).toHaveBeenCalledWith({ teamId: '123123123' });
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledWith({
    owner: '0x65879739f8523163300586992eC5c2d0a367ecE7',
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
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(false);
  expect(HorizonService.getAuction).toHaveBeenCalledTimes(0);
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledTimes(0);
});

test('not allowed to bid because max bid allowed by owner is 0', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: 0 });
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(false);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(repositories.selectOwnerMaxBidAllowed).toHaveBeenCalledWith({
    owner: '0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD',
  });
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to bid because bid is lower than minimum default bid and maxbidowner has row and is set to null', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: null });
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(repositories.selectOwnerMaxBidAllowed).toHaveBeenCalledWith({
    owner: '0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD',
  });
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to bid because bid is lower than minimum default bid and less thahn maxbidowner', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: 21 });
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 10,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const { allowed: result } = await bidValidation.isAllowedToBid();

  expect(result).toBe(true);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(repositories.selectOwnerMaxBidAllowed).toHaveBeenCalledWith({
    owner: '0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD',
  });
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

// 2020-12-15 11:08:25
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
