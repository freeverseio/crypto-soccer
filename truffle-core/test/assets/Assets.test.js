const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Assets = artifacts.require('Assets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('Assets', (accounts) => {
    let assets = null;
    let playerStateLib = null;
    let PLAYERS_PER_TEAM_MAX = null;
    let PLAYERS_PER_TEAM_INIT = null;
    let LEAGUES_PER_DIV = null;
    let TEAMS_PER_LEAGUE = null;
    let FREE_PLAYER_ID = null;
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const N_SKILLS = 5;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = await assets.PLAYERS_PER_TEAM_INIT().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await assets.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        LEAGUES_PER_DIV = await assets.LEAGUES_PER_DIV().should.be.fulfilled;
        TEAMS_PER_LEAGUE = await assets.TEAMS_PER_LEAGUE().should.be.fulfilled;
        FREE_PLAYER_ID = await assets.FREE_PLAYER_ID().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = PLAYERS_PER_TEAM_INIT.toNumber();
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        LEAGUES_PER_DIV = LEAGUES_PER_DIV.toNumber();
        TEAMS_PER_LEAGUE = TEAMS_PER_LEAGUE.toNumber();
        });

    // it('check initial and max number of players per team', async () =>  {
    //     PLAYERS_PER_TEAM_INIT.should.be.equal(18);
    //     PLAYERS_PER_TEAM_MAX.should.be.equal(25);
    //     LEAGUES_PER_DIV.should.be.equal(16);
    //     TEAMS_PER_LEAGUE.should.be.equal(8);
    // });

    // it('check initial setup of timeZones', async () =>  {
    //     nCountries = await assets.getNCountriesInTZ(0).should.be.rejected;
    //     nCountries = await assets.getNCountriesInTZ(25).should.be.rejected;
    //     for (tz = 1; tz<25; tz++) {
    //         nCountries = await assets.getNCountriesInTZ(tz).should.be.fulfilled;
    //         nCountries.toNumber().should.be.equal(1);
    //         nDivs = await assets.getNDivisionsInCountry(tz, countryIdxInTZ = 0).should.be.fulfilled;
    //         nDivs.toNumber().should.be.equal(1);
    //         nLeagues = await assets.getNLeaguesInCountry(tz, countryIdxInTZ).should.be.fulfilled;
    //         nLeagues.toNumber().should.be.equal(LEAGUES_PER_DIV);
    //         nTeams = await assets.getNTeamsInCountry(tz, countryIdxInTZ).should.be.fulfilled;
    //         nTeams.toNumber().should.be.equal(LEAGUES_PER_DIV * TEAMS_PER_LEAGUE);
    //     }
    // });

    // it('check teamExists for existing teams', async () =>  {
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry = nTeams - 1;
    //     for (tz = 1; tz<25; tz++) {
    //         teamExists = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
    //         teamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
    //         teamExists2 = await assets.teamExists(teamId).should.be.fulfilled;
    //         teamExists.should.be.equal(true);            
    //         teamExists2.should.be.equal(true); 
    //     }
    // });
    
    // it('check teamExists for not-created teams', async () =>  {
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry = nTeams;
    //     for (tz = 1; tz<25; tz++) {
    //         teamExists = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
    //         teamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
    //         teamExists2 = await assets.teamExists(teamId).should.be.fulfilled;
    //         teamExists.should.be.equal(false);            
    //         teamExists2.should.be.equal(false); 
    //     }
    // });
    
    // it('check teamExists for non-existing countries', async () =>  {
    //     countryIdxInTZ = 1;
    //     teamIdxInCountry = nTeams;
    //     for (tz = 1; tz<25; tz++) {
    //         teamExists = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.rejected;
    //         teamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
    //         teamExists2 = await assets.teamExists(teamId).should.be.rejected;
    //     }
    // });

    // it('check playerExists and isVirtual', async () =>  {
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry = nTeams;
    //     playerIdxInCountry = teamIdxInCountry * PLAYERS_PER_TEAM_INIT - 1;
    //     for (tz = 1; tz<25; tz++) {
    //         playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry);
    //         playerExists = await assets.playerExists(playerId).should.be.fulfilled;
    //         playerExists.should.be.equal(true);            
    //         isVirtual = await assets.isVirtualPlayer(playerId).should.be.fulfilled;
    //         isVirtual.should.be.equal(true);            
    //         playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry+1);
    //         playerExists = await assets.playerExists(playerId).should.be.fulfilled;
    //         playerExists.should.be.equal(false);            
    //         isVirtual = await assets.isVirtualPlayer(playerId).should.be.rejected;
    //     }
    // });

    // it('isBot teams', async () =>  {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry = 0;
    //     isBot = await assets.isBotTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
    //     isBot.should.be.equal(true);            
    // });
    
    // it('transfer of bot teams', async () =>  {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry1 = 0;
    //     teamIdxInCountry2 = 1;
    //     teamId1 = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry1);
    //     teamId2 = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry2);
    //     await assets.transferBotInCountryToAddr(tz, countryIdxInTZ, teamIdxInCountry1, ALICE).should.be.fulfilled;
    //     await assets.transferBotToAddr(teamId2, BOB).should.be.fulfilled;
    //     isBot = await assets.isBotTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry1).should.be.fulfilled;
    //     isBot.should.be.equal(false);
    //     isBot = await assets.isBotTeam(teamId2).should.be.fulfilled;
    //     isBot.should.be.equal(false);
    //     owner = await assets.getOwnerTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry1).should.be.fulfilled;
    //     owner.should.be.equal(ALICE);
    //     owner = await assets.getOwnerTeam(teamId2).should.be.fulfilled;
    //     owner.should.be.equal(BOB);
    // });
    
    
    // it('get team player ids', async () => {
    //     // for the first team we should find playerIdx = [0, 1,...,17, FREE, FREE, ...]
    //     teamId = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0);
    //     let ids = await assets.getPlayerIdsInTeam(teamId).should.be.fulfilled;
    //     ids.length.should.be.equal(PLAYERS_PER_TEAM_MAX);
    //     for (shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_MAX; shirtNum++) {
    //         if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
    //             ids[shirtNum].should.be.bignumber.equal(FREE_PLAYER_ID);
    //             continue;
    //         } else {
    //             decoded = await playerStateLib.decodeTZCountryAndVal(ids[shirtNum]).should.be.fulfilled;
    //             const {0: timeZone, 1: country, 2: playerIdxInCountry} = decoded;
    //             playerIdxInCountry.toNumber().should.be.equal(shirtNum);
    //         }
    //     }
    //     // for the first team we should find playerIdx = [18, 19,..., FREE, FREE, ...]
    //     teamId = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 1);
    //     ids = await assets.getPlayerIdsInTeam(teamId).should.be.fulfilled;
    //     ids.length.should.be.equal(PLAYERS_PER_TEAM_MAX);
    //     for (shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_MAX; shirtNum++) {
    //         if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
    //             ids[shirtNum].should.be.bignumber.equal(FREE_PLAYER_ID);
    //             continue;
    //         } else {
    //             decoded = await playerStateLib.decodeTZCountryAndVal(ids[shirtNum]).should.be.fulfilled;
    //             const {0: timeZone, 1: country, 2: playerIdxInCountry} = decoded;
    //             playerIdxInCountry.toNumber().should.be.equal(shirtNum + PLAYERS_PER_TEAM_INIT);
    //         }
    //     }
    // });

    // it('gameDeployMonth', async () => {
    //     const gameDeployMonth =  await assets.gameDeployMonth().should.be.fulfilled;
    //     currentBlockNum = await web3.eth.getBlockNumber()
    //     currentBlock = await web3.eth.getBlock(currentBlockNum)
    //     currentMonth = Math.floor(currentBlock.timestamp * 12 / (3600 * 24 * 365));
    //     gameDeployMonth.toNumber().should.be.equal(currentMonth);
    // });

    // it('get skills of player on creation', async () => {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     playerIdxInCountry = 1;
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     encodedSkills = await assets.getPlayerSkillsAtBirth(playerId).should.be.fulfilled;
    //     skills = await playerStateLib.getSkillsVec(encodedSkills).should.be.fulfilled; 
    //     expected = [43, 57, 58, 61, 31];
    //     for (sk = 0; sk < N_SKILLS; sk++) {
    //         skills[sk].toNumber().should.be.equal(expected[sk])
    //     }
    //     newId =  await playerStateLib.getPlayerIdFromSkills(encodedSkills).should.be.fulfilled; 
    //     newId.should.be.bignumber.equal(playerId);
    //     monthOfBirth =  await playerStateLib.getMonthOfBirth(encodedSkills).should.be.fulfilled; 
    //     monthOfBirth.toNumber().should.be.equal(248);
    //     ageInMonths = await assets.getPlayerAgeInMonths(playerId).should.be.fulfilled;
    //     ageInMonths.toNumber().should.be.equal(348); // 29 years
    // });

    // it('get state of player on creation', async () => {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     // test for players on the first team
    //     playerIdxInCountry = 1;
    //     teamIdxInCountry = Math.floor(playerIdxInCountry / PLAYERS_PER_TEAM_INIT);
    //     teamIdxInCountry.should.be.equal(0);
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     state = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     newId =  await playerStateLib.getPlayerIdFromState(state).should.be.fulfilled; 
    //     newId.should.be.bignumber.equal(playerId);
    //     expectedTeamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
    //     teamId =  await playerStateLib.getCurrentTeamId(state).should.be.fulfilled; 
    //     teamId.should.be.bignumber.equal(expectedTeamId);
    //     shirtNum =  await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled; 
    //     shirtNum.toNumber().should.be.equal(1);
    //     // test for players on the second team
    //     playerIdxInCountry = 18;
    //     teamIdxInCountry = Math.floor(playerIdxInCountry / PLAYERS_PER_TEAM_INIT);
    //     teamIdxInCountry.should.be.equal(1);
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     state = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     newId =  await playerStateLib.getPlayerIdFromState(state).should.be.fulfilled; 
    //     newId.should.be.bignumber.equal(playerId);
    //     expectedTeamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
    //     teamId =  await playerStateLib.getCurrentTeamId(state).should.be.fulfilled; 
    //     teamId.should.be.bignumber.equal(expectedTeamId);
    //     shirtNum =  await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled; 
    //     shirtNum.toNumber().should.be.equal(0);
    // });
    
    // it('get player state of unexistent player', async () => {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT - 1; // last player that exists
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     state = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT; // player not existing
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     state = await assets.getPlayerState(playerId).should.be.rejected;
    //     tz = 0; // dummy timeZone
    //     countryIdxInTZ = 0;
    //     playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT - 1;
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     state = await assets.getPlayerState(playerId).should.be.rejected;
    //     tz = 1; 
    //     countryIdxInTZ = 1; // country not existing
    //     playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT - 1;
    //     playerId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
    //     state = await assets.getPlayerState(playerId).should.be.rejected;
    // });


    // it('isFreeShirt', async () => {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry = 0; 
    //     teamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
    //     // cannot query about a Bot Team
    //     isFree = await assets.isFreeShirt(teamId,shirtNum = 3).should.be.rejected
    //     // so transfer and query again
    //     await assets.transferBotToAddr(teamId, ALICE).should.be.fulfilled;
    //     isBot = await assets.isBotTeam(teamId).should.be.fulfilled;
    //     isBot.should.be.equal(false);
    //     isFree = await assets.isFreeShirt(teamId, shirtNum = 3).should.be.fulfilled
    //     isFree.should.be.equal(false)
    //     isFree = await assets.isFreeShirt(teamId, shirtNum = 18).should.be.fulfilled
    //     isFree.should.be.equal(true)
    // });

    // it('getFreeShirt', async () => {
    //     tz = 1;
    //     countryIdxInTZ = 0;
    //     teamIdxInCountry = 0; 
    //     teamId = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
    //     // cannot query about a Bot Team
    //     shirtNum = await assets.getFreeShirt(teamId).should.be.rejected
    //     // so transfer and query again
    //     await assets.transferBotToAddr(teamId, ALICE).should.be.fulfilled;
    //     isBot = await assets.isBotTeam(teamId).should.be.fulfilled;
    //     isBot.should.be.equal(false);
    //     shirtNum = await assets.getFreeShirt(teamId).should.be.fulfilled
    //     shirtNum.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 1);
    // });

    
    // it('transferPlayer', async () => {
    //     playerId    = await playerStateLib.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
    //     teamId1     = await playerStateLib.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
    //     teamId2     = await playerStateLib.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 2).should.be.fulfilled; 

    //     // state before selling:
    //     state = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     obtainedTeamId = await playerStateLib.getCurrentTeamId(state).should.be.fulfilled;
    //     obtainedTeamId.should.be.bignumber.equal(teamId1);
    //     shirt = await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled;
    //     shirt.toNumber().should.be.equal(playerIdxInCountry1);        

    //     await assets.transferBotToAddr(teamId1, ALICE).should.be.fulfilled;
    //     await assets.transferBotToAddr(teamId2, BOB).should.be.fulfilled;
    //     await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;

    //     // state of player after selling:
    //     state = await assets.getPlayerState(playerId).should.be.fulfilled;
    //     obtainedTeamId = await playerStateLib.getCurrentTeamId(state).should.be.fulfilled;
    //     obtainedTeamId.should.be.bignumber.equal(teamId2);
    //     shirt = await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled;
    //     shirt.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 1);        

    //     // states of teams after selling
    //     isFree = await assets.isFreeShirt(teamId1, shirtNum = playerIdxInCountry1).should.be.fulfilled
    //     isFree.should.be.equal(true);
    //     isFree = await assets.isFreeShirt(teamId2, shirtNum = PLAYERS_PER_TEAM_MAX - 1).should.be.fulfilled
    //     isFree.should.be.equal(false);
    //     shirtNum = await assets.getFreeShirt(teamId2).should.be.fulfilled
    //     shirtNum.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 2);
    // });

    // it('transferPlayer different team works', async () => {
    //     playerId    = await playerStateLib.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
    //     teamId1     = await playerStateLib.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
    //     teamId2     = await playerStateLib.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 2).should.be.fulfilled; 
    //     await assets.transferBotToAddr(teamId1, ALICE).should.be.fulfilled;
    //     await assets.transferBotToAddr(teamId2, ALICE).should.be.fulfilled;
    //     await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;
    // });

    // it('transferPlayer same team fails', async () => {
    //     playerId    = await playerStateLib.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
    //     teamId1     = await playerStateLib.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
    //     await assets.transferBotToAddr(teamId1, ALICE).should.be.fulfilled;
    //     await assets.transferPlayer(playerId, teamId1).should.be.rejected;
    // });

    // it('transferPlayer fails when at least one team involved is a bot', async () => {
    //     playerId1   = await playerStateLib.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
    //     playerId2   = await playerStateLib.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ1 = 0, playerIdxInCountry1 = 8).should.be.fulfilled; 
    //     teamId1     = await playerStateLib.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
    //     teamId2     = await playerStateLib.encodeTZCountryAndVal(tz2, countryIdxInTZ2, teamIdxInCountry = 0).should.be.fulfilled; 
    //     // both teams are bots: fails
    //     await assets.transferPlayer(playerId1, teamId2).should.be.rejected;
    //     // only buyer team is bot: fails
    //     await assets.transferBotToAddr(teamId1, ALICE).should.be.fulfilled;
    //     await assets.transferPlayer(playerId1, teamId2).should.be.rejected;
    //     // only seller team is bot: fails
    //     await assets.transferPlayer(playerId2, teamId1).should.be.rejected;
    //     // // both are owned: works
    //     await assets.transferBotToAddr(teamId2, ALICE).should.be.fulfilled;
    //     await assets.transferPlayer(playerId1, teamId2).should.be.fulfilled;
    //     await assets.transferPlayer(playerId2, teamId1).should.be.fulfilled;
    // });
    
    // it('transferPlayer to already full team', async () => {
    //     teamId     = await playerStateLib.encodeTZCountryAndVal(tz2, countryIdxInTZ2, teamIdxInCountry = 0).should.be.fulfilled; 
    //     for (playerIdxInCountry = 0; playerId < PLAYERS_PER_TEAM_MAX-PLAYERS_PER_TEAM_INIT; playerId++) {
    //         playerId   = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry).should.be.fulfilled; 
    //         await assets.transferPlayer(playerId, teamId).should.be.fulfilled;
    //     }
    //     playerId   = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry+1).should.be.fulfilled; 
    //     await assets.transferPlayer(playerId, teamId).should.be.rejected;
    // });

    // it('team exists', async () => {
    //     teamId     = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled; 
    //     result = await assets.teamExists(teamId).should.be.fulfilled;
    //     result.should.be.equal(true);
    //     teamId     = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = TEAMS_PER_LEAGUE * LEAGUES_PER_DIV - 1).should.be.fulfilled; 
    //     result = await assets.teamExists(teamId).should.be.fulfilled;
    //     result.should.be.equal(true);
    //     teamId     = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = TEAMS_PER_LEAGUE * LEAGUES_PER_DIV).should.be.fulfilled; 
    //     result = await assets.teamExists(teamId).should.be.fulfilled;
    //     result.should.be.equal(false);
    //     teamId     = await playerStateLib.encodeTZCountryAndVal(tz = 0, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled; 
    //     result = await assets.teamExists(teamId).should.be.rejected;
    // });

    // it('initial number of countries', async () => {
    //     const count = await assets.countCountries(tz = 1).should.be.fulfilled;
    //     count.toNumber().should.be.equal(1);
    // });

    // it('initial number of teams', async () => {
    //     const count = await assets.countTeams(tz = 1, countryIdxInTZ = 0).should.be.fulfilled;
    //     count.toNumber().should.be.equal(TEAMS_PER_LEAGUE * LEAGUES_PER_DIV);
    // });

    // it('existence of null player', async () => {
    //     const exists = await assets.playerExists(playerId = 0).should.be.fulfilled;
    //     exists.should.be.equal(false);
    // });

    // it('isVirtual after sale', async () => {
    //     playerId    = await playerStateLib.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
    //     teamId1     = await playerStateLib.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
    //     teamId2     = await playerStateLib.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 2).should.be.fulfilled; 
    //     await assets.transferBotToAddr(teamId1, ALICE).should.be.fulfilled;
    //     await assets.transferBotToAddr(teamId2, ALICE).should.be.fulfilled;
    //     isVirtual = await assets.isVirtualPlayer(playerId).should.be.fulfilled;
    //     isVirtual.should.be.equal(true)
    //     await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;
    //     isVirtual = await assets.isVirtualPlayer(playerId).should.be.fulfilled;
    //     isVirtual.should.be.equal(false)
    // });

    
    // it('computed skills with rnd = 0 is 50 each', async () => {
    //     let skills = await assets.computeSkills(0).should.be.fulfilled;
    //     skills.forEach(skill => (skill.toNumber().should.be.equal(50)));
    // });

    // it('sum of computed skills is 250', async () => {
    //     for (let i = 0; i < 10; i++) {
    //         const seed = 32123;
    //         const skills = await assets.computeSkills(seed).should.be.fulfilled;
    //         const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
    //         sum.should.be.equal(250);
    //     }
    // });

    it('get shirtNum in team for all players in a country', async () => {
        tz = 0;
        countryIdxInTZ = 0;
        playersInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT
        for (let playerIdxInCountry = 0; playerId < playersInCountry ; playerId++){
            playerId    = await playerStateLib.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
            const playerState = await assets.getPlayerState(playerId).should.be.fulfilled;
            const shirtNum = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
            shirtNum.toNumber().should.be.equal(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        }
    })
    
    
    return;

    // it('get existing virtual player skills', async () => {
    //     const numSkills = await assets.NUM_SKILLS().should.be.fulfilled;
    //     await assets.createTeam("Barca",ALICE).should.be.fulfilled;
    //     const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
    //     const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
    //     skills.length.should.be.equal(numSkills.toNumber());
    //     skills[0].should.be.bignumber.equal('78');
    //     skills[1].should.be.bignumber.equal('65');
    //     skills[2].should.be.bignumber.equal('35');
    //     skills[3].should.be.bignumber.equal('35');
    //     skills[4].should.be.bignumber.equal('37');
    //     const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
    //     sum.should.be.equal(250);
    // });


    // it('get existing non virtual player skills', async () => {
    //     await assets.createTeam("Barca",ALICE).should.be.fulfilled;
    //     const state = await playerStateLib.playerStateCreate(
    //         defence = 1,
    //         speed = 2,
    //         pass = 3,
    //         shoot = 4,
    //         endurance = 5,
    //         monthOfBirthInUnixTime = 6,
    //         playerId = 10,
    //         currentTeamId = 1,
    //         currentShirtNum = 3,
    //         prevLeagueId = 3,
    //         prevTeamPosInLeague = 3,
    //         prevShirtNumInLeague = 3,
    //         lastSaleBlock = 3
    //     ).should.be.fulfilled;
    //     await assets.setPlayerState(state).should.be.fulfilled;
    //     const playerState = await assets.getPlayerState(playerId = 10).should.be.fulfilled;
    //     const skills = await playerStateLib.getSkillsVec(playerState).should.be.fulfilled;
    //     skills[0].should.be.bignumber.equal('1');
    //     skills[1].should.be.bignumber.equal('2');
    //     skills[2].should.be.bignumber.equal('3');
    //     skills[3].should.be.bignumber.equal('4');
    //     skills[4].should.be.bignumber.equal('5');
    // });

    it('compute player birth', async () => {
        const birth = await assets.computeBirth(0, 1557495456).should.be.fulfilled;
        birth.should.be.bignumber.equal('406');
    });

    // it('get non virtual player team', async () => {
    //     await assets.createTeam("Barca",ALICE).should.be.fulfilled;
    //     await assets.createTeam("Madrid",ALICE).should.be.fulfilled;
    //     let playerState = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
    //     const teamBefore = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
    //     const state = await playerStateLib.playerStateCreate(
    //         defence = 3,
    //         speed = 3,
    //         pass = 3,
    //         shoot = 3,
    //         endurance = 3,
    //         monthOfBirthInUnixTime = 3,
    //         playerId = 1,
    //         currentTeamId = 2,
    //         currentShirtNum = 3,
    //         prevLeagueId = 3,
    //         prevTeamPosInLeague = 3,
    //         prevShirtNumInLeague = 3,
    //         lastSaleBlock = 3
    //     ).should.be.fulfilled;
    //     await assets.setPlayerState(state).should.be.fulfilled;
    //     playerState = await assets.getPlayerState(playerId = 1).should.be.fulfilled;
    //     const teamAfter = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
    //     teamAfter.should.be.bignumber.not.equal(teamBefore);
    //     teamAfter.should.be.bignumber.equal('2');
    // });

    // it('create team', async () => {
    //     const receipt = await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
    //     const count = await assets.countTeams().should.be.fulfilled;
    //     count.toNumber().should.be.equal(1);
    //     const teamId = receipt.logs[0].args.id.toNumber();
    //     teamId.should.be.equal(1);
    //     teamName = await assets.getTeamName(teamId).should.be.fulfilled;
    //     teamName.should.be.equal("Barca",ALICE);
    // });

    // it('get playersId from teamId and pos in team', async () => {
    //     await assets.generateVirtualPlayerId(teamId = 1, posInTeam=0).should.be.rejected;
    //     await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
    //     await assets.generateVirtualPlayerId(teamId = 1, posInTeam=PLAYERS_PER_TEAM_MAX).should.be.rejected;
    //     let playerId = await assets.generateVirtualPlayerId(teamId = 1, posInTeam=0).should.be.fulfilled;
    //     playerId.toNumber().should.be.equal(1);
    //     playerId = await assets.generateVirtualPlayerId(teamId = 1, posInTeam=PLAYERS_PER_TEAM_MAX-1).should.be.fulfilled;
    //     playerId.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX);

    // });

    // it('sign team to league', async () => {
    //     await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
    //     await assets.createTeam(name = "Barca",ALICE).should.be.fulfilled;
    //     await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.fulfilled;
    //     const currentHistory = await assets.getTeamCurrentHistory(1).should.be.fulfilled;
    //     currentHistory.currentLeagueId.should.be.bignumber.equal('1');
    //     currentHistory.posInCurrentLeague.should.be.bignumber.equal('3');
    //     currentHistory.prevLeagueId.should.be.bignumber.equal('0');
    //     currentHistory.posInPrevLeague.should.be.bignumber.equal('0');
    // });

    // it('sign team to league twice should fail', async () => {
    //     await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
    //     await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.rejected;
    // });
    
    // it('transfer team', async () => {
    //     await assets.createTeam(name = "Barca", ALICE).should.be.fulfilled;
    //     const currentOwner = await assets.getTeamOwner(name).should.be.fulfilled;
    //     currentOwner.should.be.equal(ALICE);
    //     let tx = await assets.transferTeam(teamId = 1, BOB).should.be.fulfilled;
    //     const newOwner = await assets.getTeamOwner(name).should.be.fulfilled;
    //     newOwner.should.be.equal(BOB);
    //     truffleAssert.eventEmitted(tx, "TeamTransfer", (ev) => {
    //         return ev.teamId == 1 && ev.to == BOB;
    //     });
    // });

    // it ('transfer invalid team 0', async () => {
    //     await assets.transferTeam(teamId = 0, BOB).should.be.rejected;
    // });
        
    // it('transfer non-exisiting team', async () => {
    //     await assets.transferTeam(teamId = 1, BOB).should.be.rejected;
    // });

    // it('transfer team accross same owner', async () => {
    //     await assets.createTeam(name = "Barca", ALICE).should.be.fulfilled;
    //     await assets.transferTeam(teamId = 1, ALICE).should.be.rejected;
    // });
});