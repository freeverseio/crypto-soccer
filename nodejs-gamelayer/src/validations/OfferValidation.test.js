const repositories = require('../repositories');
const HorizonService = require('../services/HorizonService.js');
const OfferValidation = require('./OfferValidation.js');
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

test('Offer Prefixed Hash', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b',
    web3,
  });

  const hash = offerValidation.prefixedHash();

  expect(hash).toBe('0x48280aaca3224b385bcc4e4b94662cbf17f989f99426943da0e1a10cd2e5a4a0');
});

test('Offer expected owner', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });

  const addressRecovered = await offerValidation.signerAddress();

  expect(addressRecovered).toBe('0x83A909262608c650BD9b0ae06E29D90D0F67aC5e');
});

test('Compute auction id', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '25723578238440869144533393071649442553899076447028039543423578',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });

  const auctionId = offerValidation.computeAuctionId();

  expect(auctionId).toBe('a51ac22be9bcc29236b88373d494c30c5dff2108059aa14c00e4c2a1d4061eba');
});

test('allowed to offer because offer is lower than minimum default offer', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0xCC952556aff48bBFE0D48edA178f29c928E25448');
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 10,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(true);
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('not allowed to offer because offer is grater than minimum default offer', async () => {
  const web3 = new Web3('');
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(false);
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to offer because has auction passs', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x83A909262608c650BD9b0ae06E29D90D0F67aC5e');
  HorizonService.hasAuctionPass.mockReturnValue(true);
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(true);
  expect(HorizonService.hasAuctionPass).toHaveBeenCalledWith({ owner: '0x83A909262608c650BD9b0ae06E29D90D0F67aC5e' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalledWith();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to offer because has spent in worldplayers', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x83A909262608c650BD9b0ae06E29D90D0F67aC5e');
  HorizonService.hasAuctionPass.mockReturnValue(false);
  HorizonService.hasSpentInWorldPlayers.mockReturnValue(true);
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(true);
  expect(HorizonService.hasAuctionPass).toHaveBeenCalledWith({ owner: '0x83A909262608c650BD9b0ae06E29D90D0F67aC5e' });
  expect(HorizonService.hasSpentInWorldPlayers).toHaveBeenCalledWith({ teamId: '123123123' });
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalledWith();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to offer because has spent in offers', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x83A909262608c650BD9b0ae06E29D90D0F67aC5e');
  HorizonService.hasAuctionPass.mockReturnValue(false);
  HorizonService.hasSpentInWorldPlayers.mockReturnValue(false);
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(true);
  expect(HorizonService.hasAuctionPass).toHaveBeenCalledWith({ owner: '0x83A909262608c650BD9b0ae06E29D90D0F67aC5e' });
  expect(HorizonService.getTeamsByOwner).toHaveBeenCalledWith({ owner: '0x83A909262608c650BD9b0ae06E29D90D0F67aC5e' });
  expect(HorizonService.hasSpentInWorldPlayers).toHaveBeenCalledWith({ teamId: '123123123' });
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledWith({
    owner: '0x83A909262608c650BD9b0ae06E29D90D0F67aC5e',
  });
});

test('not allowed to offer because signer is not owner', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x0000');
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(false);
  expect(HorizonService.getBidsPayedByOwner).toHaveBeenCalledTimes(0);
});

test('not allowed to offer because max offer allowed by owner is 0', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0x83A909262608c650BD9b0ae06E29D90D0F67aC5e');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: 0 });
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 41234,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(false);
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(repositories.selectOwnerMaxBidAllowed).toHaveBeenCalledWith({
    owner: '0x83A909262608c650BD9b0ae06E29D90D0F67aC5e',
  });
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to offer because offer is lower than minimum default offer and maxbidowner has row and is set to null', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0xCC952556aff48bBFE0D48edA178f29c928E25448');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_bid_allowed: null });
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 10,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(true);
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(repositories.selectOwnerMaxBidAllowed).toHaveBeenCalledWith({
    owner: '0xCC952556aff48bBFE0D48edA178f29c928E25448',
  });
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});

test('allowed to offer because offer is lower than minimum default offer and less thahn maxbidowner', async () => {
  const web3 = new Web3('');
  HorizonService.getTeamOwner.mockReturnValue('0xCC952556aff48bBFE0D48edA178f29c928E25448');
  repositories.selectOwnerMaxBidAllowed.mockReturnValue({ max_offer_allowed: 21 });
  offerValidation = new OfferValidation({
    currencyId: 1,
    playerId: '10',
    price: 10,
    validUntil: '2000000000',
    buyerTeamId: '20',
    rnd: 42321,
    signature:
      '84beb98c6770b70000e37e56adf4a46b0c208bb207e4a2fc5a510d42f2186a500e1ab0c2586fdebf8fed7843d0f636c438de7a111f926de6adcbcb52da6b63141b---',
    web3,
  });
  const result = await offerValidation.isAllowedToOffer();

  expect(result).toBe(true);
  expect(HorizonService.getTeamsByOwner).not.toHaveBeenCalled();
  expect(repositories.selectOwnerMaxBidAllowed).toHaveBeenCalledWith({
    owner: '0xCC952556aff48bBFE0D48edA178f29c928E25448',
  });
  expect(HorizonService.hasAuctionPass).not.toHaveBeenCalled();
  expect(HorizonService.hasSpentInWorldPlayers).not.toHaveBeenCalled();
  expect(HorizonService.getBidsPayedByOwner).not.toHaveBeenCalled();
});
