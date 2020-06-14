/*
 Tests for all functions in EncodingMatchLog.sol and contracts inherited by it
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const Utils = artifacts.require('Utils');
const logUtils = require('../utils/matchLogUtils.js');
const debug = require('../utils/debugUtils.js');

contract('EncodingMatchLog', (accounts) => {

    const UNDEF = undefined;
    const MAX_GOALS = 12;
    
    beforeEach(async () => {
        encoding = await EncodingMatchLog.new().should.be.fulfilled;
        utils = await Utils.new().should.be.fulfilled;
    });
    
    it('encode and decode matchlog', async () =>  {
        nGoals = 15;
        assistersIdx = Array.from(new Array(MAX_GOALS), (x,i) => 15-i%4);
        shootersIdx  = Array.from(new Array(MAX_GOALS), (x,i) => 15-i%4);
        shooterForwardPos  = Array.from(new Array(MAX_GOALS), (x,i) => i % 4);
        penalties  = Array.from(new Array(7), (x,i) => (i % 2 == 0));
        outOfGames = [14, 13];
        outOfGameRounds = [14, 15];
        typesOutOfGames = [2, 3];
        isHomeStadium = true;
        ingameSubs1 = [3, 2, 3];
        ingameSubs2 = [2, 3, 2];
        yellowCards1 = [14, 15];
        yellowCards2 = [15, 14];
        halfTimeSubstitutions = [31, 30, 31];
        nGKAndDefs1 = 14;
        nGKAndDefs2 = 15;
        nTot1 = 15;
        nTot2 = 14;
        winner = 3;
        teamSumSkills = (2**24)-1;
        trainingPoints = (2**12)-1;
        
        log = await logUtils.encodeLog(encoding, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, 
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nGKAndDefs1, nGKAndDefs2, nTot1, nTot2, winner, teamSumSkills, trainingPoints
        );

        await logUtils.checkExpectedLog(encoding, log, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, 
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nGKAndDefs1, nGKAndDefs2, nTot1, nTot2, winner, teamSumSkills, trainingPoints
        );
            
        // mini test that once showed a bug:
        result = await encoding.getIsHomeStadium(log).should.be.fulfilled;
        result.should.be.equal(isHomeStadium)
        result = await encoding.getTeamSumSkills(log).should.be.fulfilled;
        result.toNumber().should.be.equal(teamSumSkills)
        log = await encoding.setIsHomeStadium(log, !isHomeStadium).should.be.fulfilled;
        result = await encoding.getIsHomeStadium(log).should.be.fulfilled;
        result.should.be.equal(!isHomeStadium)
        result = await encoding.getTeamSumSkills(log).should.be.fulfilled;
        result.toNumber().should.be.equal(teamSumSkills)
        
        // HALF 1
        result = await utils.fullDecodeMatchLog(log, is2ndHalf = false).should.be.fulfilled;
        expected = [
            teamSumSkills,
            winner,
            nGoals,
            trainingPoints1stHalf = 0,
            outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
            yellowCards1[0], yellowCards1[1],
            ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
            0, 0, 0
        ]
        debug.compareArrays(result, expected, toNum = true);

        // HALF 2
        result = await utils.fullDecodeMatchLog(log, is2ndHalf = true).should.be.fulfilled;
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
        debug.compareArrays(result, expected, toNum = true);
    });
});