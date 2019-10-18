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
        nGoals = 3;
        assistersIdx = Array.from(new Array(14), (x,i) => i);
        shootersIdx  = Array.from(new Array(14), (x,i) => 15-i);
        shooterForwardPos  = Array.from(new Array(14), (x,i) => i % 4);
        assistersShootersForwardsPos = assistersIdx.concat(shootersIdx).concat(shooterForwardPos);
        penalties  = Array.from(new Array(7), (x,i) => (i % 2 == 0));
        typesOutOfGames = [1, 2];
        outOfGameRounds = [7, 4];
        yellowCardedDidNotFinish1stHalf = [false, true];
        isHomeStadium = [true];
        yellowCardedDidNotFinish1stHalfAndIsHomeStadium = yellowCardedDidNotFinish1stHalf.concat(isHomeStadium);
        ingameSubs = [
            0, 1, 2,  // half 1
            2, 1, 0,   // half 2
        ]
        outOfGamesAndYellowCards = [10, 4, 9, 6, 3, 0];
        // made out from these two:
        //      outOfGames = [10, 4];
        //      yellowCards = [9, 6, 3, 0];

        halfTimeSubstitutions = [9, 7, 10];
        numDefTotWinner = [10, 4, 3, 1]; // [nDefsHalf1, nDefsHalf2, nTotHalf2, winner]
        result = await encoding.encodeMatchLog(
            nGoals, 
            assistersShootersForwardsPos, 
            penalties, 
            outOfGamesAndYellowCards, 
            outOfGameRounds, 
            typesOutOfGames, 
            yellowCardedDidNotFinish1stHalfAndIsHomeStadium,
            halfTimeSubstitutions, 
            ingameSubs,
            numDefTotWinner
        );
        result = await encoding.decodeMatchLog(result);
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