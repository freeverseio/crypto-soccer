const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const Encoding = artifacts.require('EncodingState');

contract('EncodingState', (accounts) => {

    beforeEach(async () => {
        encoding = await Encoding.new().should.be.fulfilled;
    });
    
    it('encode decode player state', async () => {
        const playerId = 231;
        const currentTeamId = 432432;
        const currentShirtNum = 12;
        const prevPlayerTeamId = 32123;
        const lastSaleBlock = 3221;
        // check the initial setting of a player state (from empty)
        const state = await encoding.encodePlayerState(playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock).should.be.fulfilled;
        result = await encoding.getPlayerIdFromState(state).should.be.fulfilled;
        result.toNumber().should.be.equal(playerId);
        result = await encoding.getCurrentTeamId(state).should.be.fulfilled;
        result.toNumber().should.be.equal(currentTeamId);
        result = await encoding.getCurrentShirtNum(state).should.be.fulfilled;
        result.toNumber().should.be.equal(currentShirtNum);
        result = await encoding.getPrevPlayerTeamId(state).should.be.fulfilled;
        result.toNumber().should.be.equal(prevPlayerTeamId);
        result = await encoding.getLastSaleBlock(state).should.be.fulfilled;
        result.toNumber().should.be.equal(lastSaleBlock);
        // check the individual changes (from non-empty)
        newState = await encoding.setCurrentTeamId(state, newval = 43).should.be.fulfilled;
        result = await encoding.getCurrentTeamId(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encoding.setCurrentShirtNum(state, newval = 2).should.be.fulfilled;
        result = await encoding.getCurrentShirtNum(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encoding.setPrevPlayerTeamId(state, newval = 43643).should.be.fulfilled;
        result = await encoding.getPrevPlayerTeamId(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encoding.setLastSaleBlock(state, newval = 11223).should.be.fulfilled;
        result = await encoding.getLastSaleBlock(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
    });


});