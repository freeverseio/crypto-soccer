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
        penalties  = Array.from(new Array(7), (x,i) => (i % 2 == 0));
        typesOutOfGames = [1, 2];
        outOfGameRounds = [7, 4];
        yellowCardedFinished1stHalf = [false, true];
        ingameSubsCancelled = [
            false, true, false,  // half 1
            true, true, false,   // half 2
        ]
        outOfGamesAndYellowCards = [10, 4, 9, 6, 3, 0]
        // made out from these two:
        //      outOfGames = [10, 4];
        //      yellowCards = [9, 6, 3, 0];

        halfTimeSubstitutions = [9, 7, 10]
        result = await encoding.encodeMatchLog(
            nGoals, 
            assistersIdx, 
            shootersIdx, 
            shooterForwardPos, 
            penalties, 
            outOfGamesAndYellowCards, 
            outOfGameRounds, 
            typesOutOfGames, 
            yellowCardedFinished1stHalf,
            halfTimeSubstitutions, 
            ingameSubsCancelled
        );
        result = await encoding.decodeMatchLog(result);
        let {0: nGo, 1: ass, 2: sho, 3: fwd, 4: pen, 5: outsAndYels, 6: outRounds, 7: typ, 8: yelFin, 9: halfSubs, 10: inGameSubsCancl} = result;
        nGo.toNumber().should.be.equal(nGoals);        
        for (i = 0; i < assistersIdx.length; i++) ass[i].toNumber().should.be.equal(assistersIdx[i]); 
        for (i = 0; i < shootersIdx.length; i++) sho[i].toNumber().should.be.equal(shootersIdx[i]); 
        for (i = 0; i < shooterForwardPos.length; i++) fwd[i].toNumber().should.be.equal(shooterForwardPos[i]); 
        for (i = 0; i < penalties.length; i++) pen[i].should.be.equal(penalties[i]); 
        for (i = 0; i < outOfGamesAndYellowCards.length; i++) outsAndYels[i].toNumber().should.be.equal(outOfGamesAndYellowCards[i]); 
        for (i = 0; i < outOfGameRounds.length; i++) outRounds[i].toNumber().should.be.equal(outOfGameRounds[i]); 
        for (i = 0; i < typesOutOfGames.length; i++) typ[i].toNumber().should.be.equal(typesOutOfGames[i]); 
        for (i = 0; i < yellowCardedFinished1stHalf.length; i++) yelFin[i].should.be.equal(yellowCardedFinished1stHalf[i]); 
        for (i = 0; i < halfTimeSubstitutions.length; i++) halfSubs[i].toNumber().should.be.equal(halfTimeSubstitutions[i]); 
        for (i = 0; i < ingameSubsCancelled.length; i++) inGameSubsCancl[i].should.be.equal(ingameSubsCancelled[i]); 
    });
    

});