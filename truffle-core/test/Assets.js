const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const Assets = artifacts.require('Assets');

contract('Assets', (accounts) => {
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const N_SKILLS = 5;

    beforeEach(async () => {
        assets = await Assets.new().should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = await assets.PLAYERS_PER_TEAM_INIT().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await assets.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        LEAGUES_PER_DIV = await assets.LEAGUES_PER_DIV().should.be.fulfilled;
        TEAMS_PER_LEAGUE = await assets.TEAMS_PER_LEAGUE().should.be.fulfilled;
        FREE_PLAYER_ID = await assets.FREE_PLAYER_ID().should.be.fulfilled;
        NULL_ADDR = await assets.NULL_ADDR().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = PLAYERS_PER_TEAM_INIT.toNumber();
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        LEAGUES_PER_DIV = LEAGUES_PER_DIV.toNumber();
        TEAMS_PER_LEAGUE = TEAMS_PER_LEAGUE.toNumber();
        });

    it('check cannot initialize contract twice', async () =>  {
        await assets.init().should.be.rejected;
    });

    it('check initial and max number of players per team', async () =>  {
        PLAYERS_PER_TEAM_INIT.should.be.equal(18);
        PLAYERS_PER_TEAM_MAX.should.be.equal(25);
        LEAGUES_PER_DIV.should.be.equal(16);
        TEAMS_PER_LEAGUE.should.be.equal(8);
    });

    it('check initial setup of timeZones', async () =>  {
        nCountries = await assets.getNCountriesInTZ(0).should.be.rejected;
        nCountries = await assets.getNCountriesInTZ(25).should.be.rejected;
        for (tz = 1; tz<25; tz++) {
            nCountries = await assets.getNCountriesInTZ(tz).should.be.fulfilled;
            nCountries.toNumber().should.be.equal(1);
            nDivs = await assets.getNDivisionsInCountry(tz, countryIdxInTZ = 0).should.be.fulfilled;
            nDivs.toNumber().should.be.equal(1);
            nLeagues = await assets.getNLeaguesInCountry(tz, countryIdxInTZ).should.be.fulfilled;
            nLeagues.toNumber().should.be.equal(LEAGUES_PER_DIV);
            nTeams = await assets.getNTeamsInCountry(tz, countryIdxInTZ).should.be.fulfilled;
            nTeams.toNumber().should.be.equal(LEAGUES_PER_DIV * TEAMS_PER_LEAGUE);
        }
    });

    it('check teamExists for existing teams', async () =>  {
        countryIdxInTZ = 0;
        teamIdxInCountry = nTeams - 1;
        for (tz = 1; tz<25; tz++) {
            teamExists = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
            teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
            teamExists2 = await assets.teamExists(teamId).should.be.fulfilled;
            teamExists.should.be.equal(true);            
            teamExists2.should.be.equal(true); 
        }
    });
    
    it('check teamExists for not-created teams', async () =>  {
        countryIdxInTZ = 0;
        teamIdxInCountry = nTeams;
        for (tz = 1; tz<25; tz++) {
            teamExists = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
            teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
            teamExists2 = await assets.teamExists(teamId).should.be.fulfilled;
            teamExists.should.be.equal(false);            
            teamExists2.should.be.equal(false); 
        }
    });
    
    it('check teamExists for non-existing countries', async () =>  {
        countryIdxInTZ = 1;
        teamIdxInCountry = nTeams;
        for (tz = 1; tz<25; tz++) {
            teamExists = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.rejected;
            teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
            teamExists2 = await assets.teamExists(teamId).should.be.rejected;
        }
    });

    it('check playerExists and isVirtual', async () =>  {
        countryIdxInTZ = 0;
        teamIdxInCountry = nTeams;
        playerIdxInCountry = teamIdxInCountry * PLAYERS_PER_TEAM_INIT - 1;
        for (tz = 1; tz<25; tz++) {
            playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry);
            playerExists = await assets.playerExists(playerId).should.be.fulfilled;
            playerExists.should.be.equal(true);            
            isVirtual = await assets.isVirtualPlayer(playerId).should.be.fulfilled;
            isVirtual.should.be.equal(true);            
            playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry+1);
            playerExists = await assets.playerExists(playerId).should.be.fulfilled;
            playerExists.should.be.equal(false);            
            isVirtual = await assets.isVirtualPlayer(playerId).should.be.rejected;
        }
    });

    it('isBot teams', async () =>  {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0;
        isBot = await assets.isBotTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
        isBot.should.be.equal(true);            
    });

    it('transfer first bot to address', async () => {
        const tz = 1;
        const countryIdxInTZ = 0;
        const tx = await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "TeamTransfer", (event) => {
            return event.teamId.should.be.bignumber.equal('274877906944') && event.to.should.be.equal(ALICE);
        });
    });
    
    it('transfer of bot teams', async () =>  {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry1 = 0;
        teamIdxInCountry2 = 1;
        teamId1 = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry1);
        teamId2 = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry2);
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        tx = await assets.transferFirstBotToAddr(tz, countryIdxInTZ, BOB).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "TeamTransfer", (event) => {
            return event.teamId.toNumber() == teamId2 && event.to == BOB;
        });
        isBot = await assets.isBotTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry1).should.be.fulfilled;
        isBot.should.be.equal(false);
        isBot = await assets.isBotTeam(teamId2).should.be.fulfilled;
        isBot.should.be.equal(false);
        owner = await assets.getOwnerTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry1).should.be.fulfilled;
        owner.should.be.equal(ALICE);
        owner = await assets.getOwnerTeam(teamId2).should.be.fulfilled;
        owner.should.be.equal(BOB);
    });

    it('get team player ids', async () => {
        // for the first team we should find playerIdx = [0, 1,...,17, FREE, FREE, ...]
        teamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0);
        let ids = await assets.getPlayerIdsInTeam(teamId).should.be.fulfilled;
        ids.length.should.be.equal(PLAYERS_PER_TEAM_MAX);
        for (shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_MAX; shirtNum++) {
            if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
                ids[shirtNum].should.be.bignumber.equal(FREE_PLAYER_ID);
                continue;
            } else {
                decoded = await assets.decodeTZCountryAndVal(ids[shirtNum]).should.be.fulfilled;
                const {0: timeZone, 1: country, 2: playerIdxInCountry} = decoded;
                playerIdxInCountry.toNumber().should.be.equal(shirtNum);
            }
        }
        // for the first team we should find playerIdx = [18, 19,..., FREE, FREE, ...]
        teamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 1);
        ids = await assets.getPlayerIdsInTeam(teamId).should.be.fulfilled;
        ids.length.should.be.equal(PLAYERS_PER_TEAM_MAX);
        for (shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_MAX; shirtNum++) {
            if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
                ids[shirtNum].should.be.bignumber.equal(FREE_PLAYER_ID);
                continue;
            } else {
                decoded = await assets.decodeTZCountryAndVal(ids[shirtNum]).should.be.fulfilled;
                const {0: timeZone, 1: country, 2: playerIdxInCountry} = decoded;
                playerIdxInCountry.toNumber().should.be.equal(shirtNum + PLAYERS_PER_TEAM_INIT);
            }
        }
    });

    it('gameDeployMonth', async () => {
        const gameDeployMonth =  await assets.gameDeployMonth().should.be.fulfilled;
        currentBlockNum = await web3.eth.getBlockNumber()
        currentBlock = await web3.eth.getBlock(currentBlockNum)
        currentMonth = Math.floor(currentBlock.timestamp * 12 / (3600 * 24 * 365));
        gameDeployMonth.toNumber().should.be.equal(currentMonth);
    });

    it('get skills of a GoalKeeper on creation', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        playerIdxInCountry = 1;
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        encodedSkills = await assets.getPlayerSkillsAtBirth(playerId).should.be.fulfilled;
        skills = await assets.getSkillsVec(encodedSkills).should.be.fulfilled; 
        expected = [44, 0, 0, 0, 0];
        for (sk = 0; sk < N_SKILLS; sk++) {
            skills[sk].toNumber().should.be.equal(expected[sk])
        }
        newId =  await assets.getPlayerIdFromSkills(encodedSkills).should.be.fulfilled; 
        newId.should.be.bignumber.equal(playerId);
        monthOfBirth =  await assets.getMonthOfBirth(encodedSkills).should.be.fulfilled; 
        monthOfBirth.toNumber().should.be.equal(236);
        ageInMonths = await assets.getPlayerAgeInMonths(playerId).should.be.fulfilled;
        ageInMonths.toNumber().should.be.equal(360);
    });

    it('get state of player on creation', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        // test for players on the first team
        playerIdxInCountry = 1;
        teamIdxInCountry = Math.floor(playerIdxInCountry / PLAYERS_PER_TEAM_INIT);
        teamIdxInCountry.should.be.equal(0);
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await assets.getPlayerState(playerId).should.be.fulfilled;
        newId =  await assets.getPlayerIdFromState(state).should.be.fulfilled; 
        newId.should.be.bignumber.equal(playerId);
        expectedTeamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        teamId =  await assets.getCurrentTeamId(state).should.be.fulfilled; 
        teamId.should.be.bignumber.equal(expectedTeamId);
        shirtNum =  await assets.getCurrentShirtNum(state).should.be.fulfilled; 
        shirtNum.toNumber().should.be.equal(1);
        // test for players on the second team
        playerIdxInCountry = 18;
        teamIdxInCountry = Math.floor(playerIdxInCountry / PLAYERS_PER_TEAM_INIT);
        teamIdxInCountry.should.be.equal(1);
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await assets.getPlayerState(playerId).should.be.fulfilled;
        newId =  await assets.getPlayerIdFromState(state).should.be.fulfilled; 
        newId.should.be.bignumber.equal(playerId);
        expectedTeamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        teamId =  await assets.getCurrentTeamId(state).should.be.fulfilled; 
        teamId.should.be.bignumber.equal(expectedTeamId);
        shirtNum =  await assets.getCurrentShirtNum(state).should.be.fulfilled; 
        shirtNum.toNumber().should.be.equal(0);
    });
    
    it('get player state of unexistent player', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT - 1; // last player that exists
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await assets.getPlayerState(playerId).should.be.fulfilled;
        playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT; // player not existing
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await assets.getPlayerState(playerId).should.be.rejected;
        tz = 0; // dummy timeZone
        countryIdxInTZ = 0;
        playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT - 1;
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await assets.getPlayerState(playerId).should.be.rejected;
        tz = 1; 
        countryIdxInTZ = 1; // country not existing
        playerIdxInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT - 1;
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await assets.getPlayerState(playerId).should.be.rejected;
    });


    it('isFreeShirt', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0; 
        teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        // cannot query about a Bot Team
        isFree = await assets.isFreeShirt(teamId,shirtNum = 3).should.be.rejected
        // so transfer and query again
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        isBot = await assets.isBotTeam(teamId).should.be.fulfilled;
        isBot.should.be.equal(false);
        isFree = await assets.isFreeShirt(teamId, shirtNum = 3).should.be.fulfilled
        isFree.should.be.equal(false)
        isFree = await assets.isFreeShirt(teamId, shirtNum = 18).should.be.fulfilled
        isFree.should.be.equal(true)
    });

    it('getFreeShirt', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0; 
        teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        // cannot query about a Bot Team
        shirtNum = await assets.getFreeShirt(teamId).should.be.rejected
        // so transfer and query again
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        isBot = await assets.isBotTeam(teamId).should.be.fulfilled;
        isBot.should.be.equal(false);
        shirtNum = await assets.getFreeShirt(teamId).should.be.fulfilled
        shirtNum.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 1);
    });

    
    it('transferPlayer', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 

        // state before selling:
        state = await assets.getPlayerState(playerId).should.be.fulfilled;
        obtainedTeamId = await assets.getCurrentTeamId(state).should.be.fulfilled;
        obtainedTeamId.should.be.bignumber.equal(teamId1);
        shirt = await assets.getCurrentShirtNum(state).should.be.fulfilled;
        shirt.toNumber().should.be.equal(playerIdxInCountry1);        

        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, BOB).should.be.fulfilled;
        tx = await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;
        // teamId2.should.be.equal(false)
        truffleAssert.eventEmitted(tx, "PlayerTransfer", (event) => {
            return event.playerId == playerId.toNumber() && event.teamIdTarget == teamId2.toNumber();
        });

        // state of player after selling:
        state = await assets.getPlayerState(playerId).should.be.fulfilled;
        obtainedTeamId = await assets.getCurrentTeamId(state).should.be.fulfilled;
        obtainedTeamId.should.be.bignumber.equal(teamId2);
        shirt = await assets.getCurrentShirtNum(state).should.be.fulfilled;
        shirt.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 1);        

        // states of teams after selling
        isFree = await assets.isFreeShirt(teamId1, shirtNum = playerIdxInCountry1).should.be.fulfilled
        isFree.should.be.equal(true);
        isFree = await assets.isFreeShirt(teamId2, shirtNum = PLAYERS_PER_TEAM_MAX - 1).should.be.fulfilled
        isFree.should.be.equal(false);
        shirtNum = await assets.getFreeShirt(teamId2).should.be.fulfilled
        shirtNum.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 2);
    });

    it('get owner of player', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 

        // state before selling:
        owner = await assets.getOwnerPlayer(playerId).should.be.fulfilled;
        owner.should.be.equal(NULL_ADDR);
        // state after acquiring bot:
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        owner = await assets.getOwnerPlayer(playerId).should.be.fulfilled
        owner.should.be.equal(ALICE);
        // state after selling player:
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, BOB).should.be.fulfilled;
        await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;
        owner = await assets.getOwnerPlayer(playerId).should.be.fulfilled;
        owner.should.be.equal(BOB);
        // state after selling team:
        await assets.transferTeam(teamId2, CAROL).should.be.fulfilled;
        owner = await assets.getOwnerPlayer(playerId).should.be.fulfilled;
        owner.should.be.equal(CAROL);
    });

    it('get owner invalid player', async () => {
        await assets.getOwnerPlayer(playerId = 3).should.be.rejected;
    });

    it('transferPlayer different team works', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, ALICE).should.be.fulfilled;
        await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;
    });

    it('transferPlayer same team fails', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferPlayer(playerId, teamId1).should.be.rejected;
    });

    it('transferPlayer fails when at least one team involved is a bot', async () => {
        playerId1   = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        playerId2   = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ1 = 0, playerIdxInCountry1 = 8).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2, countryIdxInTZ2, teamIdxInCountry = 0).should.be.fulfilled; 
        // both teams are bots: fails
        await assets.transferPlayer(playerId1, teamId2).should.be.rejected;
        // only buyer team is bot: fails
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferPlayer(playerId1, teamId2).should.be.rejected;
        // only seller team is bot: fails
        await assets.transferPlayer(playerId2, teamId1).should.be.rejected;
        // // both are owned: works
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, ALICE).should.be.fulfilled;
        await assets.transferPlayer(playerId1, teamId2).should.be.fulfilled;
        await assets.transferPlayer(playerId2, teamId1).should.be.fulfilled;
    });
    
    it('transferPlayer to already full team', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz2, countryIdxInTZ2, teamIdxInCountry = 0).should.be.fulfilled; 
        for (playerIdxInCountry = 0; playerId < PLAYERS_PER_TEAM_MAX-PLAYERS_PER_TEAM_INIT; playerId++) {
            playerId   = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry).should.be.fulfilled; 
            await assets.transferPlayer(playerId, teamId).should.be.fulfilled;
        }
        playerId   = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry+1).should.be.fulfilled; 
        await assets.transferPlayer(playerId, teamId).should.be.rejected;
    });

    it('team exists', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        result = await assets.teamExists(teamId).should.be.fulfilled;
        result.should.be.equal(true);
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = TEAMS_PER_LEAGUE * LEAGUES_PER_DIV - 1).should.be.fulfilled; 
        result = await assets.teamExists(teamId).should.be.fulfilled;
        result.should.be.equal(true);
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = TEAMS_PER_LEAGUE * LEAGUES_PER_DIV).should.be.fulfilled; 
        result = await assets.teamExists(teamId).should.be.fulfilled;
        result.should.be.equal(false);
        teamId     = await assets.encodeTZCountryAndVal(tz = 0, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        result = await assets.teamExists(teamId).should.be.rejected;
    });

    it('initial number of countries', async () => {
        const count = await assets.countCountries(tz = 1).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('initial number of teams', async () => {
        const count = await assets.countTeams(tz = 1, countryIdxInTZ = 0).should.be.fulfilled;
        count.toNumber().should.be.equal(TEAMS_PER_LEAGUE * LEAGUES_PER_DIV);
    });

    it('existence of null player', async () => {
        const exists = await assets.playerExists(playerId = 0).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('isVirtual after sale', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, ALICE).should.be.fulfilled;
        isVirtual = await assets.isVirtualPlayer(playerId).should.be.fulfilled;
        isVirtual.should.be.equal(true)
        await assets.transferPlayer(playerId, teamId2).should.be.fulfilled;
        isVirtual = await assets.isVirtualPlayer(playerId).should.be.fulfilled;
        isVirtual.should.be.equal(false)
    });

    
    it('computed skills with rnd = 0 for a goal keeper', async () => {
        let computedSkills = await assets.computeSkills(rnd = 0, shirtNum = 0).should.be.fulfilled;
        const {0: skills, 1: potential, 2: prefPos} = computedSkills;
        expected = [30, 0, 0, 0, 0];
        for (sk = 0; sk < N_SKILLS; sk++) {
            skills[sk].toNumber().should.be.equal(expected[sk]);
        }
        potential.toNumber().should.be.equal(0);
        prefPos.toNumber().should.be.equal(0);
    });

    it('computed skills with rnd = 0 for non goal keepers should be 50 each', async () => {
        let computedSkills = await assets.computeSkills(rnd = 0, shirtNum = 3).should.be.fulfilled;
        const {0: skills, 1: potential, 2: prefPos} = computedSkills;
        expected = [50, 50, 50, 50, 50];
        for (sk = 0; sk < N_SKILLS; sk++) {
            skills[sk].toNumber().should.be.equal(expected[sk]);
        }
        potential.toNumber().should.be.equal(0);
        decodedPrefPos = await assets.decodePrefPos(prefPos).should.be.fulfilled;
        decodedPrefPos[0].toNumber().should.be.equal(1); // defender
        decodedPrefPos[1].toNumber().should.be.equal(1 + shirtNum);
    });


    it('computed prefPos gives correct number of defenders, mids, etc', async () => {
        expectedPos = [ 0, 0, 0, 1, 1, 1, 1, 1, 2, 2, 4, 4, 5, 5, 3, 3, 3, 3 ];
        for (let shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_INIT; shirtNum++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + shirtNum));
            computedSkills = await assets.computeSkills(seed, shirtNum).should.be.fulfilled;
            decodedPrefPos = await assets.decodePrefPos(computedSkills[2]).should.be.fulfilled;
            decodedPrefPos[0].toNumber().should.be.equal(expectedPos[shirtNum]);
            // skills = computedSkills[0];
            // for (sk = 0; sk < N_SKILLS; sk++) console.log(shirtNum, ": ", skills[sk].toNumber());
        }
    });


    it('sum of computed skills is close to 250', async () => {
        for (let i = 0; i < 10; i++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + i));
            shirtNum = 3 + (seed % 15); // avoid goalkeepers
            computedSkills = await assets.computeSkills(seed, shirtNum).should.be.fulfilled;
            skills = computedSkills[0];
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            (Math.abs(sum - 250) < 5).should.be.equal(true);
        }
    });

    it('get shirtNum in team for all players in a country', async () => {
        tz = 0;
        countryIdxInTZ = 0;
        playersInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT
        for (let playerIdxInCountry = 0; playerId < playersInCountry ; playerId++){
            playerId    = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
            const playerState = await assets.getPlayerState(playerId).should.be.fulfilled;
            const shirtNum = await assets.getCurrentShirtNum(playerState).should.be.fulfilled;
            shirtNum.toNumber().should.be.equal(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        }
    })

    it('transfer team', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE); 
        currentOwner = await assets.getOwnerTeam(teamId).should.be.fulfilled;
        currentOwner.should.be.equal(ALICE);
        tx = await assets.transferTeam(teamId, BOB).should.be.fulfilled;
        newOwner = await assets.getOwnerTeam(teamId).should.be.fulfilled;
        newOwner.should.be.equal(BOB);
        truffleAssert.eventEmitted(tx, "TeamTransfer", (event) => {
            return event.teamId.toNumber() == teamId && event.to == BOB;
        });
    });

    it ('transfer invalid team 0', async () => {
        await assets.transferTeam(teamId = 0, BOB).should.be.rejected;
    });
        
    it('transfer fails when team is a bot', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferTeam(teamId, BOB).should.be.rejected;
        await assets.transferFirstBotToAddr(tz,countryIdxInTZ, ALICE); 
        await assets.transferTeam(teamId, BOB).should.be.fulfilled;
    });

    it('transfer team accross same owner', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE); 
        await assets.transferTeam(teamId, ALICE).should.be.rejected;
    });
});