const processOffersHistories = require('../processOffersHistories');
const HorizonService = require('../../HorizonService');
const processAcceptedOffers = require('../processAcceptedOffers');
const processRejectedOffers = require('../processRejectedOffers');
const processStartedOffers = require('../processStartedOffers');
const selectLastChecked = require('../../../repositories/selectLastChecked');
const updateLastChecked = require('../../../repositories/updateLastChecked');

jest.mock('../../HorizonService.js', () => ({
  getLastOfferHistories: jest.fn().mockReturnValue([
    {
      insertedAt: '2020-09-29T10:13:47.226091+00:00',
      auctionId:
        'fa3944fc76f6d9b8dc9775ba385dbdbef34b19b0f39e2d37fa728d094885b4f4',
      playerId: '2748779076705',
      currencyId: 1,
      price: '50',
      rnd: '375503914',
      validUntil: '1601374670',
      signature:
        '6beb027ad523f3aa2a6629ab6153aa2a6a94f89597d31fbeef4e0170e7f0e0f36ee7eda8b2debcfcf90362ce6f6fa88b2c7a1850f336dc1911797e8e3a6270b61c',
      state: 'ACCEPTED',
      stateExtra: '',
      seller: '0xF7dF8923eE9De53e5ffC40F51F96df72bAcC0BA4',
      buyer: '0xaC347a9Fa330c6c23136F1460086D436ed55a3f8',
      buyerTeamId: '2748779069845',
    },
    {
      insertedAt: '2020-09-29T12:44:34.023977+00:00',
      auctionId:
        'a74a90850037f4e927721f0a55caa1be9880df764a544a6955efdd957733f876',
      playerId: '2748779076846',
      currencyId: 1,
      price: '50',
      rnd: '1418804107',
      validUntil: '1601383774',
      signature:
        '041d52d836b56db813e03f9cfee37f294c82917874f8bbe7c8732db5174c2af634780a55c3385140bce3d1fba1f2fdcdd35525942a2ac1aa8bb9a952eea6551e1b',
      state: 'STARTED',
      stateExtra: '',
      seller: '0x9f46F66b079d469920e4e72a99ef42D8A3447C10',
      buyer: '0xAD13847726E798B9faC50748F0784421A86dfBB4',
      buyerTeamId: '2748779069850',
    },
    {
      insertedAt: '2020-09-29T12:45:01.718129+00:00',
      auctionId:
        '3829046fe1b9f3a74ad3ee8c28edf2c79761737e785dc22d7339e756449b280a',
      playerId: '2748779076846',
      currencyId: 1,
      price: '100',
      rnd: '1418804107',
      validUntil: '1601383774',
      signature:
        'd175435f7dd55ee4edb8388b112df707cc7d9cc12da43486870b679a277dc6023a141ecc59b77ea753b35916e8e57d0c027ba4e84e7110dbe3a0c2ee39f3e8d81b',
      state: 'REJECTED',
      stateExtra: '',
      seller: '0x9f46F66b079d469920e4e72a99ef42D8A3447C10',
      buyer: '0xf1735CAdC166a115CA170A1423DF4DB6dAE53bf8',
      buyerTeamId: '2748779069852',
    },
  ]),
}));

jest.mock('../processAcceptedOffers', () => jest.fn());
jest.mock('../processRejectedOffers', () => jest.fn());
jest.mock('../processStartedOffers', () => jest.fn());
jest.mock('../../../repositories/selectLastChecked', () => jest.fn());
jest.mock('../../../repositories/updateLastChecked', () => jest.fn());

afterEach(() => {
  jest.clearAllMocks();
});

test('processOffersHistories works correctly', async () => {
  await processOffersHistories();

  expect(selectLastChecked).toHaveBeenCalled();
  expect(updateLastChecked).toHaveBeenCalledWith({
    entity: 'offer',
    lastChecked: '2020-09-29T12:45:01.718129+00:00',
  });
  expect(HorizonService.getLastOfferHistories).toHaveBeenCalledTimes(1);
  expect(processStartedOffers).toHaveBeenCalledTimes(1);
  expect(processAcceptedOffers).toHaveBeenCalledTimes(1);
  expect(processRejectedOffers).toHaveBeenCalledTimes(1);
});
