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
const ConstantsGetters = artifacts.require('ConstantsGetters');

async function stateWrapper(state) {
    var result = {
        encodedState: state.toString(),
        currentTeamId: 0, 
        currentShirtNum: 0,
        prevPlayerTeamId: 0,
        lastSaleBlocknum: 0,
        isInTransit: false
    };
    result.currentTeamId = Number(await encodingState.getCurrentTeamIdFromPlayerState(state).should.be.fulfilled);
    result.currentShirtNum = Number(await encodingState.getCurrentShirtNum(state).should.be.fulfilled);
    result.prevPlayerTeamId = Number(await encodingState.getPrevPlayerTeamId(state).should.be.fulfilled);
    result.lastSaleBlocknum = Number(await encodingState.getLastSaleBlock(state).should.be.fulfilled);
    result.isInTransit = await encodingState.getIsInTransitFromState(state).should.be.fulfilled;

    return result;
}

contract('EncodingState', (accounts) => {

    beforeEach(async () => {
        encodingState = await EncodingState.new().should.be.fulfilled;
        encodingIDs = await EncodingIDs.new().should.be.fulfilled;
        constants = await ConstantsGetters.new().should.be.fulfilled;
    });
    
    it('encode decode player state', async () => {
        const writeMode = true;
        toWrite = [];

        const currentTeamId = await encodingIDs.encodeTZCountryAndVal(tz = 1, countryIdx = 0, teamIDx = 0).should.be.fulfilled;
        const currentShirtNum = 12;
        const prevPlayerTeamId = await encodingIDs.encodeTZCountryAndVal(tz = 1, countryIdx = 0, teamIDx = 1).should.be.fulfilled;
        const lastSaleBlock = 3221;
        // check the initial setting of a player state (from empty)
        const state = await encodingState.encodePlayerState(currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await stateWrapper(state))}

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
        if (writeMode) { await toWrite.push(await stateWrapper(newState))}
        result = await encodingState.getCurrentTeamIdFromPlayerState(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encodingState.setCurrentShirtNum(state, newval = 2).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await stateWrapper(newState))}
        result = await encodingState.getCurrentShirtNum(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encodingState.setPrevPlayerTeamId(state, newval = 43643).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await stateWrapper(newState))}
        result = await encodingState.getPrevPlayerTeamId(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await encodingState.setLastSaleBlock(state, newval = 11223).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await stateWrapper(newState))}
        result = await encodingState.getLastSaleBlock(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        result = await encodingState.getIsInTransitFromState(newState).should.be.fulfilled;
        result.should.be.equal(false);
        IN_TRANSIT_SHIRTNUM = await constants.get_IN_TRANSIT_SHIRTNUM().should.be.fulfilled;
        newState = await encodingState.setCurrentShirtNum(newState, newval = IN_TRANSIT_SHIRTNUM.toNumber()).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await stateWrapper(newState))}
        result = await encodingState.getIsInTransitFromState(newState).should.be.fulfilled;
        result.should.be.equal(true);
        newState = await encodingState.setCurrentShirtNum(newState, newval = 13).should.be.fulfilled;
        if (writeMode) { await toWrite.push(await stateWrapper(newState))}
        result = await encodingState.getIsInTransitFromState(newState).should.be.fulfilled;
        result.should.be.equal(false);

        const fs = require('fs');
        if (writeMode) {
            fs.writeFileSync('test/testdata/encodingStateTestData.json', JSON.stringify(toWrite), function(err) {
                if (err) {
                    console.log(err);
                }
            });
        }             
        
        writtenData = fs.readFileSync('test/testdata/encodingStateTestData.json', 'utf8');
        assert.equal(
            web3.utils.keccak256(writtenData),
            "0xf42362d5f2917296bb8c0d6b15caebd08af89a3946fa69fd9e60190709bea73e",
            "written testdata for encodingState State does not match expected result"
        );
    });


});