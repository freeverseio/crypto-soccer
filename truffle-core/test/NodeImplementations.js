/*
 Tests node js implementations agains solidity ones
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const timeTravel = require('../utils/TimeTravel.js');
const deployUtils = require('../utils/deployUtils.js');
const { assert } = require('chai');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Challenges = artifacts.require('Challenges');
const Merkle = artifacts.require('Merkle');
const Stakers = artifacts.require("Stakers")
const Utils = artifacts.require('Utils');

const UniverseInfo = artifacts.require('UniverseInfo');
const EncodingSkills = artifacts.require('EncodingSkills');
const EncodingState = artifacts.require('EncodingState');
const EncodingSkillsSetters = artifacts.require('EncodingSkillsSetters');
const UpdatesBase = artifacts.require('UpdatesBase');

const seedrandom = require('seedrandom');
const {nextTimeZoneToPlay, getCurrentRound, getMatch1stHalfUTC, getMatchHalfUTC, calendarInfo} = require('../utils/calendarUtils.js');

contract('Updates', (accounts) => {
    const inheritedArtfcts = [UniverseInfo, EncodingSkills, EncodingState, EncodingSkillsSetters, UpdatesBase];
    const nLevelsInOneChallenge = 11;
    const nNonNullLeafsInLeague = 640;
    const nLevelsInLastChallenge = 10; // must be nearest exponent to 640 ... 1024
    
    const it2 = async(text, f) => {};
    
    async function deployAndConfigureStakers(Stakers, updates, setup) {
        const { singleTimezone, owners, requiredStake } = setup;
        const stakers  = await Stakers.new(updates.address, requiredStake).should.be.fulfilled;

        for (trustedParty of owners.trustedParties) {
            await stakers.addTrustedParty(trustedParty, {from: owners.COO}).should.be.fulfilled;
        }
        for (trustedParty of owners.trustedParties) {
            await stakers.enrol({from:trustedParty, value: requiredStake}).should.be.fulfilled;
        }
        return stakers;
    }
    



    beforeEach(async () => {
        defaultSetup = deployUtils.getDefaultSetup(accounts);
        owners = defaultSetup.owners;
        depl = await deployUtils.deploy(owners, Proxy, Assets, Market, Updates, Challenges, inheritedArtfcts);
        [proxy, assets, market, updates, challenges] = depl;
        await deployUtils.setProxyContractOwners(proxy, assets, owners, owners.company).should.be.fulfilled;
        await updates.setChallengeTime(60, {from: owners.COO}).should.be.fulfilled;
        stakers = await deployAndConfigureStakers(Stakers, updates, defaultSetup);
        await updates.setStakersAddress(stakers.address, {from: owners.superuser}).should.be.fulfilled;
        await stakers.setGameOwner(updates.address, {from:owners.COO}).should.be.fulfilled;
        
        utils = await Utils.new().should.be.fulfilled;
        constants = await ConstantsGetters.new().should.be.fulfilled;
        merkle = await Merkle.new().should.be.fulfilled;
        blockChainTimeSec = Math.floor(Date.now()/1000);
        await updates.initUpdates(blockChainTimeSec, {from: owners.COO}).should.be.fulfilled;
        await updates.setChallengeLevels(nLevelsInOneChallenge, nNonNullLeafsInLeague, nLevelsInLastChallenge, {from: owners.relay}).should.be.fulfilled;
        NULL_TIMEZONE = await constants.get_NULL_TIMEZONE().should.be.fulfilled;
        NULL_TIMEZONE = NULL_TIMEZONE.toNumber();
        VERSES_PER_DAY = await constants.get_VERSES_PER_DAY().should.be.fulfilled;
        VERSES_PER_ROUND = await constants.get_VERSES_PER_ROUND().should.be.fulfilled;
    });

    it('calendarInfo', async () =>  {
        let info;
        const NULL_TIMEZONE = 0;

        info = calendarInfo(verse = 0, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1, "matchDay": 0, "half": 0, "leagueRound": 0, "timestamp": firstVerseTimeStamp});

        info = calendarInfo(verse = 0, TZForRound1 = 14, firstVerseTimeStamp = 55550);
        assert.deepEqual(info, {"timezone": TZForRound1, "matchDay": 0, "half": 0, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = 1, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1, "matchDay": 0, "half": 1, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = 2, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": NULL_TIMEZONE, "matchDay": null, "half": null, "leagueRound": null, "timestamp": null});

        info = calendarInfo(verse = 3, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": NULL_TIMEZONE, "matchDay": null, "half": null, "leagueRound": null, "timestamp": null});

        info = calendarInfo(verse = 4, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1 + 1, "matchDay": 0, "half": 0, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = 5, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1 + 1, "matchDay": 0, "half": 1, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = 6, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": NULL_TIMEZONE, "matchDay": null, "half": null, "leagueRound": null, "timestamp": null});

        info = calendarInfo(verse = 7, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": NULL_TIMEZONE, "matchDay": null, "half": null, "leagueRound": null, "timestamp": null});

        info = calendarInfo(verse = 8, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1 + 2, "matchDay": 0, "half": 0, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = 9, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1 + 2, "matchDay": 0, "half": 1, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = 10, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": NULL_TIMEZONE, "matchDay": null, "half": null, "leagueRound": null, "timestamp": null});

        info = calendarInfo(verse = 11, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": NULL_TIMEZONE, "matchDay": null, "half": null, "leagueRound": null, "timestamp": null});

        // after one leage:
        const VERSES_PER_ROUND = 672; /// 96 * 7days
        info = calendarInfo(verse = VERSES_PER_ROUND, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1, "matchDay": 0, "half": 0, "leagueRound": 1, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = VERSES_PER_ROUND + 1, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": TZForRound1, "matchDay": 0, "half": 1, "leagueRound": 1, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = VERSES_PER_ROUND + 2, TZForRound1 = 1, firstVerseTimeStamp = 0);
        assert.deepEqual(info, {"timezone": 16, "matchDay": 13, "half": 0, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});

        info = calendarInfo(verse = VERSES_PER_ROUND + 2, TZForRound1 = 1, firstVerseTimeStamp = 534535);
        assert.deepEqual(info, {"timezone": 16, "matchDay": 13, "half": 0, "leagueRound": 0, "timestamp": firstVerseTimeStamp + 900 * verse});
    });

    it('TimezonetoUptate bug from field', async () =>  {
        const bcResult = await updates.timeZoneToUpdatePure(12289,24).should.be.fulfilled;
        bcResult.timezone.toNumber().should.be.equal(24);
        bcResult.turnInDay.toNumber().should.be.equal(1);
        bcResult.day.toNumber().should.be.equal(4);

        const nodeResult = nextTimeZoneToPlay(12289,24);
        nodeResult.timezone.should.be.equal(24);
        nodeResult.half.should.be.equal(1);
        nodeResult.matchDay.should.be.equal(4);
    });

    it('TimezonetoUptate bug from field with random inputs', async () =>  {
        const numberOfTests = 100;
        const rng = seedrandom('dummyseed');

        function getRandomValue(min, max) {
            return Math.floor(rng() * (max - min + 1)) + min;
        }
    
        for (let i = 0; i < numberOfTests; i++) {
            const randomVerse = getRandomValue(1, 1000000);
            const randomTZForRound1 = getRandomValue(0, 24); // Timezones range from 1 to 24
    
            const bcResult = await updates.timeZoneToUpdatePure(randomVerse, randomTZForRound1).should.be.fulfilled;
            const nodeResult = nextTimeZoneToPlay(randomVerse, randomTZForRound1);
    
            bcResult.timezone.toNumber().should.be.equal(nodeResult.timezone);
            bcResult.turnInDay.toNumber().should.be.equal(nodeResult.half);
            bcResult.day.toNumber().should.be.equal(nodeResult.matchDay);
        }
    });

    it('test getMatchUTC', async () =>  {
        let bcUTC;

        const firstVerseTimestamp = Number(await updates.getNextVerseTimestamp());
        const TZForRound1 = await updates.getTimeZoneForRound1().should.be.fulfilled;

        bcUTC = await updates.getMatchUTC(tz = TZForRound1, round = 0, matchDay = 0).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));

        bcUTC = await updates.getMatchUTC(tz = TZForRound1, round = 0, matchDay = 2).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp + 24*3600);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));

        bcUTC = await updates.getMatchUTC(tz = TZForRound1, round = 0, matchDay = 1).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp + 9.5*3600);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));

        bcUTC = await updates.getMatchUTC(tz = TZForRound1, round = 1, matchDay = 1).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp + 9.5*3600 + 7*24*3600);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));

        bcUTC = await updates.getMatchUTC(tz = TZForRound1, round = 1, matchDay = 2).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp + 24*3600 + 7*24*3600);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));


        // tests for other timezones
        tz = 1;
        let deltaN = (tz >= TZForRound1) ? (tz-TZForRound1) : (24+tz-TZForRound1); 
        bcUTC = await updates.getMatchUTC(tz, round = 0, matchDay = 0).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp + deltaN * 3600);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));

        tz = 24;
        deltaN = (tz >= TZForRound1) ? (tz-TZForRound1) : (24+tz-TZForRound1); 
        bcUTC = await updates.getMatchUTC(tz, round = 0, matchDay = 0).should.be.fulfilled;
        bcUTC.toNumber().should.be.equal(firstVerseTimestamp + deltaN * 3600);
        bcUTC.toNumber().should.be.equal(getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimestamp));
        bcUTC.toNumber().should.be.equal(getMatchHalfUTC(tz, round, matchDay, half = 0, TZForRound1, firstVerseTimestamp));
    });

    it('test getCurrentRound', async () =>  {
        result = await assets.getCurrentRoundPure(tz = 5, TZForRound1 = 5, verse = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 24, TZForRound1 = 5, verse = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 4, TZForRound1 = 5, verse = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        VERSES_DAY = 24*4;
        VERSES_ROUND = 7 * VERSES_DAY;
        // move to start of round 1 for 1st tz:
        result = await assets.getCurrentRoundPure(tz = 5, TZForRound1 = 5, verse = VERSES_ROUND).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 4, TZForRound1 = 5, verse = VERSES_ROUND).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 24, TZForRound1 = 5, verse = VERSES_ROUND).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        // move to start of round 1 for 1st tz after TZForRound1:
        result = await assets.getCurrentRoundPure(tz = 5, TZForRound1 = 5, verse = VERSES_ROUND + 4).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 6, TZForRound1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 7, TZForRound1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 24, TZForRound1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        // move to start of round 1 for last tz to reach it:
        result = await assets.getCurrentRoundPure(tz = 5, TZForRound1 = 5, verse = 2 * VERSES_ROUND - 4).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 4, TZForRound1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 24, TZForRound1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));

        result = await assets.getCurrentRoundPure(tz = 16, TZForRound1 = 1, verse = 674).should.be.fulfilled;
        result.toNumber().should.be.equal(getCurrentRound(tz, TZForRound1, verse));
        result.toNumber().should.be.equal(0);
    });
});