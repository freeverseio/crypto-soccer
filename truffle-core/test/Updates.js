const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const timeTravel = require('../utils/TimeTravel.js');
const delegateUtils = require('../utils/delegateCallUtils.js');
const merkleUtils = require('../utils/merkleUtils.js');
const chllUtils = require('../utils/challengeUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Merkle = artifacts.require('Merkle');
const Stakers = artifacts.require("Stakers")



contract('Updates', (accounts) => {
    const nullHash = web3.eth.abi.encodeParameter('bytes32','0x0');
    const nLevelsInOneChallenge = 11;
    const nNonNullLeafsInLeague = 640;
    const nLevelsInLastChallenge = 10; // must be nearest exponent to 640 ... 1024
    
    const it2 = async(text, f) => {};
    
    function normalizeTZ(tz) {
        return 1 + ((tz - 1) % 24);
    }

    const moveToNextVerse = async (updates, extraSecs = 0) => {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.getNextVerseTimestamp().should.be.fulfilled;
        await timeTravel.advanceTime(nextTime - now + extraSecs);
        await timeTravel.advanceBlock().should.be.fulfilled;
    };

    function isCloseEnough(timeResult, timeExpected) {
        // everything is in secs
        allowedError = 4;
        closeEnough  = (timeResult > timeExpected - allowedError); 
        closeEnough = closeEnough && (timeResult < timeExpected + allowedError);
        return closeEnough;
    };
    
    function arrayToHex(x) {
        y = [...x];
        for (i = 0; i < x.length; i++) {
            y[i] = web3.utils.toHex(x[i]);
        }
        return y;
    }
    
    async function asyncForEach(array, callback) {
        for (let index = 0; index < array.length; index++) {
            await callback(array[index], index, array);
        }
    }
    
    async function addTrustedParties(contract, owner, addresses) {
        await asyncForEach(addresses, async (address) => {
            contract.addTrustedParty(address, {from:owner})
        });
    }
    async function enroll(contract, stake, addresses) {
        await asyncForEach(addresses, async (address) => {
            await contract.enroll({from:address, value: stake})
        });
    }
    async function unenroll(contract, addresses) {
        await asyncForEach(addresses, async (address) => {
            await contract.unEnroll({from:address})
        });
    }

    async function deployAndConfigureStakers(owner, parties, updates) {
        stakers  = await Stakers.new({from:owner});
        stakers.setGameOwner(updates.address, {from:owner}).should.be.fulfilled;
        stake = await stakers.kRequiredStake();
        await addTrustedParties(stakers, owner, parties);
        await enroll(stakers, stake, parties);
        await updates.setStakersAddress(stakers.address).should.be.fulfilled;
        return stakers;
    }
    
    beforeEach(async () => {
        depl =  await delegateUtils.deploy(versionNumber = 0, Proxy, '0x0', Assets, Market, Updates);
        proxy  = depl[0];
        assets = depl[1];
        market = depl[2];
        updates = depl[3];
        // // done with delegate calls

        constants = await ConstantsGetters.new().should.be.fulfilled;
        merkle = await Merkle.new().should.be.fulfilled;
        await updates.initUpdates().should.be.fulfilled;
        await updates.setChallengeLevels(nLevelsInOneChallenge, nNonNullLeafsInLeague, nLevelsInLastChallenge).should.be.fulfilled;
        NULL_TIMEZONE = await constants.get_NULL_TIMEZONE().should.be.fulfilled;
        NULL_TIMEZONE = NULL_TIMEZONE.toNumber();
        snapShot = await timeTravel.takeSnapshot();
        snapshotId = snapShot['result'];
        VERSES_PER_DAY = await constants.get_VERSES_PER_DAY().should.be.fulfilled;
        VERSES_PER_ROUND = await constants.get_VERSES_PER_ROUND().should.be.fulfilled;

    });
    
    it('test getAllMatchdaysUTCInRound', async () =>  {
        nextVerseTimestamp = await updates.getNextVerseTimestamp().should.be.fulfilled;
        nextVerseTimestamp = nextVerseTimestamp.toNumber();
        timeZoneForRound1 = await updates.getTimeZoneForRound1().should.be.fulfilled;
        // tests for init timezone
        utc = await updates.getAllMatchdaysUTCInRound(tz = timeZoneForRound1, round = 0).should.be.fulfilled;
        nMatchesPerRound = 14;
        for (matchDay = 0; matchDay < nMatchesPerRound/2; matchDay++) {
            utc[2 * matchDay].toNumber().should.be.equal(nextVerseTimestamp + 24*3600*matchDay);
            utc[1 + 2 * matchDay].toNumber().should.be.equal(nextVerseTimestamp + 19 * 1800 + 24*3600*matchDay);
        }
        // tests for last timezone
        utc = await updates.getAllMatchdaysUTCInRound(tz = (timeZoneForRound1 - 1), round = 0).should.be.fulfilled;
        nMatchesPerRound = 14;
        for (matchDay = 0; matchDay < nMatchesPerRound/2; matchDay++) {
            utc[2 * matchDay].toNumber().should.be.equal(nextVerseTimestamp + 24*3600*matchDay + 23*3600);
            utc[1 + 2 * matchDay].toNumber().should.be.equal(nextVerseTimestamp + 19 * 1800 + 24*3600*matchDay + 23*3600);
        }
        // tests for first timezone, round = 1
        utc = await updates.getAllMatchdaysUTCInRound(tz = timeZoneForRound1, round = 1).should.be.fulfilled;
        nMatchesPerRound = 14;
        for (matchDay = 0; matchDay < nMatchesPerRound/2; matchDay++) {
            utc[2 * matchDay].toNumber().should.be.equal(nextVerseTimestamp + 24*3600*matchDay + 7*24*3600);
            utc[1 + 2 * matchDay].toNumber().should.be.equal(nextVerseTimestamp + 19 * 1800 + 24*3600*matchDay + 7*24*3600);
        }    
    });
    
    it('test getCurrentRoundPure', async () =>  {
        result = await assets.getCurrentRoundPure(tz = 5, tz1 = 5, verse = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await assets.getCurrentRoundPure(tz = 24, tz1 = 5, verse = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await assets.getCurrentRoundPure(tz = 4, tz1 = 5, verse = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        VERSES_DAY = 24*4;
        VERSES_ROUND = 7 * VERSES_DAY;
        // move to start of round 1 for 1st tz:
        result = await assets.getCurrentRoundPure(tz = 5, tz1 = 5, verse = VERSES_ROUND).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await assets.getCurrentRoundPure(tz = 4, tz1 = 5, verse = VERSES_ROUND).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await assets.getCurrentRoundPure(tz = 24, tz1 = 5, verse = VERSES_ROUND).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        // move to start of round 1 for 1st tz after tz1:
        result = await assets.getCurrentRoundPure(tz = 5, tz1 = 5, verse = VERSES_ROUND + 4).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await assets.getCurrentRoundPure(tz = 6, tz1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await assets.getCurrentRoundPure(tz = 7, tz1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await assets.getCurrentRoundPure(tz = 24, tz1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        // move to start of round 1 for last tz to reach it:
        result = await assets.getCurrentRoundPure(tz = 5, tz1 = 5, verse = 2 * VERSES_ROUND - 4).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await assets.getCurrentRoundPure(tz = 4, tz1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await assets.getCurrentRoundPure(tz = 24, tz1 = 5, verse).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
    });

    it('test getMatchUTC', async () =>  {
        nextVerseTimestamp = await updates.getNextVerseTimestamp().should.be.fulfilled;
        nextVerseTimestamp = nextVerseTimestamp.toNumber();
        timeZoneForRound1 = await updates.getTimeZoneForRound1().should.be.fulfilled;
        // tests for init timezone
        utc = await updates.getMatchUTC(tz = timeZoneForRound1, round = 0, matchDay = 0).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp);
        utc = await updates.getMatchUTC(tz = timeZoneForRound1, round = 0, matchDay = 2).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp + 24*3600);
        utc = await updates.getMatchUTC(tz = timeZoneForRound1, round = 0, matchDay = 1).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp + 9.5*3600);
        utc = await updates.getMatchUTC(tz = timeZoneForRound1, round = 1, matchDay = 1).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp + 9.5*3600 + 7*24*3600);
        utc = await updates.getMatchUTC(tz = timeZoneForRound1, round = 1, matchDay = 2).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp + 24*3600 + 7*24*3600);
        // tests for other timezones
        tz = 1;
        deltaN = (tz >= timeZoneForRound1) ? (tz-timeZoneForRound1) : (24+tz-timeZoneForRound1); 
        utc = await updates.getMatchUTC(tz, round = 0, matchDay = 0).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp + deltaN * 3600);
        tz = 24;
        deltaN = (tz >= timeZoneForRound1) ? (tz-timeZoneForRound1) : (24+tz-timeZoneForRound1); 
        utc = await updates.getMatchUTC(tz, round = 0, matchDay = 0).should.be.fulfilled;
        utc.toNumber().should.be.equal(nextVerseTimestamp + deltaN * 3600);
    });


    
    it('test that cannot initialize updates twice', async () =>  {
        await updates.initUpdates().should.be.rejected;
    });
    
    it('check timezones for this verse', async () =>  {
        TZForRound1 = 2;
        result = "";
        for (verse = 0; verse < 10*VERSES_PER_DAY.toNumber(); verse += 13) {
            var {0: tz, 1: matchday, 2: turn} = await updates._timeZoneToUpdatePure(verse, TZForRound1).should.be.fulfilled;
            day = Math.floor(0.25 * verse / 24);
            thisResult = " | verse = " + verse + 
                ", tz = " + tz.toNumber() + 
                ", matchday = " + matchday.toNumber() +
                ", turn = " + turn.toNumber();
            result += thisResult;
        }
        expected = " | verse = 0, tz = 2, matchday = 0, turn = 0 | verse = 13, tz = 5, matchday = 0, turn = 1 | verse = 26, tz = 0, matchday = 0, turn = 0 | verse = 39, tz = 2, matchday = 1, turn = 1 | verse = 52, tz = 15, matchday = 0, turn = 0 | verse = 65, tz = 18, matchday = 0, turn = 1 | verse = 78, tz = 12, matchday = 1, turn = 0 | verse = 91, tz = 15, matchday = 1, turn = 1 | verse = 104, tz = 4, matchday = 2, turn = 0 | verse = 117, tz = 7, matchday = 2, turn = 1 | verse = 130, tz = 1, matchday = 1, turn = 0 | verse = 143, tz = 4, matchday = 3, turn = 1 | verse = 156, tz = 17, matchday = 2, turn = 0 | verse = 169, tz = 20, matchday = 2, turn = 1 | verse = 182, tz = 14, matchday = 3, turn = 0 | verse = 195, tz = 17, matchday = 3, turn = 1 | verse = 208, tz = 6, matchday = 4, turn = 0 | verse = 221, tz = 9, matchday = 4, turn = 1 | verse = 234, tz = 3, matchday = 5, turn = 0 | verse = 247, tz = 6, matchday = 5, turn = 1 | verse = 260, tz = 19, matchday = 4, turn = 0 | verse = 273, tz = 22, matchday = 4, turn = 1 | verse = 286, tz = 16, matchday = 5, turn = 0 | verse = 299, tz = 19, matchday = 5, turn = 1 | verse = 312, tz = 8, matchday = 6, turn = 0 | verse = 325, tz = 11, matchday = 6, turn = 1 | verse = 338, tz = 5, matchday = 7, turn = 0 | verse = 351, tz = 8, matchday = 7, turn = 1 | verse = 364, tz = 21, matchday = 6, turn = 0 | verse = 377, tz = 24, matchday = 6, turn = 1 | verse = 390, tz = 18, matchday = 7, turn = 0 | verse = 403, tz = 21, matchday = 7, turn = 1 | verse = 416, tz = 10, matchday = 8, turn = 0 | verse = 429, tz = 13, matchday = 8, turn = 1 | verse = 442, tz = 7, matchday = 9, turn = 0 | verse = 455, tz = 10, matchday = 9, turn = 1 | verse = 468, tz = 23, matchday = 8, turn = 0 | verse = 481, tz = 2, matchday = 10, turn = 1 | verse = 494, tz = 20, matchday = 9, turn = 0 | verse = 507, tz = 23, matchday = 9, turn = 1 | verse = 520, tz = 12, matchday = 10, turn = 0 | verse = 533, tz = 15, matchday = 10, turn = 1 | verse = 546, tz = 9, matchday = 11, turn = 0 | verse = 559, tz = 12, matchday = 11, turn = 1 | verse = 572, tz = 1, matchday = 10, turn = 0 | verse = 585, tz = 4, matchday = 12, turn = 1 | verse = 598, tz = 22, matchday = 11, turn = 0 | verse = 611, tz = 1, matchday = 11, turn = 1 | verse = 624, tz = 14, matchday = 12, turn = 0 | verse = 637, tz = 17, matchday = 12, turn = 1 | verse = 650, tz = 11, matchday = 13, turn = 0 | verse = 663, tz = 14, matchday = 13, turn = 1 | verse = 676, tz = 3, matchday = 0, turn = 0 | verse = 689, tz = 6, matchday = 0, turn = 1 | verse = 702, tz = 24, matchday = 13, turn = 0 | verse = 715, tz = 3, matchday = 1, turn = 1 | verse = 728, tz = 16, matchday = 0, turn = 0 | verse = 741, tz = 19, matchday = 0, turn = 1 | verse = 754, tz = 13, matchday = 1, turn = 0 | verse = 767, tz = 16, matchday = 1, turn = 1 | verse = 780, tz = 5, matchday = 2, turn = 0 | verse = 793, tz = 8, matchday = 2, turn = 1 | verse = 806, tz = 2, matchday = 3, turn = 0 | verse = 819, tz = 5, matchday = 3, turn = 1 | verse = 832, tz = 18, matchday = 2, turn = 0 | verse = 845, tz = 21, matchday = 2, turn = 1 | verse = 858, tz = 15, matchday = 3, turn = 0 | verse = 871, tz = 18, matchday = 3, turn = 1 | verse = 884, tz = 7, matchday = 4, turn = 0 | verse = 897, tz = 10, matchday = 4, turn = 1 | verse = 910, tz = 4, matchday = 5, turn = 0 | verse = 923, tz = 7, matchday = 5, turn = 1 | verse = 936, tz = 20, matchday = 4, turn = 0 | verse = 949, tz = 23, matchday = 4, turn = 1";
        result.should.be.equal(expected);
    });
    
    it('require that BC and local time are less than 15 sec out of sync', async () =>  {
        blockChainTimeSec = await updates.getNow().should.be.fulfilled;
        localTimeMs = Date.now();
        // the substraction is in miliseconds:
        // require less than 3 hours
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 3*3600*1000).should.be.equal(true);
        // require less than 1 hour
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 1*3600*1000).should.be.equal(true);
        // require less than 30 min
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 30*60*1000).should.be.equal(true);
        // require less than 10 min
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 10*60*1000).should.be.equal(true);
        // require less than 5 min
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 5*60*1000).should.be.equal(true);
        // require less than 1 min
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 60*1000).should.be.equal(true);
        // require less than 20 sec
        (Math.abs(blockChainTimeSec.toNumber()*1000 - localTimeMs) < 20*1000).should.be.equal(true);
    });
    
    it('check BC is set up in agreement with the local time', async () =>  {
        nextVerseTimestamp = await updates.getNextVerseTimestamp().should.be.fulfilled;
        timeZoneForRound1 = await updates.getTimeZoneForRound1().should.be.fulfilled;
        localTimeMs = Date.now();
        now = new Date(localTimeMs);
        nextVerse = new Date(nextVerseTimestamp.toNumber() * 1000);
        if (now.getUTCMinutes() < 27) {
            expectedHour = now.getUTCHours();
        } else {
            expectedHour = now.getUTCHours() + 1;
        }
        nextVerse.getUTCFullYear().should.be.equal(now.getUTCFullYear());
        nextVerse.getUTCMonth().should.be.equal(now.getUTCMonth());
        nextVerse.getUTCDate().should.be.equal(now.getUTCDate());
        nextVerse.getUTCHours().should.be.equal(expectedHour);
        nextVerse.getUTCMinutes().should.be.equal(30);
        nextVerse.getUTCSeconds().should.be.equal(0);
        timeZoneForRound1.toNumber().should.be.equal(expectedHour);
    });
    
    it('wait some minutes', async () =>  {
        now = await updates.getNow().should.be.fulfilled;
        block = await web3.eth.getBlockNumber().should.be.fulfilled;
        extraTime = 3*60
        await timeTravel.advanceTime(extraTime).should.be.fulfilled;
        await timeTravel.advanceBlock().should.be.fulfilled;
        newNow = await updates.getNow().should.be.fulfilled;
        newBlock = await web3.eth.getBlockNumber().should.be.fulfilled;
        newBlock.should.be.equal(block+1);
        await isCloseEnough(newNow.toNumber(), now.toNumber() + extraTime).should.be.equal(true);
        await timeTravel.revertToSnapShot(snapshotId);
        newNow = await updates.getNow().should.be.fulfilled;
        isCloseEnough(newNow.toNumber(), now.toNumber()).should.be.equal(true)
    });
    
    it('submitActions to timezone', async () =>  {
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        verseBefore = await updates.getCurrentVerse().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        await moveToNextVerse(updates, extraTime = -10)        
        await timeTravel.advanceTime(20);
        await timeTravel.advanceBlock().should.be.fulfilled;
        const cif = "ciao";
        tx = await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboys"), nullHash, nullHash, 2, cif).should.be.fulfilled;
        timeZoneToUpdate = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        verse = await updates.getCurrentVerse().should.be.fulfilled;
        verse.toNumber().should.be.equal(verseBefore.toNumber() + 1); 
        timeZoneToUpdate[0].toNumber().should.be.equal(timeZoneToUpdateBefore[0].toNumber()); // tz to update does not change during the first 4 verses
        seed1 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        seed1.should.not.be.equal(seed0);
        now = await updates.getNow().should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "ActionsSubmission", (event) => {
            return event.seed == seed1 && isCloseEnough(event.submissionTime.toNumber(), now.toNumber());
        });
    });

    it('update Timezone once', async () =>  {
        const [owner, gameAddr, alice, bob, carol, dave, erin, frank] = accounts;
        parties = [alice, bob, carol, dave, erin, frank]
        stakes = await deployAndConfigureStakers(owner, parties, updates);

        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        await moveToNextVerse(updates, extraSecs = -10);
        await timeTravel.advanceTime(20);
        const cif = "ciao2";
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
        timeZoneToUpdate = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        now = await updates.getNow().should.be.fulfilled;
        await updates.updateTZ(root =  web3.utils.keccak256("hiboyz"), {from:erin}).should.be.fulfilled;
        submissionTime = await updates.getLastActionsSubmissionTime(timeZoneToUpdateBefore[0].toNumber()).should.be.fulfilled;
        timeZoneToUpdateAfter = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        isCloseEnough(timeZoneToUpdate[0].toNumber(), timeZoneToUpdateBefore[0].toNumber()).should.be.equal(true)
        isCloseEnough(submissionTime.toNumber(), now.toNumber()).should.be.equal(true)
    });

    it('moveToNextVerse', async () =>  {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.getNextVerseTimestamp().should.be.fulfilled;
        (nextTime - now > 0).should.be.equal(true)
        await moveToNextVerse(updates, extraSecs = 0);
        now = await updates.getNow().should.be.fulfilled;
        (nextTime - now > 0).should.be.equal(false)
        
    });

    it('update Timezone many times', async () =>  {
        result = await assets.getCurrentRound(tz = 1).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await assets.getCurrentRound(tz = 24).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        console.log("warning: the next test lasts about 20 secs...")
        await moveToNextVerse(updates, extraSecs = 10);
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        const cif = "ciao3";
        for (verse = 0; verse < 110; verse++) {
            await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
            await moveToNextVerse(updates, extraSecs = 10);
        }
    });

    it('timeZoneToUpdateBefore only increases turnInDay by one after submiteActionsRoot', async () =>  {
        await moveToNextVerse(updates, extraSecs = 2);
        var {0: tzBefore, 1: dayBefore, 2: turnInDayBefore} = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        const cif = "ciao3";
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
        var {0: tzAfter, 1: dayAfter, 2: turnInDayAfter} = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        tzAfter.toNumber().should.be.equal(tzBefore.toNumber());
        dayAfter.toNumber().should.be.equal(dayBefore.toNumber());
        (turnInDayAfter.toNumber() - turnInDayBefore.toNumber()).should.be.equal(1);
    });
    
    // level 0: root
    // level 1: 2048 league Roots
    // level 2: 640 leafs for each
    
    it('challenging a tz', async () =>  {
        const [owner, gameAddr, alice, bob, carol, dave, erin, frank] = accounts;
        parties = [alice, bob, carol, dave, erin, frank]
        stakes = await deployAndConfigureStakers(owner, parties, updates);

        // level 0 can only challenge leaf 0, as there is only 1 root
        challengePos = [0];
        var level = 0;

        // move to next verse adn submit actions
        await moveToNextVerse(updates, extraSecs = 2);
        var {0: tz} = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        const cif = "ciao3";
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
        tzZeroBased = 2;

        // create leafs by building them from an orgmap:
        const {0: orgMapHeader, 1: orgMap, 2: userActions} = await chllUtils.createOrgMap(assets, nCountriesPerTZ = 2, nActiveUsersPerCountry = 6)
        const {0: leafsADecimal, 1: nLeaguesInTzA} = chllUtils.createLeafsForOrgMap(day = 3, half = 0, orgMapHeader[tzZeroBased], nNonNullLeafsInLeague);
        const {0: leafsBDecimal, 1: nLeaguesInTzB} = chllUtils.createLeafsForOrgMap(day = 13, half = 1, orgMapHeader[tzZeroBased], nNonNullLeafsInLeague);
        leafsA = chllUtils.leafsToBytes32(leafsADecimal);
        leafsB = chllUtils.leafsToBytes32(leafsBDecimal);

        // set the levelVerifiableByBC to adjust to as many leagues as you have
        nLeafsPerRoot = 2**nLevelsInOneChallenge;
        levelVerifiableByBC = merkleUtils.computeLevelVerifiableByBC(nLeaguesInTzA, nLeafsPerRoot);
        await updates.setLevelVerifiableByBC(levelVerifiableByBC).should.be.fulfilled;

        // build merkle structs for 2 different days
        merkleStructA = merkleUtils.buildMerkleStruct(leafsA, nLeafsPerRoot, levelVerifiableByBC);
        merkleStructB = merkleUtils.buildMerkleStruct(leafsB, nLeafsPerRoot, levelVerifiableByBC);
        
        // get data to challenge at level 0 (level is inferred from the length of challengePos).
        var {0: challValA, 1: proofA, 2: roots2SubmitA} = merkleUtils.getDataToChallenge(challengePos, leafsA, merkleStructA, nLeafsPerRoot, levelVerifiableByBC);
        var {0: challValB, 1: proofB, 2: roots2SubmitB} = merkleUtils.getDataToChallenge(challengePos, leafsB, merkleStructB, nLeafsPerRoot, levelVerifiableByBC);

        // First challenge fails because the TZ has not been updated yet with a root
        await updates.challengeTZ(challVal = nullHash, challengePos[level], proof = [], roots2SubmitA, {from:alice}).should.be.rejected;

        // So let's update with rootA...
        await updates.updateTZ(root = merkleStructA[lev = 0][pos = 0], {from:alice}).should.be.fulfilled;

        // We can not challenge with something compatible with rootA:
        await updates.challengeTZ(challVal = nullHash, challengePos[level], proof = [], roots2SubmitA, {from:bob}).should.be.rejected;

        // We can not challenge with by alice again:
        await updates.challengeTZ(challVal = nullHash, challengePos[level], proof = [], roots2SubmitA, {from:alice}).should.be.rejected;

        // ...but we can challenge with rootsB, that differ from rootsA:
        assert.notEqual(merkleStructA[lev = 0][pos = 0], merkleStructB[lev = 0][pos = 0], "wrong leafsA should lead to different root");
        await updates.challengeTZ(challVal = nullHash, challengePos[level], proof = [], roots2SubmitB, {from:bob}).should.be.fulfilled;

        // check that level increased:
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        isSet.should.be.equal(false);
        level = lev.toNumber();
        
        // TODO: test that vals are gotten from events
        // Challenge one of the leagues:
        challengePos.push(newChallengePos = 1);
        var {0: challValA, 1: proofA, 2: roots2SubmitA} = merkleUtils.getDataToChallenge(challengePos, leafsA, merkleStructA, nLeafsPerRoot, levelVerifiableByBC);
        var {0: challValB, 1: proofB, 2: roots2SubmitB} = merkleUtils.getDataToChallenge(challengePos, leafsB, merkleStructB, nLeafsPerRoot, levelVerifiableByBC);

        assert.equal(merkleUtils.merkleRoot(roots2SubmitB, nLevelsInLastChallenge), merkleStructB[1][1], "wrong selection of roots2submit");
        assert.equal(merkleUtils.merkleRoot(roots2SubmitB, nLevelsInLastChallenge), challValB, "wrong selection of roots2submit");
        assert.equal(
            merkleUtils.merkleRoot(roots2SubmitB, nLevelsInLastChallenge), 
            await merkle.merkleRoot(roots2SubmitB, nLevelsInLastChallenge), 
            "nonmatching merkleRoots"
        );
        assert.equal(merkleUtils.merkleRoot(roots2SubmitA, nLevelsInLastChallenge), merkleStructA[1][1], "wrong selection of roots2submit");
        assert.equal(merkleUtils.merkleRoot(roots2SubmitA, nLevelsInLastChallenge), challValA, "wrong selection of roots2submit");
        
        // As always, first check that we cannot submit roots that coinicide with previous:
        await updates.challengeTZ(challValB, challengePos[level], proofB, roots2SubmitB, {from:carol}).should.be.rejected;
        
        // But we can with differing ones:
        await updates.challengeTZ(challValB, challengePos[level], proofB, roots2SubmitA, {from:carol}).should.be.fulfilled;

        // Check that we move to level 2
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(2);
        challValB_backup = challValB;
        newChallengePos_backup = newChallengePos;
        proofB_backup = [...proofB];
        roots2SubmitA_backup = [...roots2SubmitA];
        
        // finally, the last challenge, is one that the BC can check
        // we will to a challenge of level 3 that will instantaneously resolve into killing the level2 and reverting to level1
        // try with wrong leaves:
        await updates.BCVerifableChallengeFake([...roots2SubmitB], forceSuccess = true, {from: dave}).should.be.rejected;
        // but I can submit different ones. In this case the BC decides according to forceSuccess
        await updates.BCVerifableChallengeFake([...roots2SubmitA], forceSuccess = false, {from: dave}).should.be.rejected;
        await updates.BCVerifableChallengeFake([...roots2SubmitA], forceSuccess = true, {from: dave}).should.be.fulfilled;
        
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        isSet.should.be.equal(false);
        level = lev.toNumber();
        
        // challenge again to move to level2, and now we will wait time
        await updates.challengeTZ(challValB_backup, newChallengePos_backup, proofB_backup, roots2SubmitA_backup, {from: erin}).should.be.fulfilled;
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(2);
        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(2);
        isSet.should.be.equal(false);
        level = lev.toNumber();
        
        challengeTime = await constants.get_CHALLENGE_TIME().should.be.fulfilled;
        await timeTravel.advanceTime(challengeTime.toNumber() + 10).should.be.fulfilled;
        await timeTravel.advanceBlock().should.be.fulfilled;

        // note that getStatus realises that we moved to level 0, but not the written stuff
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(2);
        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(0);
        nJumps.toNumber().should.be.equal(1);
        isSet.should.be.equal(false);
        level = lev.toNumber();

        // I should not be able to provide a new update, nor new actionRoots, for 2 reasons:
        //      we're not in the next verse yet
        //      the previous verse is not settled yet
        // In this case, it fails because of the first reason. TODO: add test for 2nd.
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.rejected;
        await updates.setLevelVerifiableByBC(3).should.be.fulfilled;

        await updates.updateTZ(root = merkleStructA[lev = 0][pos = 0], {from: erin}).should.be.rejected;
        
        await timeTravel.advanceTime(challengeTime.toNumber() + 10).should.be.fulfilled;
        await timeTravel.advanceBlock().should.be.fulfilled;
        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(0);
        isSet.should.be.equal(true);
    });
    
    
    
    // it('(takes a long time!) challenging a tz beyond the next timezone!', async () =>  {
    //     await moveToNextVerse(updates, extraSecs = 2);
    //     var {0: tz} = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
    //     const cif = "ciao3";
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
    //     await updates.setLevelVerifiableByBC(3).should.be.fulfilled;

    //     nLeafsPerRoot = 16;
    //     nChallenges = 3;
    //     nTotalLeafs = nLeafsPerRoot**3;
    //     nTotalLevels = Math.log2(nTotalLeafs);
    //     nLevelsPerRoot = Math.log2(nLeafsPerRoot);
    //     leafsA = Array.from(new Array(nTotalLeafs), (x,i) => web3.utils.keccak256(i.toString()));
    //     leafsB = Array.from(new Array(nTotalLeafs), (x,i) => web3.utils.keccak256((i+1).toString()));
    //     merkleStructA = merkleUtils.buildMerkleStruct(leafsA, nLeafsPerRoot);
    //     merkleStructB = merkleUtils.buildMerkleStruct(leafsB, nLeafsPerRoot);
    //     // We update with the correct root...
    //     await updates.updateTZ(root = merkleUtils.merkleRoot(leafsA, nTotalLevels)).should.be.fulfilled;

    //     secsBetweenVerses = await constants.get_SECS_BETWEEN_VERSES().should.be.fulfilled;
    //     challengeTime = await constants.get_CHALLENGE_TIME().should.be.fulfilled;
    //     nInterations = Math.floor(secsBetweenVerses.toNumber()/challengeTime.toNumber())

    //     // prepare for challenges to level 1 -> level 2
    //     newChallengePos = 7;
    //     challengePos = [];
    //     challengePos.push(newChallengePos);
    //     var {0: challValA, 1: proofA, 2: roots2SubmitA} = merkleUtils.getDataToChallenge(challengePos, merkleStructA, nLeafsPerRoot);
    //     var {0: challValB, 1: proofB, 2: roots2SubmitB} = merkleUtils.getDataToChallenge(challengePos, merkleStructB, nLeafsPerRoot);
    
    //     forceSuccess = false;
    //     for (iter = 0; iter < nInterations - 1; iter++) {
    //         var {0: level, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
    //         level.toNumber().should.be.equal(0);
    //         isSet.should.be.equal(false);
    //         // challenge 0 -> 1
    //         await updates.challengeTZ(challVal = nullHash, challengePos = 0, proof = [], merkleStructB[1]).should.be.fulfilled;
    //         // challenge 1 -> 2
    //         await updates.challengeTZ(challValB, newChallengePos, proofB, roots2SubmitA).should.be.fulfilled;
    //         var {0: level, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
    //         level.toNumber().should.be.equal(2);
    //         nJumps.toNumber().should.be.equal(0);
    //         isSet.should.be.equal(false);

    //         // wait so that this challenge of level 2 is successful
    //         await timeTravel.advanceTime(challengeTime.toNumber() + 2).should.be.fulfilled;
    //         await timeTravel.advanceBlock().should.be.fulfilled;

    //         var {0: level, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
    //         level.toNumber().should.be.equal(0);
    //         nJumps.toNumber().should.be.equal(1);
    //         isSet.should.be.equal(false);
    //     }
    //     // Time-wise, we are ready for next TZ actions root submission, but extraordinarily,
    //     // the previous timezone is not settled yet:
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.rejected;
    //     // so just wait one CHLL period extra.
    //     await timeTravel.advanceTime(challengeTime.toNumber() + 10).should.be.fulfilled;
    //     await timeTravel.advanceBlock().should.be.fulfilled;
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
    //     await updates.setLevelVerifiableByBC(3).should.be.fulfilled;
    //     await updates.updateTZ(root = merkleUtils.merkleRoot(leafsA, nTotalLevels)).should.be.fulfilled;
    // });
    
    
    // it('(takes a long time!) challenging a tz beyond the next timezone! -- almost', async () =>  {
    //     // identical to previous test but we wait 1 challenge time less!
    //     // so at the very end, we're not allowed to submit actions because the time has not come for next timezone
    //     await moveToNextVerse(updates, extraSecs = 2);
    //     var {0: tz} = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
    //     const cif = "ciao3";
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
    //     await updates.setLevelVerifiableByBC(3).should.be.fulfilled;

    //     nLeafsPerRoot = 16;
    //     nChallenges = 3;
    //     nTotalLeafs = nLeafsPerRoot**3;
    //     nTotalLevels = Math.log2(nTotalLeafs);
    //     nLevelsPerRoot = Math.log2(nLeafsPerRoot);
    //     leafsA = Array.from(new Array(nTotalLeafs), (x,i) => web3.utils.keccak256(i.toString()));
    //     leafsB = Array.from(new Array(nTotalLeafs), (x,i) => web3.utils.keccak256((i+1).toString()));
    //     merkleStructA = merkleUtils.buildMerkleStruct(leafsA, nLeafsPerRoot);
    //     merkleStructB = merkleUtils.buildMerkleStruct(leafsB, nLeafsPerRoot);
    //     // We update with the correct root...
    //     await updates.updateTZ(root = merkleUtils.merkleRoot(leafsA, nTotalLevels)).should.be.fulfilled;

    //     secsBetweenVerses = await constants.get_SECS_BETWEEN_VERSES().should.be.fulfilled;
    //     challengeTime = await constants.get_CHALLENGE_TIME().should.be.fulfilled;
    //     nInterations = Math.floor(secsBetweenVerses.toNumber()/challengeTime.toNumber())

    //     // prepare for challenges to level 1 -> level 2
    //     newChallengePos = 7;
    //     challengePos = [];
    //     challengePos.push(newChallengePos);
    //     var {0: challValA, 1: proofA, 2: roots2SubmitA} = merkleUtils.getDataToChallenge(challengePos, merkleStructA, nLeafsPerRoot);
    //     var {0: challValB, 1: proofB, 2: roots2SubmitB} = merkleUtils.getDataToChallenge(challengePos, merkleStructB, nLeafsPerRoot);
    
    //     forceSuccess = false;
    //     for (iter = 0; iter < nInterations - 2; iter++) {
    //         var {0: level, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
    //         level.toNumber().should.be.equal(0);
    //         isSet.should.be.equal(false);
    //         // challenge 0 -> 1
    //         await updates.challengeTZ(challVal = nullHash, challengePos = 0, proof = [], merkleStructB[1]).should.be.fulfilled;
    //         // challenge 1 -> 2
    //         await updates.challengeTZ(challValB, newChallengePos, proofB, roots2SubmitA).should.be.fulfilled;
    //         var {0: level, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
    //         level.toNumber().should.be.equal(2);
    //         nJumps.toNumber().should.be.equal(0);
    //         isSet.should.be.equal(false);

    //         // wait so that this challenge of level 2 is successful
    //         await timeTravel.advanceTime(challengeTime.toNumber() + 2).should.be.fulfilled;
    //         await timeTravel.advanceBlock().should.be.fulfilled;

    //         var {0: level, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
    //         level.toNumber().should.be.equal(0);
    //         nJumps.toNumber().should.be.equal(1);
    //         isSet.should.be.equal(false);
    //     }
    //     // Time-wise, we are ready for next TZ actions root submission, but extraordinarily,
    //     // the previous timezone is not settled yet:
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.rejected;
    //     // so just wait one CHLL period extra.
    //     await timeTravel.advanceTime(challengeTime.toNumber() + 10).should.be.fulfilled;
    //     await timeTravel.advanceBlock().should.be.fulfilled;
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.rejected;
    //     await updates.setLevelVerifiableByBC(3).should.be.fulfilled;
    //     await updates.updateTZ(root = merkleUtils.merkleRoot(leafsA, nTotalLevels)).should.be.rejected;
    // });

    
    it('true status of timezone challenge', async () =>  {
        challengeTime = await constants.get_CHALLENGE_TIME().should.be.fulfilled;
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(0.5*challengeTime), lastUpdate = 0, writtenLevel = 0).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(1.5*challengeTime), lastUpdate = 0, writtenLevel = 0).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(1.5*challengeTime), lastUpdate = 0, writtenLevel = 1).should.be.fulfilled;
        level.toNumber().should.be.equal(1);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(3.5*challengeTime), lastUpdate = 0, writtenLevel = 1).should.be.fulfilled;
        level.toNumber().should.be.equal(1);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(0.5*challengeTime), lastUpdate = 0, writtenLevel = 2).should.be.fulfilled;
        level.toNumber().should.be.equal(2);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(1.5*challengeTime), lastUpdate = 0, writtenLevel = 2).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(2.5*challengeTime), lastUpdate = 0, writtenLevel = 2).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(3.5*challengeTime), lastUpdate = 0, writtenLevel = 2).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(0.5*challengeTime), lastUpdate = 0, writtenLevel = 3).should.be.fulfilled;
        level.toNumber().should.be.equal(3);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(1.5*challengeTime), lastUpdate = 0, writtenLevel = 3).should.be.fulfilled;
        level.toNumber().should.be.equal(1);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(2.5*challengeTime), lastUpdate = 0, writtenLevel = 3).should.be.fulfilled;
        level.toNumber().should.be.equal(1);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(3.5*challengeTime), lastUpdate = 0, writtenLevel = 3).should.be.fulfilled;
        level.toNumber().should.be.equal(1);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(0.5*challengeTime), lastUpdate = 0, writtenLevel = 4).should.be.fulfilled;
        level.toNumber().should.be.equal(4);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(1.5*challengeTime), lastUpdate = 0, writtenLevel = 4).should.be.fulfilled;
        level.toNumber().should.be.equal(2);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(2.5*challengeTime), lastUpdate = 0, writtenLevel = 4).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(false);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(3.5*challengeTime), lastUpdate = 0, writtenLevel = 4).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(true);
        var {0: level, 1: nJumps, 2: isSet} = await updates.getStatusPure(nowTime = Math.floor(4.5*challengeTime), lastUpdate = 0, writtenLevel = 4).should.be.fulfilled;
        level.toNumber().should.be.equal(0);
        isSet.should.be.equal(true);
    });
    
    
    // A = correct day and half
    // B = correct day, incorrect half
    // C = incorrect day and incorrect half
    // 0: A, 1: B, 2: A => so the leafs provided by A are the correct ones and everyone fails to challenge A.
    it('vefiable challenge', async () =>  {
        const [owner, gameAddr, alice, bob, carol, dave, erin, frank] = accounts;
        parties = [alice, bob, carol, dave, erin, frank]
        stakes = await deployAndConfigureStakers(owner, parties, updates);

        // level 0 can only challenge leaf 0, as there is only 1 root
        challengePos = [0];
        var level = 0;

        // move to next verse adn submit actions
        await moveToNextVerse(updates, extraSecs = 2);
        var {0: tz,  1: day, 2: half} = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        tz      = tz.toNumber();
        day     = day.toNumber();
        half    = half.toNumber();
        differentDay = (day == 7) ? 8 : 7;
        const cif = "ciao3";
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), nullHash, nullHash, 2, cif).should.be.fulfilled;
        tzZeroBased = tz-1;
        // create leafs by building them from an orgmap:
        const {0: orgMapHeader, 1: orgMap, 2: userActions} = await chllUtils.createOrgMap(assets, nCountriesPerTZ = 2, nActiveUsersPerCountry = 6)
        const {0: leafsADecimal, 1: nLeaguesInTzA} = chllUtils.createLeafsForOrgMap(day, half, orgMapHeader[tzZeroBased], nNonNullLeafsInLeague);
        const {0: leafsBDecimal, 1: nLeaguesInTzB} = chllUtils.createLeafsForOrgMap(day, 1 - half, orgMapHeader[tzZeroBased], nNonNullLeafsInLeague);
        const {0: leafsCDecimal, 1: nLeaguesInTzC} = chllUtils.createLeafsForOrgMap(differentDay, 1 - half, orgMapHeader[tzZeroBased], nNonNullLeafsInLeague);
        leafsA = chllUtils.leafsToBytes32(leafsADecimal);
        leafsB = chllUtils.leafsToBytes32(leafsBDecimal);
        leafsC = chllUtils.leafsToBytes32(leafsCDecimal);

        // set the levelVerifiableByBC to adjust to as many leagues as you have
        nLeafsPerRoot = 2**nLevelsInOneChallenge;
        levelVerifiableByBC = merkleUtils.computeLevelVerifiableByBC(nLeaguesInTzA, nLeafsPerRoot);
        await updates.setLevelVerifiableByBC(levelVerifiableByBC).should.be.fulfilled;

        // build merkle structs for 2 different days
        merkleStructA = merkleUtils.buildMerkleStruct(leafsA, nLeafsPerRoot, levelVerifiableByBC);
        merkleStructB = merkleUtils.buildMerkleStruct(leafsB, nLeafsPerRoot, levelVerifiableByBC);
        merkleStructC = merkleUtils.buildMerkleStruct(leafsC, nLeafsPerRoot, levelVerifiableByBC);
        
        // get data to challenge at level 0 (level is inferred from the length of challengePos).
        var {0: challValA, 1: proofA, 2: roots2SubmitA} = merkleUtils.getDataToChallenge(challengePos, leafsA, merkleStructA, nLeafsPerRoot, levelVerifiableByBC);
        var {0: challValB, 1: proofB, 2: roots2SubmitB} = merkleUtils.getDataToChallenge(challengePos, leafsB, merkleStructB, nLeafsPerRoot, levelVerifiableByBC);
        var {0: challValC, 1: proofC, 2: roots2SubmitC} = merkleUtils.getDataToChallenge(challengePos, leafsC, merkleStructC, nLeafsPerRoot, levelVerifiableByBC);

        // Level0: A
        await updates.updateTZ(root = merkleStructA[lev = 0][pos = 0], {from:alice}).should.be.fulfilled;

        // Level1: B
        await updates.challengeTZ(challVal = nullHash, challengePos[level], proof = [], roots2SubmitB, {from:bob}).should.be.fulfilled;

        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        level = lev.toNumber();
        
        // Level2: C
        challengePos.push(newChallengePos = 1);
        var {0: challValA, 1: proofA, 2: roots2SubmitA} = merkleUtils.getDataToChallenge(challengePos, leafsA, merkleStructA, nLeafsPerRoot, levelVerifiableByBC);
        var {0: challValB, 1: proofB, 2: roots2SubmitB} = merkleUtils.getDataToChallenge(challengePos, leafsB, merkleStructB, nLeafsPerRoot, levelVerifiableByBC);
        var {0: challValC, 1: proofC, 2: roots2SubmitC} = merkleUtils.getDataToChallenge(challengePos, leafsC, merkleStructC, nLeafsPerRoot, levelVerifiableByBC);

        await updates.challengeTZ(challValB, challengePos[level], proofB, roots2SubmitC, {from:carol}).should.be.fulfilled;

        // Check that we move to level 2
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(2);
        
        // finally, the last challenge, is one that the BC can check
        // must provide the same leafs as the last person (C)
        await updates.BCVerifableChallengeZeros([...roots2SubmitA], {from: erin}).should.be.rejected;
        await updates.BCVerifableChallengeZeros([...roots2SubmitB], {from: erin}).should.be.rejected;

        // we succed to prove that C was wrong:
        await updates.BCVerifableChallengeZeros([...roots2SubmitC], {from: erin}).should.be.fulfilled;

        // we go back to level 1
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        var {0: lev, 1: nJumps, 2: isSet} = await updates.getStatus(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(1);
        isSet.should.be.equal(false);
        level = lev.toNumber();

        // Level2: B
        await updates.challengeTZ(challValB, challengePos[level], proofB, roots2SubmitA, {from:dave}).should.be.fulfilled;

        // Check that we move to level 2
        var {0: idx, 1: lev, 2: maxLev} = await updates.getChallengeData(tz, current = true).should.be.fulfilled; 
        lev.toNumber().should.be.equal(2);
        
        // finally, the last challenge, is one that the BC can check
        // must provide the same leafs as the last person (C)
        await updates.BCVerifableChallengeZeros([...roots2SubmitB]).should.be.rejected;
        await updates.BCVerifableChallengeZeros([...roots2SubmitC]).should.be.rejected;

        // we fail to succed to prove that A was wrong with zeros:
        assert.equal(
            chllUtils.areThereUnexpectedZeros([...roots2SubmitA], day, half, nNonNullLeafsInLeague),
            false,
            "unexpected"
        );
        result = await updates.areThereUnexpectedZeros([...roots2SubmitA], day, half).should.be.fulfilled;
        result.should.be.equal(false);
        await updates.BCVerifableChallengeZeros([...roots2SubmitA]).should.be.rejected;
    });
    
    
    

});