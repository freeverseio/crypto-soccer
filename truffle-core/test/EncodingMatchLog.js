const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const logUtils = require('./matchLogUtils.js');

contract('EncodingMatchLog', (accounts) => {

    const UNDEF = undefined;
    
    beforeEach(async () => {
        encoding = await EncodingMatchLog.new().should.be.fulfilled;
    });
    
    it('encode and decode matchlog', async () =>  {
        nGoals = 3;
        assistersIdx = Array.from(new Array(14), (x,i) => i);
        shootersIdx  = Array.from(new Array(14), (x,i) => 15-i);
        shooterForwardPos  = Array.from(new Array(14), (x,i) => i % 4);
        penalties  = Array.from(new Array(7), (x,i) => (i % 2 == 0));
        outOfGames = [10, 4];
        outOfGameRounds = [7, 4];
        typesOutOfGames = [1, 2];
        yellowCardedDidNotFinish1stHalf = [false, true];
        isHomeStadium = true;
        ingameSubs1 = [0, 1, 2];
        ingameSubs2 = [2, 1, 0];
        yellowCards1 = [9, 6];
        yellowCards2 = [3, 0];
        halfTimeSubstitutions = [9, 7, 10];
        nDefs1 = 4;
        nDefs2 = 3;
        nTot = 10;
        winner = 1;
        
        log = await logUtils.encodeLog(encoding, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner);

        await logUtils.checkExpectedLog(encoding, log, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner);
    });
    

});