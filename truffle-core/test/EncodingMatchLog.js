const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingMatchLog = artifacts.require('EncodingMatchLog');

contract('EncodingMatchLog', (accounts) => {

    async function encodeLog(encoding, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
        outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
        isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
        halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner
    ) {
        log = await encoding.addNGoals(log, nGoals);
        for (p = 0; p <assistersIdx.length; p++) log = await encoding.addAssister(log, assistersIdx[p], p);
        for (p = 0; p <shootersIdx.length; p++) log = await encoding.addShooter(log, shootersIdx[p], p);
        for (p = 0; p <shootersIdx.length; p++) log = await encoding.addForwardPos(log, shooterForwardPos[p], p);
        for (p = 0; p <penalties.length; p++) log = await encoding.addPenalty(log, penalties[p], p);
        is2ndHalfs = [false, true];
        for (p = 0; p <outOfGames.length; p++) {
            log = await encoding.addOutOfGame(log, outOfGames[p], outOfGameRounds[p], typesOutOfGames[p], is2ndHalfs[p]);
        }
        for (p = 0; p <penalties.length; p++) {
            if ( yellowCardedDidNotFinish1stHalf[p]) {
                log = await encoding.setYellowedDidNotFinished1stHalf(log, p);
            }
        }
        if (isHomeStadium) log = await encoding.addIsHomeStadium(log)
        for (p = 0; p <ingameSubs1.length; p++) log = await encoding.setInGameSubsHappened(log, ingameSubs1[p], p, is2nd = false);
        for (p = 0; p <ingameSubs2.length; p++) log = await encoding.setInGameSubsHappened(log, ingameSubs2[p], p, is2nd = true);
        for (p = 0; p <yellowCards1.length; p++) log = await encoding.addYellowCard(log, yellowCards1[p], p, is2nd = false);
        for (p = 0; p <yellowCards2.length; p++) log = await encoding.addYellowCard(log, yellowCards2[p], p, is2nd = true);
        for (p = 0; p <halfTimeSubstitutions.length; p++) log = await encoding.addHalfTimeSubs(log, halfTimeSubstitutions[p], p);
        
        log = await encoding.addNDefs(log, nDefs1, is2nd = false);
        log = await encoding.addNDefs(log, nDefs2, is2nd = true);
        log = await encoding.addNTot2ndHalf(log, nTot);
        log = await encoding.addWinner(log, winner);
        return log;
    }
    
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
        isHomeStadium = [true];
        ingameSubs1 = [0, 1, 2];
        ingameSubs2 = [2, 1, 0];
        yellowCards1 = [9, 6];
        yellowCards2 = [3, 0];
        halfTimeSubstitutions = [9, 7, 10];
        nDefs1 = 4;
        nDefs2 = 3;
        nTot = 10;
        winner = 1;
        
        log = await encodeLog(encoding, nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner);

        ///////////////
        
        
        result = await encoding.decodeMatchLog(log);
        // let {0: nGo, 1: ass, 2: sho, 3: fwd, 4: pen, 5: outsAndYels, 6: outRounds, 7: typ, 8: yelFin, 9: halfSubs, 10: inGameSubs} = result;
        let {0: nGo, 1: assShoFwd, 2: pen, 3: outsAndYels, 4: outRounds, 5: typ, 6: yelFinHome, 7: halfSubs, 8: inGameSubs, 9: defTotWin} = result;
        nGo.toNumber().should.be.equal(nGoals);        
        for (i = 0; i < assistersShootersForwardsPos.length; i++) assShoFwd[i].toNumber().should.be.equal(assistersShootersForwardsPos[i]); 
        for (i = 0; i < penalties.length; i++) pen[i].should.be.equal(penalties[i]); 
        for (i = 0; i < outOfGamesAndYellowCards.length; i++) outsAndYels[i].toNumber().should.be.equal(outOfGamesAndYellowCards[i]); 
        for (i = 0; i < outOfGameRounds.length; i++) outRounds[i].toNumber().should.be.equal(outOfGameRounds[i]); 
        for (i = 0; i < typesOutOfGames.length; i++) typ[i].toNumber().should.be.equal(typesOutOfGames[i]); 
        for (i = 0; i < yellowCardedDidNotFinish1stHalfAndIsHomeStadium.length; i++) yelFinHome[i].should.be.equal(yellowCardedDidNotFinish1stHalfAndIsHomeStadium[i]); 
        for (i = 0; i < halfTimeSubstitutions.length; i++) halfSubs[i].toNumber().should.be.equal(halfTimeSubstitutions[i]); 
        for (i = 0; i < ingameSubs.length; i++) inGameSubs[i].toNumber().should.be.equal(ingameSubs[i]); 
        for (i = 0; i < numDefTotWinner.length; i++) defTotWin[i].toNumber().should.be.equal(numDefTotWinner[i]); 
    });
    

});