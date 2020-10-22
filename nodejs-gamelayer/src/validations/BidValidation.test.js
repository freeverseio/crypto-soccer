const repositories = require('../repositories');
const HorizonService = require('../services/HorizonService.js');
const BidValidation = require('./BidValidation.js');
const Web3 = require('web3');

jest.mock('../services/HorizonService.js', () => ({
  getBidsPayedByOwner: jest.fn().mockReturnValue([{ extraPrice: 1, auctionByAuctionId: { price: 5 } }]),
  getAuction: jest.fn().mockReturnValue({ price: 10 }),
  getTeamOwner: jest.fn().mockReturnValue('0x7AAB315885FB74a292781e78c33E130be0e326c4'),
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

test('not allowed to bid because computed maximum bid is lower than total price', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x65879739f8523163300586992eC5c2d0a367ecE7');
  bidValidation = new BidValidation({
    teamId: '234',
    rnd: 12345,
    auctionId: '555',
    extraPrice: 2000,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });
  const result = await bidValidation.isAllowedToBid();

  expect(result).toBe(false);
  expect(HorizonService.getAuction).toHaveBeenCalledWith({ auctionId: '555' });
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledWith({
    owner: '0x65879739f8523163300586992eC5c2d0a367ecE7',
  });
});

test('allowed to bid because computed is greater', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD');
  HorizonService.getBidsPayedByOwner.mockRestore();
  HorizonService.getBidsPayedByOwner.mockReturnValueOnce([{ extraPrice: 1001, auctionByAuctionId: { price: 50 } }]);

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
    owner: '0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD',
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

test('allowed to bid because bid is less than minimum default bid', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD');
  HorizonService.getBidsPayedByOwner.mockRestore();
  HorizonService.getBidsPayedByOwner.mockReturnValueOnce([]);

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
    owner: '0x6a3F80b7171db8EdD14fd2b1f265BcA7F0d839fD',
  });
});
