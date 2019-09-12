const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const timeTravel = require('../utils/TimeTravel.js');

const Assets = artifacts.require('Assets');
const Updates = artifacts.require('Updates');

contract('Updates', (accounts) => {
    const VERSES_PER_DAY = 96;
    const VERSES_PER_ROUND = 96*16;
   
    const moveToNextVerse = async (updates, extraSecs = 0) => {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.nextVerseTimestamp().should.be.fulfilled;
        await timeTravel.advanceTime(nextTime - now + extraSecs);
        await timeTravel.advanceBlock().should.be.fulfilled;
    };

    
    beforeEach(async () => {
        assets = await Assets.new().should.be.fulfilled;
        encoding = assets;
        await assets.init().should.be.fulfilled;
        updates = await Updates.new().should.be.fulfilled;
        await updates.initUpdates(assets.address).should.be.fulfilled;
        snapShot = await timeTravel.takeSnapshot();
        snapshotId = snapShot['result'];
        });

    afterEach(async() => {
        await timeTravel.revertToSnapShot(snapshotId);
    });
            

    it('check BC has the correct time', async () =>  {
        nextVerseTimestamp = await updates.nextVerseTimestamp().should.be.fulfilled;
        timeZoneForRound1 = await updates.timeZoneForRound1().should.be.fulfilled;
        nextVerse = new Date(nextVerseTimestamp.toNumber() * 1000)
        now = new Date()
        if (now.getUTCMinutes() < 42) {
            expectedHour = now.getUTCHours();
        } else {
            expectedHour = now.getUTCHours() + 1;
        }
        nextVerse.getUTCFullYear().should.be.equal(now.getUTCFullYear())
        nextVerse.getUTCMonth().should.be.equal(now.getUTCMonth())
        nextVerse.getUTCDate().should.be.equal(now.getUTCDate())
        nextVerse.getUTCHours().should.be.equal(expectedHour)
        nextVerse.getUTCMinutes().should.be.equal(45)
        nextVerse.getUTCSeconds().should.be.equal(0)
        timeZoneForRound1.toNumber().should.be.equal(1 + expectedHour)
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
        (newNow.toNumber() - now.toNumber() - extraTime < 4).should.be.equal(true); // we should be within 4 secs of (before + exrtaTime)
        await timeTravel.revertToSnapShot(snapshotId);
        newNow = await updates.getNow().should.be.fulfilled;
        newNow.toNumber().should.be.equal(now.toNumber());
    });

    it('submitActions to timezone too early', async () =>  {
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
    });

    it('submitActions to timezone', async () =>  {
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        verseBefore = await updates.currentVerse().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.nextVerseTimestamp().should.be.fulfilled;
        (nextTime-now > 0).should.be.equal(true);
        await timeTravel.advanceTime(nextTime-now);
        await timeTravel.advanceBlock().should.be.fulfilled;
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
        await timeTravel.advanceTime(1);
        await timeTravel.advanceBlock().should.be.fulfilled;
        tx = await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboys")).should.be.fulfilled;
        timeZoneToUpdate = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        verse = await updates.currentVerse().should.be.fulfilled;
        verse.toNumber().should.be.equal(verseBefore.toNumber() + 1); 
        timeZoneToUpdate[0].toNumber().should.be.equal(timeZoneToUpdateBefore[0].toNumber()); // tz to update does not change during the first 4 verses
        seed1 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        seed1.should.not.be.equal(seed0);
        now = await updates.getNow().should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "ActionsSubmission", (event) => {
            return event.seed == seed1 && event.submissionTime.toNumber() == now.toNumber() && event.timeZone.toNumber() == timeZoneToUpdate[0].toNumber();
        });
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
    });

    it('update Timezone once', async () =>  {
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        seed0 = await updates.getCurrentVerseSeed().should.be.fulfilled;
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.nextVerseTimestamp().should.be.fulfilled;
        (nextTime-now > 0).should.be.equal(true);
        await timeTravel.advanceTime(nextTime-now);
        await timeTravel.advanceBlock().should.be.fulfilled;
        await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.rejected;
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
        await timeTravel.advanceTime(1);
        await timeTravel.advanceBlock().should.be.fulfilled;
        await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.rejected;
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.fulfilled;
        now = await updates.getNow().should.be.fulfilled;
        await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.fulfilled;
        submissionTime = await assets.getLastActionsSubmissionTime(timeZoneToUpdateBefore[0].toNumber()).should.be.fulfilled;
        timeZoneToUpdateAfter = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        timeZoneToUpdate[0].toNumber().should.be.equal(timeZoneToUpdateBefore[0].toNumber()); // tz to update does not change during the first 4 verses
        submissionTime.should.be.bignumber.equal(now);        
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
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        for (verse = 0; verse < 4; verse++) {
            await moveToNextVerse(updates, extraSecs = 10);
            await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.fulfilled;
            await updates.updateTZ(root =  web3.utils.keccak256("hiboyz")).should.be.fulfilled;
            timeZoneToUpdateAfter = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
            diff = timeZoneToUpdateAfter[0].toNumber() - timeZoneToUpdateBefore[0].toNumber();
            console.log(verse, '  ', diff, '  ', timeZoneToUpdateAfter[0].toNumber())
            // diff.should.be.equal(verse % 4);
        }
    });

    it('timeZoneToUpdate selected edge choices', async () =>  {
        result = await updates._timeZoneToUpdatePure.call(verse = 0, TZ1 = 1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(1)
        result.day.toNumber().should.be.equal(1);
        result.turnInDay.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 3, TZ1 = 1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(1)
        result.day.toNumber().should.be.equal(1);
        result = await updates._timeZoneToUpdatePure.call(verse = 4, TZ1 = 1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(2)
        result.day.toNumber().should.be.equal(1);
        result.turnInDay.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 95, TZ1 = 4).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(3);
        result.day.toNumber().should.be.equal(1);
        result.turnInDay.toNumber().should.be.equal(3);
        result = await updates._timeZoneToUpdatePure.call(verse = 96, TZ1 = 2).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 97, TZ1 = 24).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(24);
        result.day.toNumber().should.be.equal(2);
        result.turnInDay.toNumber().should.be.equal(0);
        result = await updates._timeZoneToUpdatePure.call(verse = 1535, TZ1 = 24).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(23);
        result.day.toNumber().should.be.equal(16);
        result.turnInDay.toNumber().should.be.equal(2);
    });

    return;
    
    
    it('timeZoneToUpdate exhaustive', async () =>  {
        TZ1 = 3;
        // for day 1
        for (verse = 0; verse < VERSES_PER_DAY; verse+=5){
            result = await updates._timeZoneToUpdatePure.call(verse, TZ1).should.be.fulfilled;
            result.timeZone.toNumber().should.be.equal(1+(TZ1 - 1 + Math.floor(verse / 4))%24)
            result.day.toNumber().should.be.equal(1);
            result.turnInDay.toNumber().should.be.equal(verse % 4)
        } 
        // there is an empty spot at the end of day 1
        result = await updates._timeZoneToUpdatePure(VERSES_PER_DAY, TZ1).should.be.fulfilled;
        result.timeZone.toNumber().should.be.equal(0)
        // beyond day 1
        for (verse = VERSES_PER_DAY + 1; verse < VERSES_PER_ROUND; verse+=7){
            result = await updates._timeZoneToUpdatePure.call(verse, TZ1).should.be.fulfilled;
            result.timeZone.toNumber().should.be.equal(1+(TZ1 - 1 + Math.floor((verse-1) / 4))%24);
            result.day.toNumber().should.be.equal(1 + Math.floor((verse - 1) / VERSES_PER_DAY));
            result.turnInDay.toNumber().should.be.equal( (verse-1) % 4);
        } 
    });



});