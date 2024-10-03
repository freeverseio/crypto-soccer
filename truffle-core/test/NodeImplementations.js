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
    

    // Inputs:
    // - verse: the verse to be played
    // - TZForRound1: the timezone that played the very first games
    // Outputs: for the verse to be played:
    // - timezone: which timezone plays
    // - matchDay: what matchDay of the league it corresponds to: a number in [0, 13]
    // - turn: whether it corresponds to the first half (turn = 0), or to the second half (turn = 1)
    function timeZoneToUpdatePure(verse, TZForRound1) {
        const NULL_TIMEZONE = 0; 
        const VERSES_PER_DAY = 96; 
        const MATCHDAYS_PER_ROUND = 14;

        // if _currentVerse = 0, we should be updating TZForRound1
        // recall that timeZones range from 1...24 (not from 0...24)

        let turn = verse % 4;
        let delta = 9 * 4 + turn;
        let tz, dia;

        // if the turn is greater or equal to 2 and verse is less than delta, return NULL_TIMEZONE
        if (turn >= 2 && verse < delta) return { timezone: NULL_TIMEZONE, day: 0, turnInDay: 0 };

        if (turn < 2) {
            tz = TZForRound1 + Math.floor((verse - turn) % VERSES_PER_DAY / 4);
            dia = 2 * Math.floor((verse - 4 * (tz - TZForRound1) - turn) / VERSES_PER_DAY);
        } else {
            tz = TZForRound1 + Math.floor((verse - delta) % VERSES_PER_DAY / 4);
            dia = 1 + 2 * Math.floor((verse - 4 * (tz - TZForRound1) - delta) / VERSES_PER_DAY);
            turn -= 2;
        }

        let timezone = normalizeTZ(tz);
        let day = dia % MATCHDAYS_PER_ROUND;

        return { timezone, day, turnInDay: turn };
    }

    function normalizeTZ(tz) {
        return 1 + ((24 + tz - 1) % 24);
    }

    // Function to get the current round (pure function)
    function getCurrentRound(tz, TZForRound1, verse) {
        const VERSES_PER_ROUND = 672; /// 96 * 7days
        if (verse < VERSES_PER_ROUND) return 0;
        let round = Math.floor(verse / VERSES_PER_ROUND);
        let deltaN = (tz >= TZForRound1) ? (tz - TZForRound1) : ((tz + 24) - TZForRound1);
        if (verse < 4 * deltaN + round * VERSES_PER_ROUND) {
            return round - 1;
        } else {
            return round;
        }
    }


    // Returns the Unix timestamp in UTC (seconds) corresponding to the start of a match's first half 
    // Inputs:
    // - tz: the timezone where the match belongs
    // - round: the round of a league (the first league played is round 0, the next league is round 1, etc.)
    // - matchDay: what matchDay of the league it corresponds to: a number in [0, 13]
    // - TZForRound1: the timezone that played the very first games
    // - firstVerseTimeStamp: the timestamp the very first games where played at 
    // Outputs:
    // - timeUTC: the Unix timestamp in UTC (seconds) corresponding to the start of a match's first half 
    function getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimeStamp) {
        const DAYS_PER_ROUND = 7;
        if (tz <= 0 || tz >= 25) {
            throw new Error("timezone out of range");
        }

        const deltaN = (tz >= TZForRound1) ? 
            (tz - TZForRound1) : 
            ((tz + 24) - TZForRound1);

        let timeUTC;
        if (matchDay % 2 === 0) {
            timeUTC = firstVerseTimeStamp + (deltaN + 12 * matchDay + 24 * DAYS_PER_ROUND * round) * 3600;
        } else {
            timeUTC = firstVerseTimeStamp + (19 + 2 * deltaN + 24 * (matchDay - 1) + 48 * DAYS_PER_ROUND * round) * 1800;
        }
        return timeUTC;
    }

    function getMatchHalfUTC(tz, round, matchDay, half, TZForRound1, firstVerseTimeStamp) {
        const SECS_BETWEEN_VERSES = 900; /// 15 mins
        const extraSeconds = half == 0 ? 0 : SECS_BETWEEN_VERSES;
        return getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimeStamp) + extraSeconds;
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

    it2('TimezonetoUptate bug from field', async () =>  {
        const bcResult = await updates.timeZoneToUpdatePure(12289,24).should.be.fulfilled;
        bcResult.timezone.toNumber().should.be.equal(24);
        bcResult.turnInDay.toNumber().should.be.equal(1);
        bcResult.day.toNumber().should.be.equal(4);

        const nodeResult = timeZoneToUpdatePure(12289,24);
        nodeResult.timezone.should.be.equal(24);
        nodeResult.turnInDay.should.be.equal(1);
        nodeResult.day.should.be.equal(4);
    });

    it2('TimezonetoUptate bug from field with random inputs', async () =>  {
        const numberOfTests = 100;
        const rng = seedrandom('dummyseed');

        function getRandomValue(min, max) {
            return Math.floor(rng() * (max - min + 1)) + min;
        }
    
        for (let i = 0; i < numberOfTests; i++) {
            const randomVerse = getRandomValue(1, 1000000);
            const randomTZForRound1 = getRandomValue(0, 24); // Timezones range from 1 to 24
    
            const bcResult = await updates.timeZoneToUpdatePure(randomVerse, randomTZForRound1).should.be.fulfilled;
            const nodeResult = timeZoneToUpdatePure(randomVerse, randomTZForRound1);
    
            bcResult.timezone.toNumber().should.be.equal(nodeResult.timezone);
            bcResult.turnInDay.toNumber().should.be.equal(nodeResult.turnInDay);
            bcResult.day.toNumber().should.be.equal(nodeResult.day);
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
    });
});