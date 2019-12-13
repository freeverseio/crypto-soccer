const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const UtilsMatchLog = artifacts.require('UtilsMatchLog');
const logUtils = require('../utils/matchLogUtils.js');
const debug = require('../utils/debugUtils.js');

contract('EncodingMatchLog', (accounts) => {

    const UNDEF = undefined;
    const MAX_GOALS = 12;
    
    beforeEach(async () => {
        encoding = await EncodingMatchLog.new().should.be.fulfilled;
        utilsML = await UtilsMatchLog.new().should.be.fulfilled;
    });
    
    it('encode and decode matchlog', async () =>  {
        nGoals = 3;
        assistersIdx = Array.from(new Array(MAX_GOALS), (x,i) => i);
        shootersIdx  = Array.from(new Array(MAX_GOALS), (x,i) => 15-i);
        shooterForwardPos  = Array.from(new Array(MAX_GOALS), (x,i) => i % 4);
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
        teamSumSkills = 1;
        trainingPoints = 2333;
        
        log = await logUtils.encodeLog(encoding, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPoints);

        await logUtils.checkExpectedLog(encoding, log, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPoints);
            
            
    //  teamSumSkills 
    //  winner: 0 = home, 1 = away, 2 = draw
    //  nGoals
    //  trainingPoints
    //  uint8 memory outOfGames
    //  uint8 memory typesOutOfGames, 
    //  uint8 memory outOfGameRounds
    //  uint8[2] memory yellowCards
    //  uint8[3] memory ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen  
    //  uint8[3] memory halfTimeSubstitutions: 0...10 the player in the starting 11 that was changed during half time            
    // HALF 1
    result = await utilsML.fullDecodeMatchLog(log, is2ndHalf = false).should.be.fulfilled;
    expected = [
        teamSumSkills,
        winner,
        nGoals,
        trainingPoints,
        outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
        yellowCards1[0], yellowCards1[1],
        ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
        0, 0, 0
    ]
    debug.compareArrays(result, expected, toNum = true, verbose = false);

    // HALF 2
    result = await utilsML.fullDecodeMatchLog(log, is2ndHalf = true).should.be.fulfilled;
    expected = [
        teamSumSkills,
        winner,
        nGoals,
        trainingPoints,
        outOfGames[1], typesOutOfGames[1], outOfGameRounds[1],
        yellowCards2[0], yellowCards2[1],
        ingameSubs2[0], ingameSubs2[1], ingameSubs2[2],
        halfTimeSubstitutions[0], halfTimeSubstitutions[1], halfTimeSubstitutions[2]
    ]
    debug.compareArrays(result, expected, toNum = true, verbose = false);
});
    

});