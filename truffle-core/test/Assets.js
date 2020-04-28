const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const delegateUtils = require('../utils/delegateCallUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');



contract('Assets', (accounts) => {
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const N_SKILLS = 5;
    let N_DIVS_AT_START;
    let initTx = null;
    let N_TEAMS_AT_START;
    
    // Skills: shoot, speed, pass, defence, endurance
    const SK_SHO = 0;
    const SK_SPE = 1;
    const SK_PAS = 2;
    const SK_DEF = 3;
    const SK_END = 4;
    
    var assets;
    var market;
    

    
    const it2 = async(text, f) => {};
    function toBytes32(name) { return web3.utils.utf8ToHex(name); }

    beforeEach(async () => {
        depl = await delegateUtils.deploy(versionNumber = 0, Proxy, proxyAddress = '0x0', Assets, Market, Updates);
        proxy = depl[0]
        assets = depl[1]
        market = depl[2]
        
        constants = await ConstantsGetters.new().should.be.fulfilled;
        initTx = await assets.init().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = await constants.get_PLAYERS_PER_TEAM_INIT().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await constants.get_PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        LEAGUES_PER_DIV = await constants.get_LEAGUES_PER_DIV().should.be.fulfilled;
        TEAMS_PER_LEAGUE = await constants.get_TEAMS_PER_LEAGUE().should.be.fulfilled;
        FREE_PLAYER_ID = await constants.get_FREE_PLAYER_ID().should.be.fulfilled;
        NULL_ADDR = await constants.get_NULL_ADDR().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = PLAYERS_PER_TEAM_INIT.toNumber();
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        LEAGUES_PER_DIV = LEAGUES_PER_DIV.toNumber();
        TEAMS_PER_LEAGUE = TEAMS_PER_LEAGUE.toNumber();
        
        N_DIVS_AT_START = await assets.getNDivisionsInCountry(1,0).should.be.fulfilled;;
        N_DIVS_AT_START = N_DIVS_AT_START.toNumber();
        N_TEAMS_AT_START = N_DIVS_AT_START * LEAGUES_PER_DIV * TEAMS_PER_LEAGUE;
    });
        
    it('create special players', async () => {
        sk = [16383, 13, 4, 56, 456]
        sumSkills = sk.reduce((a, b) => a + b, 0);
        specialPlayerId = await assets.encodePlayerSkills(
            sk,
            dayOfBirth = 4*365, 
            generation = 0,
            playerId = 144321433,
            [potential = 5,
            forwardness = 3,
            leftishness = 4,
            aggressiveness = 1],
            alignedEndOfLastHalf = true,
            redCardLastGame = true,
            gamesNonStopping = 2,
            injuryWeeksLeft = 6,
            substitutedLastHalf = true,
            sumSkills
        ).should.be.fulfilled;
        result = await assets.getPlayerSkillsAtBirth(specialPlayerId).should.be.rejected;
        specialPlayerId = await assets.addIsSpecial(specialPlayerId).should.be.fulfilled;
        skills = await assets.getPlayerSkillsAtBirth(specialPlayerId).should.be.fulfilled;
        result = await assets.getSkill(skills, SK_SHO).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[0]);        
    });

    it('check DivisionCreation event on init', async () => {
        let timezone = 0;
        truffleAssert.eventEmitted(initTx, "DivisionCreation", (event) => {
            timezone++;
            return event.timezone.toString() === timezone.toString() && event.countryIdxInTZ.toString() === '0' && event.divisionIdxInCountry.toString() === '0';
        });
    });

    it('check DivisionCreation event on initSingleTz', async () => {
        const {0: proxy2, 1: assets2, 2: markV0, 3: updV0} =  await delegateUtils.deploy(versionNumber = 0, Proxy, '0x0', Assets, Market, Updates);
        tx = await assets2.initSingleTZ(tz = 4).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "DivisionCreation", (event) => {
            return event.timezone.toString() === tz.toString() && event.countryIdxInTZ.toString() === '0' && event.divisionIdxInCountry.toString() === '0';
        });
    });
    
    
    it('check cannot initialize contract twice', async () => {
        await assets.init().should.be.rejected;
    });

    it('emit event upon creation', async () => {
        truffleAssert.eventEmitted(initTx, "AssetsInit", (event) => {
            return event.creatorAddr.should.be.equal(accounts[0]);
        });
    });

    it('check initial and max number of players per team', async () =>  {
        PLAYERS_PER_TEAM_INIT.should.be.equal(18);
        PLAYERS_PER_TEAM_MAX.should.be.equal(25);
        LEAGUES_PER_DIV.should.be.equal(16);
        TEAMS_PER_LEAGUE.should.be.equal(8);
    });

    it('check initial setup of timeZones', async () =>  {
        nCountries = await assets.countCountries(0).should.be.rejected;
        nCountries = await assets.countCountries(25).should.be.rejected;
        for (tz = 1; tz<25; tz++) {
            nCountries = await assets.countCountries(tz).should.be.fulfilled;
            nCountries.toNumber().should.be.equal(1);
            nDivs = await assets.getNDivisionsInCountry(tz, countryIdxInTZ = 0).should.be.fulfilled;
            nDivs.toNumber().should.be.equal(N_DIVS_AT_START);
            nLeagues = await assets.getNLeaguesInCountry(tz, countryIdxInTZ).should.be.fulfilled;
            nLeagues.toNumber().should.be.equal(N_DIVS_AT_START*LEAGUES_PER_DIV);
            nTeams = await assets.getNTeamsInCountry(tz, countryIdxInTZ).should.be.fulfilled;
            nTeams.toNumber().should.be.equal(N_DIVS_AT_START*LEAGUES_PER_DIV * TEAMS_PER_LEAGUE);
        }
    });

    it('check teamWasCreatedVirtually for existing teams', async () =>  {
        countryIdxInTZ = 0;
        teamIdxInCountry = N_TEAMS_AT_START - 1;
        for (tz = 1; tz<25; tz++) {
            teamWasCreatedVirtually = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
            teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
            teamExists2 = await market.teamWasCreatedVirtually(teamId).should.be.fulfilled;
            teamWasCreatedVirtually.should.be.equal(true);            
            teamExists2.should.be.equal(true); 
        }
    });
    
    it('check teamWasCreatedVirtually for not-created teams', async () =>  {
        countryIdxInTZ = 0;
        teamIdxInCountry = N_TEAMS_AT_START;
        for (tz = 1; tz<25; tz++) {
            teamWasCreatedVirtually = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
            teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
            teamExists2 = await market.teamWasCreatedVirtually(teamId).should.be.fulfilled;
            teamWasCreatedVirtually.should.be.equal(false);            
            teamExists2.should.be.equal(false); 
        }
    });
    
    it('check teamWasCreatedVirtually for non-existing countries', async () =>  {
        countryIdxInTZ = 1;
        teamIdxInCountry = N_TEAMS_AT_START;
        for (tz = 1; tz<25; tz++) {
            teamWasCreatedVirtually = await assets._teamExistsInCountry(tz, countryIdxInTZ, teamIdxInCountry).should.be.rejected;
            teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
            teamExists2 = await market.teamWasCreatedVirtually(teamId).should.be.rejected;
        }
    });

    it('check playerExists and isPlayerWritten', async () =>  {
        countryIdxInTZ = 0;
        teamIdxInCountry = N_TEAMS_AT_START;
        playerIdxInCountry = teamIdxInCountry * PLAYERS_PER_TEAM_INIT - 1;
        for (tz = 1; tz<25; tz++) {
            playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry);
            playerExists = await assets.playerExists(playerId).should.be.fulfilled;
            playerExists.should.be.equal(true);            
            isPlayerWritten = await market.isPlayerWritten(playerId).should.be.fulfilled;
            isPlayerWritten.should.be.equal(false);            
            playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry+1);
            playerExists = await assets.playerExists(playerId).should.be.fulfilled;
            playerExists.should.be.equal(false);            
            isPlayerWritten = await market.isPlayerWritten(playerId).should.be.fulfilled;
            isPlayerWritten.should.be.equal(false);            
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

    it('add users until you need a new division (it can take several seconds)', async () => {
        const tz = 1;
        const countryIdxInTZ = 0;
        nTeamsPerDiv = 128
        for (user = 0; user < (nTeamsPerDiv - 1); user++) {
            await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        }
        tx = await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "DivisionCreation", (event) => {
            return event.timezone.toString() === tz.toString() && event.countryIdxInTZ.toString() === countryIdxInTZ.toString() && event.divisionIdxInCountry.toString() === '1';
        });

    });


    it('transfer 2 bots to address to estimate cost', async () => {
        const tz = 1;
        const countryIdxInTZ = 0;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, BOB).should.be.fulfilled;
    });



    it('transfer of bot teams', async () =>  {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry1 = 0;
        teamIdxInCountry2 = 1;
        teamId1 = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry1);
        teamId2 = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry2);
        addresses = [ALICE, BOB];
        teamIds = [teamId1, teamId2];
        tx = await assets.transferFirstBotsToAddresses([tz, tz], [countryIdxInTZ, countryIdxInTZ], addresses).should.be.fulfilled;
        let count = -1;
        truffleAssert.eventEmitted(tx, "TeamTransfer", (event) => {
            count++;
            return event.teamId.toNumber() == teamIds[count] && event.to == addresses[count];
        });
        isBot = await assets.isBotTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry1).should.be.fulfilled;
        isBot.should.be.equal(false);
        isBot = await market.isBotTeam(teamId2).should.be.fulfilled;
        isBot.should.be.equal(false);
        owner = await assets.getOwnerTeamInCountry(tz, countryIdxInTZ, teamIdxInCountry1).should.be.fulfilled;
        owner.should.be.equal(ALICE);
        owner = await market.getOwnerTeam(teamId2).should.be.fulfilled;
        owner.should.be.equal(BOB);
    });

    it('get team player ids', async () => {
        // for the first team we should find playerIdx = [0, 1,...,17, FREE, FREE, ...]
        teamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0);
        let ids = await market.getPlayerIdsInTeam(teamId).should.be.fulfilled;
        ids.length.should.be.equal(PLAYERS_PER_TEAM_MAX);
        for (shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_MAX; shirtNum++) {
            if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
                ids[shirtNum].toNumber().should.be.equal(0);
            } else {
                decoded = await assets.decodeTZCountryAndVal(ids[shirtNum]).should.be.fulfilled;
                const {0: timeZone, 1: country, 2: playerIdxInCountry} = decoded;
                playerIdxInCountry.toNumber().should.be.equal(shirtNum);
            }
        }
        // for the first team we should find playerIdx = [18, 19,..., FREE, FREE, ...]
        teamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 1);
        ids = await market.getPlayerIdsInTeam(teamId).should.be.fulfilled;
        ids.length.should.be.equal(PLAYERS_PER_TEAM_MAX);
        for (shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_MAX; shirtNum++) {
            if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
                ids[shirtNum].toNumber().should.be.equal(0);
            } else {
                decoded = await assets.decodeTZCountryAndVal(ids[shirtNum]).should.be.fulfilled;
                const {0: timeZone, 1: country, 2: playerIdxInCountry} = decoded;
                playerIdxInCountry.toNumber().should.be.equal(shirtNum + PLAYERS_PER_TEAM_INIT);
            }
        }
    });

    it('gameDeployDay', async () => {
        const gameDeployDay =  await assets.gameDeployDay().should.be.fulfilled;
        currentBlockNum = await web3.eth.getBlockNumber()
        currentBlock = await web3.eth.getBlock(currentBlockNum)
        currentDay = Math.floor(currentBlock.timestamp / (3600 * 24));
        gameDeployDay.toNumber().should.be.equal(currentDay);
    });

    it('get skills of a GoalKeeper on creation', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        playerIdxInCountry = 1;
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        encodedSkills = await assets.getPlayerSkillsAtBirth(playerId).should.be.fulfilled;
        expectedSkills = [ 1427, 1016, 853, 974, 726 ];
        resultSkills = [];
        for (sk = 0; sk < N_SKILLS; sk++) {
            resultSkills.push(await assets.getSkill(encodedSkills, sk).should.be.fulfilled);
        }
        debug.compareArrays(resultSkills, expectedSkills, toNum = true, verbose = false);

        newId =  await assets.getPlayerIdFromSkills(encodedSkills).should.be.fulfilled; 
        newId.should.be.bignumber.equal(playerId);
        gameDeployDay = await assets.gameDeployDay().should.be.fulfilled;
        dayOfBirth =  await assets.getBirthDay(encodedSkills).should.be.fulfilled; 
        ageInDays = await assets.getPlayerAgeInDays(playerId).should.be.fulfilled;
        (Math.abs(ageInDays.toNumber() - 11521) <= 7).should.be.equal(true); // we cannot guarantee exactness +/- 1
        // check that the ageInDay can be obtained by 7 * (now - dayOfBirth), where
        // now is approximately gameDeployDay. There is an uncertainty of about 7 days due to rounding.
        (Math.abs(7*(gameDeployDay.toNumber()-dayOfBirth.toNumber())-ageInDays) < 8).should.be.equal(true);
    });

    it('get state of player on creation', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        // test for players on the first team
        playerIdxInCountry = 1;
        teamIdxInCountry = Math.floor(playerIdxInCountry / PLAYERS_PER_TEAM_INIT);
        teamIdxInCountry.should.be.equal(0);
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await market.getPlayerState(playerId).should.be.fulfilled;
        newId =  await assets.getPlayerIdFromState(state).should.be.fulfilled; 
        newId.should.be.bignumber.equal(playerId);
        expectedTeamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        teamId =  await assets.getCurrentTeamIdFromPlayerState(state).should.be.fulfilled; 
        teamId.should.be.bignumber.equal(expectedTeamId);
        shirtNum =  await assets.getCurrentShirtNum(state).should.be.fulfilled; 
        shirtNum.toNumber().should.be.equal(1);
        // test for players on the second team
        playerIdxInCountry = 18;
        teamIdxInCountry = Math.floor(playerIdxInCountry / PLAYERS_PER_TEAM_INIT);
        teamIdxInCountry.should.be.equal(1);
        playerId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
        state = await market.getPlayerState(playerId).should.be.fulfilled;
        newId =  await assets.getPlayerIdFromState(state).should.be.fulfilled; 
        newId.should.be.bignumber.equal(playerId);
        expectedTeamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        teamId =  await assets.getCurrentTeamIdFromPlayerState(state).should.be.fulfilled; 
        teamId.should.be.bignumber.equal(expectedTeamId);
        shirtNum =  await assets.getCurrentShirtNum(state).should.be.fulfilled; 
        shirtNum.toNumber().should.be.equal(0);
    });

    it('isFreeShirt', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0; 
        teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled;
        let ids = await market.getPlayerIdsInTeam(teamId).should.be.fulfilled;
        shirtNum = 3;
        isFree = await market.isFreeShirt(ids[shirtNum], shirtNum = 18).should.be.fulfilled
        isFree.should.be.equal(false);
        shirtNum = 18;
        isFree = await market.isFreeShirt(ids[shirtNum], shirtNum = 18).should.be.fulfilled
        isFree.should.be.equal(true);
    });

    it('getFreeShirt', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0; 
        teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry).should.be.fulfilled; 
        shirtNum = await market.getFreeShirt(teamId).should.be.fulfilled
        shirtNum.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 1);
    });

    
    it('transferPlayer', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 

        // state before selling:
        state = await market.getPlayerState(playerId).should.be.fulfilled;
        obtainedTeamId = await assets.getCurrentTeamIdFromPlayerState(state).should.be.fulfilled;
        obtainedTeamId.should.be.bignumber.equal(teamId1);
        shirt = await assets.getCurrentShirtNum(state).should.be.fulfilled;
        shirt.toNumber().should.be.equal(playerIdxInCountry1);        

        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, BOB).should.be.fulfilled;
        tx = await market.transferPlayer(playerId, teamId2).should.be.fulfilled;

        // state of player after selling:
        state = await market.getPlayerState(playerId).should.be.fulfilled;
        obtainedTeamId = await assets.getCurrentTeamIdFromPlayerState(state).should.be.fulfilled;
        obtainedTeamId.should.be.bignumber.equal(teamId2);
        shirt = await assets.getCurrentShirtNum(state).should.be.fulfilled;
        shirt.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 1);        

        
        truffleAssert.eventEmitted(tx, "PlayerStateChange", (event) => {
            return event.playerId.should.be.bignumber.equal(playerId) && event.state.should.be.bignumber.equal(state);
        });

        // states of teams after selling
        let ids1 = await market.getPlayerIdsInTeam(teamId1).should.be.fulfilled;
        let ids2 = await market.getPlayerIdsInTeam(teamId2).should.be.fulfilled;
        
        shirtNum = playerIdxInCountry1;
        isFree = await market.isFreeShirt(ids1[shirtNum], shirtNum).should.be.fulfilled
        isFree.should.be.equal(true);
        shirtNum = PLAYERS_PER_TEAM_MAX - 1;
        isFree = await market.isFreeShirt(ids2[shirtNum], shirtNum).should.be.fulfilled
        isFree.should.be.equal(false);
        shirtNum = await market.getFreeShirt(teamId2).should.be.fulfilled
        shirtNum.toNumber().should.be.equal(PLAYERS_PER_TEAM_MAX - 2);
    });

    it('get owner of player', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 

        // state before selling:
        owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
        owner.should.be.equal(NULL_ADDR);
        // state after acquiring bot:
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        owner = await market.getOwnerPlayer(playerId).should.be.fulfilled
        owner.should.be.equal(ALICE);
        // state after selling player:
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, BOB).should.be.fulfilled;
        await market.transferPlayer(playerId, teamId2).should.be.fulfilled;
        owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
        owner.should.be.equal(BOB);
        // state after selling team:
        await market.transferTeam(teamId2, CAROL).should.be.fulfilled;
        owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
        owner.should.be.equal(CAROL);
    });

    it('get owner invalid player', async () => {
        await market.getOwnerPlayer(playerId = 3).should.be.rejected;
    });

    it('transferPlayer different team works', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, ALICE).should.be.fulfilled;
        await market.transferPlayer(playerId, teamId2).should.be.fulfilled;
    });

    it('transferPlayer same team fails', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await market.transferPlayer(playerId, teamId1).should.be.rejected;
    });

    it('transferPlayer fails when at least one team involved is a bot', async () => {
        playerId1   = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        playerId2   = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ1 = 0, playerIdxInCountry1 = 8).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2, countryIdxInTZ2, teamIdxInCountry = 0).should.be.fulfilled; 
        // both teams are bots: fails
        await market.transferPlayer(playerId1, teamId2).should.be.rejected;
        // only buyer team is bot: fails
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await market.transferPlayer(playerId1, teamId2).should.be.rejected;
        // only seller team is bot: fails
        await market.transferPlayer(playerId2, teamId1).should.be.rejected;
        // // both are owned: works
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, ALICE).should.be.fulfilled;
        await market.transferPlayer(playerId1, teamId2).should.be.fulfilled;
        await market.transferPlayer(playerId2, teamId1).should.be.fulfilled;
    });
    
    it('transferPlayer to already full team', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz2, countryIdxInTZ2, teamIdxInCountry = 0).should.be.fulfilled; 
        for (playerIdxInCountry = 0; playerId < PLAYERS_PER_TEAM_MAX-PLAYERS_PER_TEAM_INIT; playerId++) {
            playerId   = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry).should.be.fulfilled; 
            await market.transferPlayer(playerId, teamId).should.be.fulfilled;
        }
        playerId   = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry+1).should.be.fulfilled; 
        await market.transferPlayer(playerId, teamId).should.be.rejected;
    });

    it('team exists', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        result = await market.teamWasCreatedVirtually(teamId).should.be.fulfilled;
        result.should.be.equal(true);
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = N_DIVS_AT_START * TEAMS_PER_LEAGUE * LEAGUES_PER_DIV - 1).should.be.fulfilled; 
        result = await market.teamWasCreatedVirtually(teamId).should.be.fulfilled;
        result.should.be.equal(true);
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = N_DIVS_AT_START * TEAMS_PER_LEAGUE * LEAGUES_PER_DIV).should.be.fulfilled; 
        result = await market.teamWasCreatedVirtually(teamId).should.be.fulfilled;
        result.should.be.equal(false);
        teamId     = await assets.encodeTZCountryAndVal(tz = 0, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        result = await market.teamWasCreatedVirtually(teamId).should.be.rejected;
    });

    it('initial number of countries', async () => {
        const count = await assets.countCountries(tz = 1).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('initial number of teams', async () => {
        const count = await assets.getNTeamsInCountry(tz = 1, countryIdxInTZ = 0).should.be.fulfilled;
        count.toNumber().should.be.equal(N_DIVS_AT_START * TEAMS_PER_LEAGUE * LEAGUES_PER_DIV);
    });

    it('existence of null player', async () => {
        const exists = await assets.playerExists(playerId = 0).should.be.fulfilled;
        exists.should.be.equal(false);
    });

    it('isPlayerWritten after sale', async () => {
        playerId    = await assets.encodeTZCountryAndVal(tz1 = 1, countryIdxInTZ1 = 0, playerIdxInCountry1 = 3).should.be.fulfilled; 
        teamId1     = await assets.encodeTZCountryAndVal(tz1, countryIdxInTZ1, teamIdxInCountry = 0).should.be.fulfilled; 
        teamId2     = await assets.encodeTZCountryAndVal(tz2 = 2, countryIdxInTZ2 = 0, teamIdxInCountry = 0).should.be.fulfilled; 
        await assets.transferFirstBotToAddr(tz1, countryIdxInTZ1, ALICE).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz2, countryIdxInTZ2, ALICE).should.be.fulfilled;
        isPlayerWritten = await market.isPlayerWritten(playerId).should.be.fulfilled;
        isPlayerWritten.should.be.equal(false)
        await market.transferPlayer(playerId, teamId2).should.be.fulfilled;
        isPlayerWritten = await market.isPlayerWritten(playerId).should.be.fulfilled;
        isPlayerWritten.should.be.equal(true)
    });

    
    it('computed skills with rnd = 0 for a goal keeper', async () => {
        let computedSkills = await assets.computeSkills(rnd = 0, shirtNum = 0).should.be.fulfilled;
        const {0: skills, 1: birthTraits} = computedSkills;
        expected = [1000, 1000, 1000, 1000, 1000];
        debug.compareArrays(skills, expected, toNum = true, verbose = false);
        // birthTraits = [potential, forwardness, leftishness, aggressiveness]
        expected = [0, 0, 0, 0] // shirNum = 0 for a GK
        debug.compareArrays(birthTraits, expected, toNum = true, verbose = false);
    });

    
    it('test that goal keepers have great shoot=block skills', async () => {
        skillsAvg = [0,0,0,0,0];
        nTrials = 100;
        for (n = 0; n < nTrials; n++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + n));
            var {0: skills, 1: birthTraits} = await assets.computeSkills(seed , shirtNum = 0).should.be.fulfilled;
            for (sk=0; sk < 5; sk++) skillsAvg[sk] += skills[sk].toNumber();
        }
        for (sk=0; sk < 5; sk++) skillsAvg[sk] = Math.floor(skillsAvg[sk]/nTrials);
        expected = [ 1371, 909, 795, 957, 963 ];
        debug.compareArrays(skillsAvg, expected, toNum = false, verbose = false);
    });

    it('test that forwards have great shoot skills', async () => {
        skillsAvg = [0,0,0,0,0];
        nTrials = 100;
        for (n = 0; n < nTrials; n++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + n));
            var {0: skills, 1: birthTraits} = await assets.computeSkills(seed , shirtNum = 16).should.be.fulfilled;
            for (sk=0; sk < 5; sk++) skillsAvg[sk] += skills[sk].toNumber();
        }
        for (sk=0; sk < 5; sk++) skillsAvg[sk] = Math.floor(skillsAvg[sk]/nTrials);
        expected = [ 1251, 950, 989, 802, 1004 ];
        debug.compareArrays(skillsAvg, expected, toNum = false, verbose = false);
    });
    
    it('computed skills with rnd = 0 for non goal keepers should be 1000 each', async () => {
        let computedSkills = await assets.computeSkills(rnd = 0, shirtNum = 3).should.be.fulfilled;
        const {0: skills, 1: birthTraits} = computedSkills;
        expected = [1000, 1000, 1000, 1000, 1000];
        debug.compareArrays(skills, expected, toNum = true, verbose = false);
        // birthTraits = [potential, forwardness, leftishness, aggressiveness]
        expected = [0, 1, 1, 0]
        debug.compareArrays(birthTraits, expected, toNum = true, verbose = false);
    });

    it('computed prefPos gives correct number of defenders, mids, etc', async () => {
        expectedPos = [ 0, 0, 0, 1, 1, 1, 1, 1, 2, 2, 4, 4, 5, 5, 3, 3, 3, 3 ];
        for (let shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_INIT; shirtNum++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + shirtNum));
            computedSkills = await assets.computeSkills(seed, shirtNum).should.be.fulfilled;
            birthTraits = computedSkills[1];
            birthTraits[1].toNumber().should.be.equal(expectedPos[shirtNum]);
        }
    });

    it('testing aggressiveness', async () => { 
        expectedAggr = [ 3, 1, 1, 2, 2, 2, 0, 1, 0, 1, 2, 2, 3, 3, 3, 0, 2, 3 ];
        resultAggr = []
        for (let shirtNum = 0; shirtNum < PLAYERS_PER_TEAM_INIT; shirtNum++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + shirtNum));
            computedSkills = await assets.computeSkills(seed, shirtNum).should.be.fulfilled;
            birthTraits = computedSkills[1];
            resultAggr.push(birthTraits[3])
        }
        debug.compareArrays(resultAggr, expectedAggr, toNum = true, verbose = false);
    });

    it('sum of computed skills is close to 5000', async () => {
        for (let i = 0; i < 10; i++) {
            seed = web3.utils.toBN(web3.utils.keccak256("32123" + i));
            shirtNum = 3 + (seed % 15); // avoid goalkeepers
            computedSkills = await assets.computeSkills(seed, shirtNum).should.be.fulfilled;
            skills = computedSkills[0];
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            (Math.abs(sum - 5000) < 5).should.be.equal(true);
        }
    });

    it('get shirtNum in team for many players in a country', async () => {
        tz = 1;
        countryIdxInTZ = 0;
        playersInCountry = LEAGUES_PER_DIV * TEAMS_PER_LEAGUE * PLAYERS_PER_TEAM_INIT
        for (let playerIdxInCountry = 0; playerIdxInCountry < playersInCountry ; playerIdxInCountry += 77){
            playerId    = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, playerIdxInCountry).should.be.fulfilled; 
            const playerState = await market.getPlayerState(playerId).should.be.fulfilled;
            const shirtNum = await assets.getCurrentShirtNum(playerState).should.be.fulfilled;
            shirtNum.toNumber().should.be.equal(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        }
    })

    it('transfer team', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE); 
        currentOwner = await market.getOwnerTeam(teamId).should.be.fulfilled;
        currentOwner.should.be.equal(ALICE);
        tx = await market.transferTeam(teamId, BOB).should.be.fulfilled;
        newOwner = await market.getOwnerTeam(teamId).should.be.fulfilled;
        newOwner.should.be.equal(BOB);
        truffleAssert.eventEmitted(tx, "TeamTransfer", (event) => {
            return event.teamId.toNumber() == teamId && event.to == BOB;
        });
    });

    it('transfer invalid team 0', async () => {
        await market.transferTeam(teamId = 0, BOB).should.be.rejected;
    });

    it('transfer bot from a not-initialized tz', async () => {
        teamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.fulfilled; 
        teamId = await assets.encodeTZCountryAndVal(tz = 26, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE).should.be.rejected; 
    });

    it('transfer fails when team is a bot', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await market.transferTeam(teamId, BOB).should.be.rejected;
        await assets.transferFirstBotToAddr(tz,countryIdxInTZ, ALICE); 
        await market.transferTeam(teamId, BOB).should.be.fulfilled;
    });

    it('transfer team accross same owner', async () => {
        teamId     = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0).should.be.fulfilled;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE); 
        await market.transferTeam(teamId, ALICE).should.be.rejected;
    });
});