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


async function logWrapper(log) {
    // MAX_GOALS = 12
    // N_HALFS = 2
    var result = {
        encodedLog: log.toString(),
        assister: [],  // length = MAX_GOALS
        shooter: [],  // length = MAX_GOALS
        forwardPos: [],  // length = MAX_GOALS
        penalty: [],  // length = 7
        isHomeStadium: false, 
        halfTimeSubs: [], // length = 3
        nGKAndDefs: [], // length = N_HALFS
        nTot: [], // length = N_HALFS
        winner: 0,
        teamSumSkills: 0,
        trainingPoints: 0,
        nGoals: 0,
        outOfGamePlayer: [], // length = N_HALFS
        outOfGameType: [], // length = N_HALFS
        outOfGameRound: [], // length = N_HALFS
        yellowCard: [], // length = 2 * N_HALFS
        inGameSubsHappened: [], // length = 3 * N_HALFS
        changesAtHalftime: 0
    }
    result.isHomeStadium = await encoding.getIsHomeStadium(log).should.be.fulfilled;
    result.winner = Number(await encoding.getWinner(log).should.be.fulfilled);
    result.teamSumSkills = Number(await encoding.getTeamSumSkills(log).should.be.fulfilled);
    result.trainingPoints = Number(await encoding.getTrainingPoints(log).should.be.fulfilled);
    result.nGoals = Number(await encoding.getNGoals(log).should.be.fulfilled);
    result.changesAtHalftime = Number(await encoding.getChangesAtHalfTime(log).should.be.fulfilled);

    return result;
}



contract('EncodingMatchLog', (accounts) => {

    const UNDEF = undefined;
    const MAX_GOALS = 12;
    
    beforeEach(async () => {
        encoding = await EncodingMatchLog.new().should.be.fulfilled;
        utils = await Utils.new().should.be.fulfilled;
    });
    
    it('encode and decode matchlog', async () =>  {
        const writeMode = true;
        toWrite = [];

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

        if (writeMode) { await toWrite.push(await logWrapper(log))}

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
        if (writeMode) { await toWrite.push(await logWrapper(log))}
        result = await encoding.getIsHomeStadium(log).should.be.fulfilled;
        result.should.be.equal(!isHomeStadium)
        result = await encoding.getTeamSumSkills(log).should.be.fulfilled;
        result.toNumber().should.be.equal(teamSumSkills)
        
        result = await encoding.getChangesAtHalfTime(log).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await encoding.setChangesAtHalfTime(log, 3).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await logWrapper(result))}
        result = await encoding.getChangesAtHalfTime(result).should.be.fulfilled;
        result.toNumber().should.be.equal(3);

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

        const fs = require('fs');
        if (writeMode) {
            fs.writeFileSync('test/testdata/encodingMatchLogTestData.json', JSON.stringify(toWrite), function(err) {
                if (err) {
                    console.log(err);
                }
            });
        }             
        
        writtenData = fs.readFileSync('test/testdata/encodingMatchLogTestData.json', 'utf8');
        assert.equal(
            web3.utils.keccak256(writtenData),
            "0x2df0638cf5da6f466cfe2ce30bc695ccb212981f270add8ba80ef8884bfb2640",
            "written testdata for encoding MatchLog does not match expected result"
        );
    });
});