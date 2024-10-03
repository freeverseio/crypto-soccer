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
    

    function timeZoneToUpdatePure(verse, TZForRound1) {
        const NULL_TIMEZONE = 0; 
        const VERSES_PER_DAY = 96; 
        const MATCHDAYS_PER_ROUND = 14;

        // if _currentVerse = 0, we should be updating _timeZoneForRound1
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

    it('TimezonetoUptate bug from field', async () =>  {
        const bcResult = await updates.timeZoneToUpdatePure(12289,24).should.be.fulfilled;
        bcResult.timezone.toNumber().should.be.equal(24);
        bcResult.turnInDay.toNumber().should.be.equal(1);
        bcResult.day.toNumber().should.be.equal(4);

        const nodeResult = timeZoneToUpdatePure(12289,24);
        nodeResult.timezone.should.be.equal(24);
        nodeResult.turnInDay.should.be.equal(1);
        nodeResult.day.should.be.equal(4);
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
            const nodeResult = timeZoneToUpdatePure(randomVerse, randomTZForRound1);
    
            bcResult.timezone.toNumber().should.be.equal(nodeResult.timezone);
            bcResult.turnInDay.toNumber().should.be.equal(nodeResult.turnInDay);
            bcResult.day.toNumber().should.be.equal(nodeResult.day);
        }
    });

});