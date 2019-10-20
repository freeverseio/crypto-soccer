const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingMatchLog = artifacts.require('EncodingMatchLog');

contract('EncodingMatchLog', (accounts) => {

    beforeEach(async () => {
        encoding = await EncodingMatchLog.new().should.be.fulfilled;
    });
    
    it('encode and decode matchlog', async () =>  {
        log = 0;
        log = await encoding.addNGoals(log, nGoals = 3);
        assistersIdx = Array.from(new Array(14), (x,i) => i);
        for (p = 0; p <assistersIdx.length; p++) log = await encoding.addAssister(log, assistersIdx[p], p);
        shootersIdx  = Array.from(new Array(14), (x,i) => 15-i);
        for (p = 0; p <shootersIdx.length; p++) log = await encoding.addShooter(log, shootersIdx[p], p);
        shooterForwardPos  = Array.from(new Array(14), (x,i) => i % 4);
        for (p = 0; p <shootersIdx.length; p++) log = await encoding.addForwardPos(log, shooterForwardPos[p], p);
        penalties  = Array.from(new Array(7), (x,i) => (i % 2 == 0));
        for (p = 0; p <penalties.length; p++) log = await encoding.addPenalty(log, penalties[p], p);
        outOfGames = [10, 4];
        outOfGameRounds = [7, 4];
        typesOutOfGames = [1, 2];
        is2ndHalfs = [false, true];
        for (p = 0; p <outOfGames.length; p++) {
            log = await encoding.addOutOfGame(log, outOfGames[p], outOfGameRounds[p], typesOutOfGames[p], is2ndHalfs[p]);
        }
        yellowCardedDidNotFinish1stHalf = [false, true];
        for (p = 0; p <penalties.length; p++) log = await encoding.setYellowedDidNotFinished1stHalf(log, yellowCardedDidNotFinish1stHalf[p], p);

        isHomeStadium = [true];
        if (isHomeStadium) log = await encoding.addIsHomeStadium(log)

        ingameSubs1 = [0, 1, 2];
        ingameSubs2 = [2, 1, 0];
        
        for (p = 0; p <ingameSubs1.length; p++) log = await encoding.setInGameSubsHappened(log, ingameSubs1[p], p, is2nd = false);
        for (p = 0; p <ingameSubs2.length; p++) log = await encoding.setInGameSubsHappened(log, ingameSubs2[p], p, is2nd = true);

        yellowCards1 = [9, 6];
        yellowCards2 = [3, 0];

        for (p = 0; p <yellowCards1.length; p++) log = await encoding.addYellowCard(log, yellowCards1[p], p, is2nd = false);
        for (p = 0; p <yellowCards2.length; p++) log = await encoding.addYellowCard(log, yellowCards2[p], p, is2nd = true);
        
        halfTimeSubstitutions = [9, 7, 10];
        
        for (p = 0; p <halfTimeSubstitutions.length; p++) log = await encoding.addHalfTimeSubs(log, halfTimeSubstitutions[p]);
        
        numDefTotWinner = [10, 4, 3, 1]; // [nDefsHalf1, nDefsHalf2, nTotHalf2, winner]
        
        log = await encoding.addHalfTimeSubs(log, nDefs1 = 4, is2nd = false);
        log = await encoding.addHalfTimeSubs(log, nDefs2 = 3, is2nd = true);
        log = await encoding.addNTot2ndHalf(log, nTot = 10, is2nd = false);
        log = await encoding.addWinner(log, winner = 1);

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