const HorizonService = require('../../HorizonService');
const getTeamIdFromAuctionSeller = require('../getTeamIdFromAuctionSeller');

jest.mock('../../HorizonService.js', () => ({
  getTeamIdsFromOwner: jest.fn().mockReturnValue([
    {
      teamId: '2748779069857',
    },
    {
      teamId: '2748779069858',
    },
    {
      teamId: '2748779069859',
    },
  ]),
  getPlayerHistoriesLast30BlockNumberTeams: jest.fn().mockReturnValue([
    {
      teamId: '2748779069852',
      blockNumber: '12631330',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12631151',
    },
    {
      teamId: '2748779069858',
      blockNumber: '12624584',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12624396',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12614430',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12614252',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12608255',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12608110',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12598615',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12598436',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12591878',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12591698',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12581495',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12581315',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12574755',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12574575',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12564382',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12564202',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12557589',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12557410',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12547468',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12547247',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12540818',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12540639',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12530734',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12530569',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12526549',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12526500',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12514995',
    },
    {
      teamId: '2748779069852',
      blockNumber: '12514818',
    },
  ]),
}));

afterEach(() => {
  jest.clearAllMocks();
});

test('getTeamIdFromAuctionSeller from history when seller gets paid and the player has already changed from seller team to buyer team and also the owner has more than one team', async () => {
  const auction = {
    seller: '0x000',
    playerId: '234234234',
  };
  const result = await getTeamIdFromAuctionSeller({ auction });

  expect(HorizonService.getTeamIdsFromOwner).toHaveBeenCalledTimes(1);
  expect(
    HorizonService.getPlayerHistoriesLast30BlockNumberTeams
  ).toHaveBeenCalledTimes(1);
  expect(result).toBe('2748779069858');
});

test('getTeamIdFromAuctionSeller from history when seller gets paid and the player has not changed from seller team to buyer team and also the owner has more than one team', async () => {
  const auction = {
    seller: '0x000',
    playerId: '234234234',
  };
  HorizonService.getPlayerHistoriesLast30BlockNumberTeams.mockReset();
  HorizonService.getPlayerHistoriesLast30BlockNumberTeams.mockReturnValue([
    {
      teamId: '2748779069859',
      blockNumber: '12631330',
    },
    {
      teamId: '2748779069859',
      blockNumber: '12631151',
    },
    {
      teamId: '2748779069859',
      blockNumber: '12624584',
    },
    {
      teamId: '2748779069859',
      blockNumber: '12624396',
    },
    {
      teamId: '2748779069859',
      blockNumber: '12614430',
    },
  ]);
  const result = await getTeamIdFromAuctionSeller({ auction });

  expect(HorizonService.getTeamIdsFromOwner).toHaveBeenCalledTimes(1);
  expect(
    HorizonService.getPlayerHistoriesLast30BlockNumberTeams
  ).toHaveBeenCalledTimes(1);
  expect(result).toBe('2748779069859');
});

test('getTeamIdFromAuctionSeller from history when seller gets paid and the player has already changed from seller team to buyer team but the owner has one team', async () => {
  const auction = {
    seller: '0x000',
    playerId: '234234234',
  };
  HorizonService.getTeamIdsFromOwner.mockReset();
  HorizonService.getTeamIdsFromOwner.mockReturnValue([
    {
      teamId: '2748779069857',
    },
  ]);
  const result = await getTeamIdFromAuctionSeller({ auction });

  expect(HorizonService.getTeamIdsFromOwner).toHaveBeenCalledTimes(1);
  expect(
    HorizonService.getPlayerHistoriesLast30BlockNumberTeams
  ).toHaveBeenCalledTimes(0);
  expect(result).toBe('2748779069857');
});
