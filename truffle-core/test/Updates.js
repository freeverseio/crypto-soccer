const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const timeTravel = require('../utils/TimeTravel.js');
const delegateUtils = require('../utils/delegateCallUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');



contract('Updates', (accounts) => {
    
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
    
    beforeEach(async () => {
        constants = await ConstantsGetters.new().should.be.fulfilled;
        proxy = await Proxy.new().should.be.fulfilled;
        depl = await delegateUtils.deployDelegate(proxy, Assets, Market, Updates);
        updates = depl[2];
        // // done with delegate calls
        await updates.initUpdates().should.be.fulfilled;
        NULL_TIMEZONE = await constants.get_NULL_TIMEZONE().should.be.fulfilled;
        NULL_TIMEZONE = NULL_TIMEZONE.toNumber();
        snapShot = await timeTravel.takeSnapshot();
        snapshotId = snapShot['result'];
        VERSES_PER_DAY = await constants.get_VERSES_PER_DAY().should.be.fulfilled;
        VERSES_PER_ROUND = await constants.get_VERSES_PER_ROUND().should.be.fulfilled;

    });

    it2('test that cannot initialize updates twice', async () =>  {
        await updates.initUpdates().should.be.rejected;
    });
    
    it('check timezones for this verse', async () =>  {
        TZForRound1 = 2;
        result = "";
        for (verse = 0; verse < 3*VERSES_PER_DAY.toNumber(); verse += 1) {
            var {0: tz, 1: matchday, 2: turn} = await updates._timeZoneToUpdatePure(verse, TZForRound1).should.be.fulfilled;
            day = Math.floor(0.25 * verse / 24);
            thisResult = " verse = " + verse + 
                // ", day = " + day +
                // ", hour = " + (TZForRound1+ 0.5 + 0.25 * verse - day * 24) % 24 +
                ", tz = " + tz.toNumber() + 
                ", matchday = " + matchday.toNumber() +
                ", turn = " + turn.toNumber();
            // console.log(thisResult);
            result += thisResult;
            console.log(thisResult)
        }
        console.log(result)
        // expected = " | verse = 0, day = 0, hour = 2.5, tz = 2, matchday = 0, turn = 0 | verse = 13, day = 0, hour = 5.75, tz = 5, matchday = 0, turn = 1 | verse = 26, day = 0, hour = 9, tz = 0, matchday = 0, turn = 0 | verse = 39, day = 0, hour = 12.25, tz = 2, matchday = 1, turn = 1 | verse = 52, day = 0, hour = 15.5, tz = 15, matchday = 0, turn = 0 | verse = 65, day = 0, hour = 18.75, tz = 18, matchday = 0, turn = 1 | verse = 78, day = 0, hour = 22, tz = 12, matchday = 1, turn = 0 | verse = 91, day = 0, hour = 1.25, tz = 15, matchday = 1, turn = 1 | verse = 104, day = 1, hour = 4.5, tz = 4, matchday = 2, turn = 0 | verse = 117, day = 1, hour = 7.75, tz = 7, matchday = 2, turn = 1 | verse = 130, day = 1, hour = 11, tz = 1, matchday = 3, turn = 0 | verse = 143, day = 1, hour = 14.25, tz = 4, matchday = 3, turn = 1 | verse = 156, day = 1, hour = 17.5, tz = 17, matchday = 2, turn = 0 | verse = 169, day = 1, hour = 20.75, tz = 20, matchday = 2, turn = 1 | verse = 182, day = 1, hour = 0, tz = 14, matchday = 3, turn = 0 | verse = 195, day = 2, hour = 3.25, tz = 17, matchday = 5, turn = 1 | verse = 208, day = 2, hour = 6.5, tz = 6, matchday = 4, turn = 0 | verse = 221, day = 2, hour = 9.75, tz = 9, matchday = 4, turn = 1 | verse = 234, day = 2, hour = 13, tz = 3, matchday = 5, turn = 0 | verse = 247, day = 2, hour = 16.25, tz = 6, matchday = 5, turn = 1 | verse = 260, day = 2, hour = 19.5, tz = 19, matchday = 4, turn = 0 | verse = 273, day = 2, hour = 22.75, tz = 22, matchday = 4, turn = 1 | verse = 286, day = 2, hour = 2, tz = 16, matchday = 5, turn = 0 | verse = 299, day = 3, hour = 5.25, tz = 19, matchday = 7, turn = 1 | verse = 312, day = 3, hour = 8.5, tz = 8, matchday = 6, turn = 0 | verse = 325, day = 3, hour = 11.75, tz = 11, matchday = 6, turn = 1 | verse = 338, day = 3, hour = 15, tz = 5, matchday = 7, turn = 0 | verse = 351, day = 3, hour = 18.25, tz = 8, matchday = 7, turn = 1 | verse = 364, day = 3, hour = 21.5, tz = 21, matchday = 6, turn = 0 | verse = 377, day = 3, hour = 0.75, tz = 24, matchday = 6, turn = 1 | verse = 390, day = 4, hour = 4, tz = 18, matchday = 9, turn = 0 | verse = 403, day = 4, hour = 7.25, tz = 21, matchday = 9, turn = 1 | verse = 416, day = 4, hour = 10.5, tz = 10, matchday = 8, turn = 0 | verse = 429, day = 4, hour = 13.75, tz = 13, matchday = 8, turn = 1 | verse = 442, day = 4, hour = 17, tz = 7, matchday = 9, turn = 0 | verse = 455, day = 4, hour = 20.25, tz = 10, matchday = 9, turn = 1 | verse = 468, day = 4, hour = 23.5, tz = 23, matchday = 8, turn = 0 | verse = 481, day = 5, hour = 2.75, tz = 2, matchday = 10, turn = 1 | verse = 494, day = 5, hour = 6, tz = 20, matchday = 11, turn = 0 | verse = 507, day = 5, hour = 9.25, tz = 23, matchday = 11, turn = 1 | verse = 520, day = 5, hour = 12.5, tz = 12, matchday = 10, turn = 0 | verse = 533, day = 5, hour = 15.75, tz = 15, matchday = 10, turn = 1 | verse = 546, day = 5, hour = 19, tz = 9, matchday = 11, turn = 0 | verse = 559, day = 5, hour = 22.25, tz = 12, matchday = 11, turn = 1 | verse = 572, day = 5, hour = 1.5, tz = 1, matchday = 10, turn = 0 | verse = 585, day = 6, hour = 4.75, tz = 4, matchday = 12, turn = 1 | verse = 598, day = 6, hour = 8, tz = 22, matchday = 13, turn = 0 | verse = 611, day = 6, hour = 11.25, tz = 1, matchday = 13, turn = 1 | verse = 624, day = 6, hour = 14.5, tz = 14, matchday = 12, turn = 0 | verse = 637, day = 6, hour = 17.75, tz = 17, matchday = 12, turn = 1 | verse = 650, day = 6, hour = 21, tz = 11, matchday = 13, turn = 0 | verse = 663, day = 6, hour = 0.25, tz = 14, matchday = 13, turn = 1 | verse = 676, day = 7, hour = 3.5, tz = 3, matchday = 0, turn = 0 | verse = 689, day = 7, hour = 6.75, tz = 6, matchday = 0, turn = 1";
        // result.should.be.equal(expected);
    });
    
    return;
    
    
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
    
    // it('submitActions to timezone too early', async () =>  {
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
    // });

    it('submitActions to timezone', async () =>  {
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        verseBefore = await updates.getCurrentVerse().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        await moveToNextVerse(updates, extraTime = -10)        
        // await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
        await timeTravel.advanceTime(20);
        await timeTravel.advanceBlock().should.be.fulfilled;
        const cif = "ciao";
        tx = await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboys"), cif).should.be.fulfilled;
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
        // await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
    });

    it('update Timezone once', async () =>  {
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        await moveToNextVerse(updates, extraSecs = -10);
        // await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.rejected;
        // await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
        await timeTravel.advanceTime(20);
        // await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.rejected;
        const cif = "ciao2";
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), cif).should.be.fulfilled;
        timeZoneToUpdate = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        now = await updates.getNow().should.be.fulfilled;
        await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.fulfilled;
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
        console.log("warning: the next test lasts about 20 secs...")
        await moveToNextVerse(updates, extraSecs = 10);
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        const cif = "ciao3";
        for (verse = 0; verse < 110; verse++) {
            await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), cif).should.be.fulfilled;
            // await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.fulfilled;
            await moveToNextVerse(updates, extraSecs = 10);
        }
    });

    // it('timeZoneToUpdate selected edge choices', async () =>  {
    //     result = await updates._timeZoneToUpdatePure.call(verse = 0, TZ1 = 1).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(1)
    //     result.day.toNumber().should.be.equal(0);
    //     result.turnInDay.toNumber().should.be.equal(0);
    //     result = await updates._timeZoneToUpdatePure.call(verse = 3, TZ1 = 1).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(1)
    //     result.day.toNumber().should.be.equal(0);
    //     result = await updates._timeZoneToUpdatePure.call(verse = 4, TZ1 = 1).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(2)
    //     result.day.toNumber().should.be.equal(0);
    //     result.turnInDay.toNumber().should.be.equal(0);
    //     result = await updates._timeZoneToUpdatePure.call(verse = 95, TZ1 = 4).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(3);
    //     result.day.toNumber().should.be.equal(0);
    //     result.turnInDay.toNumber().should.be.equal(3);
    //     result = await updates._timeZoneToUpdatePure.call(verse = 96, TZ1 = 2).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(2);
    //     result = await updates._timeZoneToUpdatePure.call(verse = 97, TZ1 = 24).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(24);
    //     result.day.toNumber().should.be.equal(1);
    //     result.turnInDay.toNumber().should.be.equal(1);
    //     result = await updates._timeZoneToUpdatePure.call(verse = 1343, TZ1 = 24).should.be.fulfilled;
    //     result.timeZone.toNumber().should.be.equal(23);
    //     result.day.toNumber().should.be.equal(13);
    //     result.turnInDay.toNumber().should.be.equal(3);
    // });


});