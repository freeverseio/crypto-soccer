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
        outOfGames = [10, 4];
        typesOutOfGames = [1, 2];
        yellowCards = [9, 6, 3, 0];
        result = await encoding.encodeMatchLog(nGoals, assistersIdx, shootersIdx, shooterForwardPos, penalties, outOfGames, typesOutOfGames, yellowCards);
        result = await encoding.decodeMatchLog(result);
        let {0: nGo, 1: ass, 2: sho, 3: fwd, 4: pen, 5: out, 6: typ, 7: yel} = result;
        nGo.toNumber().should.be.equal(nGoals);        
        for (i = 0; i < assistersIdx.length; i++) ass[i].toNumber().should.be.equal(assistersIdx[i]); 
        for (i = 0; i < shootersIdx.length; i++) sho[i].toNumber().should.be.equal(shootersIdx[i]); 
        for (i = 0; i < shooterForwardPos.length; i++) fwd[i].toNumber().should.be.equal(shooterForwardPos[i]); 
        for (i = 0; i < penalties.length; i++) pen[i].should.be.equal(penalties[i]); 
        for (i = 0; i < outOfGames.length; i++) out[i].toNumber().should.be.equal(outOfGames[i]); 
        for (i = 0; i < typesOutOfGames.length; i++) typ[i].toNumber().should.be.equal(typesOutOfGames[i]); 
        for (i = 0; i < yellowCards.length; i++) yel[i].toNumber().should.be.equal(yellowCards[i]); 
    });
    

});