const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const PlayerState = artifacts.require('PlayerState');

/// TODO: evaluate to extract the skills part
contract('PlayerState', (accounts) => {
    let playerStateLib = null;

    beforeEach(async () => {
        playerStateLib = await PlayerState.new().should.be.fulfilled;
    });
    
    it('encoding of TZ and country in teamId and playerId', async () =>  {
        encoded = await playerStateLib.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 3, val = 4).should.be.fulfilled;
        decoded = await playerStateLib.decodeTZCountryAndVal(encoded).should.be.fulfilled;
        const {0: timeZone, 1: country, 2: value} = decoded;
        timeZone.toNumber().should.be.equal(tz);
        country.toNumber().should.be.equal(countryIdxInTZ);
        value.toNumber().should.be.equal(val);
    });

    it('encoding and decoding skills', async () => {
        const sk = [16383, 13, 4, 56, 456]
        const skills = await playerStateLib.encodePlayerSkills(
            sk,
            monthOfBirth = 4, 
            playerId = 143,
        ).should.be.fulfilled;
        result = await playerStateLib.getMonthOfBirth(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(monthOfBirth);
        result = await playerStateLib.getPlayerIdFromSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(playerId);
        result = await playerStateLib.getSkillsVec(skills).should.be.fulfilled;
        for (s=0; s < sk.length; s++) {
            result[s].toNumber().should.be.equal(sk[s]);
        }
    });

    
    it('encode decode player state', async () => {
        const playerId = 231;
        const currentTeamId = 432432;
        const currentShirtNum = 12;
        const prevPlayerTeamId = 32123;
        const lastSaleBlock = 3221;
        // check the initial setting of a player state (from empty)
        const state = await playerStateLib.encodePlayerState(playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock).should.be.fulfilled;
        result = await playerStateLib.getPlayerIdFromState(state).should.be.fulfilled;
        result.toNumber().should.be.equal(playerId);
        result = await playerStateLib.getCurrentTeamId(state).should.be.fulfilled;
        result.toNumber().should.be.equal(currentTeamId);
        result = await playerStateLib.getCurrentShirtNum(state).should.be.fulfilled;
        result.toNumber().should.be.equal(currentShirtNum);
        result = await playerStateLib.getPrevPlayerTeamId(state).should.be.fulfilled;
        result.toNumber().should.be.equal(prevPlayerTeamId);
        result = await playerStateLib.getLastSaleBlock(state).should.be.fulfilled;
        result.toNumber().should.be.equal(lastSaleBlock);
        // check the individual changes (from non-empty)
        newState = await playerStateLib.setCurrentTeamId(state, newval = 43).should.be.fulfilled;
        result = await playerStateLib.getCurrentTeamId(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await playerStateLib.setCurrentShirtNum(state, newval = 2).should.be.fulfilled;
        result = await playerStateLib.getCurrentShirtNum(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await playerStateLib.setPrevPlayerTeamId(state, newval = 43643).should.be.fulfilled;
        result = await playerStateLib.getPrevPlayerTeamId(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
        newState = await playerStateLib.setLastSaleBlock(state, newval = 11223).should.be.fulfilled;
        result = await playerStateLib.getLastSaleBlock(newState).should.be.fulfilled;
        result.toNumber().should.be.equal(newval);
    });


});