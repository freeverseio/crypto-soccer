const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const timeTravel = require('../utils/TimeTravel.js');

const Updates = artifacts.require('Updates');

contract('Updates', (accounts) => {
    const VERSES_PER_DAY = 96;
    const VERSES_PER_ROUND = 96*14;

    const it2 = async(text, f) => {};

    function normalizeTZ(tz) {
        return 1 + ((tz - 1) % 24);
    }

    const moveToNextVerse = async (updates, extraSecs = 0) => {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.nextVerseTimestamp().should.be.fulfilled;
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
        updates = await Updates.new().should.be.fulfilled;
        await updates.initUpdates().should.be.fulfilled;
        NULL_TIMEZONE = await updates.NULL_TIMEZONE().should.be.fulfilled;
        NULL_TIMEZONE = NULL_TIMEZONE.toNumber();
        snapShot = await timeTravel.takeSnapshot();
        snapshotId = snapShot['result'];
        });

    afterEach(async() => {
        await timeTravel.revertToSnapShot(snapshotId);
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
        nextVerseTimestamp = await updates.nextVerseTimestamp().should.be.fulfilled;
        timeZoneForRound1 = await updates.timeZoneForRound1().should.be.fulfilled;
        localTimeMs = Date.now();
        now = new Date(localTimeMs);
        nextVerse = new Date(nextVerseTimestamp.toNumber() * 1000);
        if (now.getUTCMinutes() < 57) {
            expectedHour = now.getUTCHours() + 1;
        } else {
            expectedHour = now.getUTCHours() + 2;
        }
        nextVerse.getUTCFullYear().should.be.equal(now.getUTCFullYear());
        nextVerse.getUTCMonth().should.be.equal(now.getUTCMonth());
        nextVerse.getUTCDate().should.be.equal(now.getUTCDate());
        nextVerse.getUTCHours().should.be.equal(expectedHour);
        nextVerse.getUTCMinutes().should.be.equal(0);
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
        verseBefore = await updates.currentVerse().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        await moveToNextVerse(updates, extraTime = -10)        
        // await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
        await timeTravel.advanceTime(20);
        await timeTravel.advanceBlock().should.be.fulfilled;
        const cif = "ciao";
        tx = await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboys"), cif).should.be.fulfilled;
        timeZoneToUpdate = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        verse = await updates.currentVerse().should.be.fulfilled;
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
        now = await updates.getNow().should.be.fulfilled;
        await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.fulfilled;
        submissionTime = await updates.getLastActionsSubmissionTime(timeZoneToUpdateBefore[0].toNumber()).should.be.fulfilled;
        timeZoneToUpdateAfter = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        isCloseEnough(timeZoneToUpdate[0].toNumber(), timeZoneToUpdateBefore[0].toNumber()).should.be.equal(true)
        isCloseEnough(submissionTime.toNumber(), now.toNumber()).should.be.equal(true)
    });

    it('moveToNextVerse', async () =>  {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.nextVerseTimestamp().should.be.fulfilled;
        (nextTime - now > 0).should.be.equal(true)
        await moveToNextVerse(updates, extraSecs = 0);
        now = await updates.getNow().should.be.fulfilled;
        (nextTime - now > 0).should.be.equal(false)
        
    });

    it('update Timezone many times', async () =>  {
        console.log("warning: the next test lasts about 20 secs...")
        await moveToNextVerse(updates, extraSecs = 10);
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        for (verse = 0; verse < 110; verse++) {
            timeZoneToUpdateAfter = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
            diff = timeZoneToUpdateAfter[0].toNumber() - timeZoneToUpdateBefore[0].toNumber();
            // timezone number wraps around: tz = 23, 24, 1, 2...
            if (diff < 0) diff += 24;
            // you change timezone every 4 verses
            diff.should.be.equal(Math.floor(verse / 4) % 24);
            const cif = "ciao3";
            await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), cif).should.be.fulfilled;
            await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.fulfilled;
            await moveToNextVerse(updates, extraSecs = 10);
        }
    });

    it('timeZoneToUpdate selected edge choices', async () =>  {
        result = await updates._timeZoneToUpdatePure.call(verse = 0, TZ1 = 1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(1)
        result.day.toNumber().should.be.equal(0);
        result.turnInDay.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 3, TZ1 = 1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(1)
        result.day.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 4, TZ1 = 1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(2)
        result.day.toNumber().should.be.equal(0);
        result.turnInDay.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 95, TZ1 = 4).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(3);
        result.day.toNumber().should.be.equal(0);
        result.turnInDay.toNumber().should.be.equal(3);
        result = await updates._timeZoneToUpdatePure.call(verse = 96, TZ1 = 2).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(2);
        result = await updates._timeZoneToUpdatePure.call(verse = 97, TZ1 = 24).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(24);
        result.day.toNumber().should.be.equal(1);
        result.turnInDay.toNumber().should.be.equal(1);
        result = await updates._timeZoneToUpdatePure.call(verse = 1343, TZ1 = 24).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(23);
        result.day.toNumber().should.be.equal(13);
        result.turnInDay.toNumber().should.be.equal(3);
    });

    it('timeZoneToUpdate exhaustive', async () =>  {
        console.log("warning: the next test lasts about 5 sec")
        TZ1 = 3;
        for (verse = 0; verse < VERSES_PER_ROUND; verse+=7){
            result = await updates._timeZoneToUpdatePure.call(verse, TZ1).should.be.fulfilled;
            result.timeZone.toNumber().should.be.equal(normalizeTZ(TZ1 + Math.floor(verse/4)) )
            result.day.toNumber().should.be.equal(Math.floor(verse/VERSES_PER_DAY));
            result.turnInDay.toNumber().should.be.equal(verse % 4)
        } 
    });

});