const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const timeTravel = require('../utils/TimeTravel.js');
const delegateUtils = require('../utils/delegateCallUtils.js');
const merkleUtils = require('../utils/merkleUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Merkle = artifacts.require('Merkle');



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
        merkle = await Merkle.new().should.be.fulfilled;
        proxy = await Proxy.new(delegateUtils.extractSelectorsFromAbi(Proxy.abi)).should.be.fulfilled;
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
    
    it2('check timezones for this verse', async () =>  {
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
        // console.log(result)
        expected = " | verse = 0, tz = 2, matchday = 0, turn = 0 | verse = 13, tz = 5, matchday = 0, turn = 1 | verse = 26, tz = 0, matchday = 0, turn = 0 | verse = 39, tz = 2, matchday = 1, turn = 1 | verse = 52, tz = 15, matchday = 0, turn = 0 | verse = 65, tz = 18, matchday = 0, turn = 1 | verse = 78, tz = 12, matchday = 1, turn = 0 | verse = 91, tz = 15, matchday = 1, turn = 1 | verse = 104, tz = 4, matchday = 2, turn = 0 | verse = 117, tz = 7, matchday = 2, turn = 1 | verse = 130, tz = 1, matchday = 1, turn = 0 | verse = 143, tz = 4, matchday = 3, turn = 1 | verse = 156, tz = 17, matchday = 2, turn = 0 | verse = 169, tz = 20, matchday = 2, turn = 1 | verse = 182, tz = 14, matchday = 3, turn = 0 | verse = 195, tz = 17, matchday = 3, turn = 1 | verse = 208, tz = 6, matchday = 4, turn = 0 | verse = 221, tz = 9, matchday = 4, turn = 1 | verse = 234, tz = 3, matchday = 5, turn = 0 | verse = 247, tz = 6, matchday = 5, turn = 1 | verse = 260, tz = 19, matchday = 4, turn = 0 | verse = 273, tz = 22, matchday = 4, turn = 1 | verse = 286, tz = 16, matchday = 5, turn = 0 | verse = 299, tz = 19, matchday = 5, turn = 1 | verse = 312, tz = 8, matchday = 6, turn = 0 | verse = 325, tz = 11, matchday = 6, turn = 1 | verse = 338, tz = 5, matchday = 7, turn = 0 | verse = 351, tz = 8, matchday = 7, turn = 1 | verse = 364, tz = 21, matchday = 6, turn = 0 | verse = 377, tz = 24, matchday = 6, turn = 1 | verse = 390, tz = 18, matchday = 7, turn = 0 | verse = 403, tz = 21, matchday = 7, turn = 1 | verse = 416, tz = 10, matchday = 8, turn = 0 | verse = 429, tz = 13, matchday = 8, turn = 1 | verse = 442, tz = 7, matchday = 9, turn = 0 | verse = 455, tz = 10, matchday = 9, turn = 1 | verse = 468, tz = 23, matchday = 8, turn = 0 | verse = 481, tz = 2, matchday = 10, turn = 1 | verse = 494, tz = 20, matchday = 9, turn = 0 | verse = 507, tz = 23, matchday = 9, turn = 1 | verse = 520, tz = 12, matchday = 10, turn = 0 | verse = 533, tz = 15, matchday = 10, turn = 1 | verse = 546, tz = 9, matchday = 11, turn = 0 | verse = 559, tz = 12, matchday = 11, turn = 1 | verse = 572, tz = 1, matchday = 10, turn = 0 | verse = 585, tz = 4, matchday = 12, turn = 1 | verse = 598, tz = 22, matchday = 11, turn = 0 | verse = 611, tz = 1, matchday = 11, turn = 1 | verse = 624, tz = 14, matchday = 12, turn = 0 | verse = 637, tz = 17, matchday = 12, turn = 1 | verse = 650, tz = 11, matchday = 13, turn = 0 | verse = 663, tz = 14, matchday = 13, turn = 1 | verse = 676, tz = 3, matchday = 0, turn = 0 | verse = 689, tz = 6, matchday = 0, turn = 1 | verse = 702, tz = 24, matchday = 13, turn = 0 | verse = 715, tz = 3, matchday = 1, turn = 1 | verse = 728, tz = 16, matchday = 0, turn = 0 | verse = 741, tz = 19, matchday = 0, turn = 1 | verse = 754, tz = 13, matchday = 1, turn = 0 | verse = 767, tz = 16, matchday = 1, turn = 1 | verse = 780, tz = 5, matchday = 2, turn = 0 | verse = 793, tz = 8, matchday = 2, turn = 1 | verse = 806, tz = 2, matchday = 3, turn = 0 | verse = 819, tz = 5, matchday = 3, turn = 1 | verse = 832, tz = 18, matchday = 2, turn = 0 | verse = 845, tz = 21, matchday = 2, turn = 1 | verse = 858, tz = 15, matchday = 3, turn = 0 | verse = 871, tz = 18, matchday = 3, turn = 1 | verse = 884, tz = 7, matchday = 4, turn = 0 | verse = 897, tz = 10, matchday = 4, turn = 1 | verse = 910, tz = 4, matchday = 5, turn = 0 | verse = 923, tz = 7, matchday = 5, turn = 1 | verse = 936, tz = 20, matchday = 4, turn = 0 | verse = 949, tz = 23, matchday = 4, turn = 1";
        result.should.be.equal(expected);
    });
    
    it2('require that BC and local time are less than 15 sec out of sync', async () =>  {
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
    
    it2('check BC is set up in agreement with the local time', async () =>  {
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
    
    it2('wait some minutes', async () =>  {
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
    
    // it2('submitActions to timezone too early', async () =>  {
    //     await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy")).should.be.rejected;
    // });

    it2('submitActions to timezone', async () =>  {
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

    it2('update Timezone once', async () =>  {
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

    it2('moveToNextVerse', async () =>  {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.getNextVerseTimestamp().should.be.fulfilled;
        (nextTime - now > 0).should.be.equal(true)
        await moveToNextVerse(updates, extraSecs = 0);
        now = await updates.getNow().should.be.fulfilled;
        (nextTime - now > 0).should.be.equal(false)
        
    });

    it2('update Timezone many times', async () =>  {
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
    
    it('challenging a tz', async () =>  {
        await moveToNextVerse(updates, extraSecs = 10);
        timeZoneToUpdateBefore = await updates.nextTimeZoneToUpdate().should.be.fulfilled;
        const cif = "ciao3";
        await updates.submitActionsRoot(actionsRoot =  web3.utils.keccak256("hiboy"), cif).should.be.fulfilled;

        nLeafsPerRoot = 16;
        nChallenges = 3;
        nTotalLeafs = nLeafsPerRoot**3;
        nTotalLevels = Math.log2(nTotalLeafs);
        nLevelsPerRoot = Math.log2(nLeafsPerRoot);
        leafs = Array.from(new Array(nTotalLeafs), (x,i) => web3.utils.keccak256(i.toString()));
        merkleStruct = merkleUtils.buildMerkleStruct(leafs, nLeafsPerRoot);
        const nullHash = '0x0';
        // First challenge fails because the TZ has not been updated yet with a root
        await updates.challengeTZ(wrongVal = nullHash, wrongPos = 0, proof = [], merkleStruct[1]).should.be.rejected;

        // We update with the correct root...
        await updates.updateTZ(root = merkleUtils.merkleRoot(leafs, nTotalLevels)).should.be.fulfilled;
        // ...so that we cannot challenge with the correct set of hashes
        await updates.challengeTZ(wrongVal = nullHash, wrongPos = 0, proof = [], merkleStruct[1]).should.be.rejected;
        // ...but we can challenge with one of them being wrong
        leafsWrong = [...leafs];
        // we will lie in a bottom leave that leads to root 7 in the first level
        // so being at pos = 7, leads to pos 7 * nLeafsPerRoot, which leads at 7*nLeafsPerRoot^2
        leafsWrong[7 * (nLeafsPerRoot**2) + 1] = web3.utils.keccak256('iAmEvil');
        merkleStructWrong =  merkleUtils.buildMerkleStruct(leafsWrong, nLeafsPerRoot);
        assert.notEqual(merkleUtils.merkleRoot(leafsWrong, nTotalLevels), merkleUtils.merkleRoot(leafs, nTotalLevels), "wrong leafs should lead to different root");
        assert.notEqual(merkleUtils.merkleRoot(merkleStructWrong[1], nLevelsPerRoot), merkleUtils.merkleRoot(merkleStruct[1], nLevelsPerRoot), "wrong leafs should lead to different merkle structs");
        
        await updates.challengeTZ(wrongVal = nullHash, wrongPos = 0, proof = [], merkleStructWrong[1]).should.be.fulfilled;
        // we can now challenge the challenger :-) with the correct hashes  
        // TODO: test that vals are gotten from events
        wrongPos = 7;
        wrongVal = merkleStructWrong[1][wrongPos];
        proof = merkleUtils.buildProof(wrongPos, merkleStructWrong[1], nLevelsPerRoot);
        roots2Submit = merkleStruct[2].slice(wrongPos*nLeafsPerRoot, (wrongPos+1)*nLeafsPerRoot);
        wrongRoot = merkleUtils.merkleRoot(merkleStructWrong[1], nLevelsPerRoot);
        assert.equal(merkleUtils.verify(wrongRoot, proof, wrongVal, wrongPos), true, "proof not working");
        // as always, first check that we cannot submit roots that coinicide with previous:
        roots2SubmitWrong = merkleStructWrong[2].slice(wrongPos*nLeafsPerRoot, (wrongPos+1)*nLeafsPerRoot);
        assert.equal(merkleUtils.merkleRoot(roots2SubmitWrong, nLevelsPerRoot), wrongVal, "wrong choice of slice");
        assert.notEqual(merkleUtils.merkleRoot(roots2SubmitWrong, nLevelsPerRoot), merkleUtils.merkleRoot(roots2Submit, nLevelsPerRoot), "wrong choice of slice");
        await updates.challengeTZ(wrongVal, wrongPos, proof, roots2SubmitWrong).should.be.rejected;
        // but we can with differing ones:
        await updates.challengeTZ(wrongVal, wrongPos, proof, roots2Submit).should.be.fulfilled;
    });

});