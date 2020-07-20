/*
 Tests for all functions in EncodingState.sol and contracts inherited by it
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingState = artifacts.require('EncodingState');
const EncodingIDs = artifacts.require('EncodingIDs');

contract('EncodingState', (accounts) => {

    beforeEach(async () => {
        encodingState = await EncodingState.new().should.be.fulfilled;
        encodingIDs = await EncodingIDs.new().should.be.fulfilled;
    });
    
    it('encode decode player state', async () => {
        const currentTeamId = await encodingIDs.encodeTZCountryAndVal(tz = 1, countryIdx = 0, teamIDx = 0).should.be.fulfilled;
        const currentShirtNum = 12;
        const prevPlayerTeamId = await encodingIDs.encodeTZCountryAndVal(tz = 1, countryIdx = 0, teamIDx = 1).should.be.fulfilled;
        const lastSaleBlock = 3221;
        // check the initial setting of a player state (from empty)
        const state = await encodingState.encodePlayerState(currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock).should.be.fulfilled;
        // console.log(state.toString())
        result = await encodingState.getCurrentTeamIdFromPlayerState(state).should.be.fulfilled;
        result.should.be.bignumber.equal(currentTeamId);
        result = await encodingState.getCurrentShirtNum(state).should.be.fulfilled;
        result.toNumber().should.be.equal(currentShirtNum);
        result = await encodingState.getPrevPlayerTeamId(state).should.be.fulfilled;
        result.should.be.bignumber.equal(prevPlayerTeamId);
        result = await encodingState.getLastSaleBlock(state).should.be.fulfilled;
        result.toNumber().should.be.equal(lastSaleBlock);
        // check the individual changes (from non-empty)
        newState = await encodingState.setCurrentTeamId(state, newval = 43).should.be.fulfilled;
        result = await encodingState.getCurrentTeamIdFromPlayerState(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encodingState.setCurrentShirtNum(state, newval = 2).should.be.fulfilled;
        result = await encodingState.getCurrentShirtNum(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encodingState.setPrevPlayerTeamId(state, newval = 43643).should.be.fulfilled;
        result = await encodingState.getPrevPlayerTeamId(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encodingState.setLastSaleBlock(state, newval = 11223).should.be.fulfilled;
        result = await encodingState.getLastSaleBlock(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        result = await encodingState.getIsInTransitFromState(newState).should.be.fulfilled;
        result.should.be.equal(false);
        newState = await encodingState.setIsInTransit(newState, newval = true).should.be.fulfilled;
        result = await encodingState.getIsInTransitFromState(newState).should.be.fulfilled;
        result.should.be.equal(newval);
        newState = await encodingState.setIsInTransit(newState, newval = false).should.be.fulfilled;
        result = await encodingState.getIsInTransitFromState(newState).should.be.fulfilled;
        result.should.be.equal(newval);
    });


});