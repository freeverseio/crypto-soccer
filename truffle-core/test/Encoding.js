const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const Encoding = artifacts.require('Encoding');

contract('Encoding', (accounts) => {

    beforeEach(async () => {
        encoding = await Encoding.new().should.be.fulfilled;
    });
    
    it('encodeTactics', async () =>  {
        PLAYERS_PER_TEAM_MAX = await encoding.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        lineup = Array.from(new Array(11), (x,i) => i);
        encoded = await encoding.encodeTactics(lineup, tacticsId = 2).should.be.fulfilled;
        decoded = await encoding.decodeTactics(encoded).should.be.fulfilled;
        let {0: line, 1: tact} = decoded;
        tact.toNumber().should.be.equal(tacticsId);
        for (p = 0; p < 11; p++) {
            line[p].toNumber().should.be.equal(lineup[p]);
        }
        // try to provide a lineup beyond range
        lineupWrong = lineup;
        lineupWrong[4] = PLAYERS_PER_TEAM_MAX;
        encoded = await encoding.encodeTactics(lineup, tacticsId = 2).should.be.rejected;
        // try to provide a tacticsId beyond range
        encoded = await encoding.encodeTactics(lineup, tacticsId = 64).should.be.rejected;
    });
    
    it('encoding of TZ and country in teamId and playerId', async () =>  {
        encoded = await encoding.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 3, val = 4).should.be.fulfilled;
        decoded = await encoding.decodeTZCountryAndVal(encoded).should.be.fulfilled;
        const {0: timeZone, 1: country, 2: value} = decoded;
        timeZone.toNumber().should.be.equal(tz);
        country.toNumber().should.be.equal(countryIdxInTZ);
        value.toNumber().should.be.equal(val);
    });

    it('encoding and decoding skills', async () => {
        const sk = [16383, 13, 4, 56, 456]
        const skills = await encoding.encodePlayerSkills(
            sk,
            monthOfBirth = 4, 
            playerId = 143,
            potential = 5,
            forwardness = 3,
            leftishness = 4
        ).should.be.fulfilled;
        result = await encoding.getMonthOfBirth(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(monthOfBirth);
        result = await encoding.getPotential(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(potential);
        result = await encoding.getForwardness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(forwardness);
        result = await encoding.getLeftishness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(leftishness);
        result = await encoding.getPlayerIdFromSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(playerId);
        result = await encoding.getSkillsVec(skills).should.be.fulfilled;
        for (s=0; s < sk.length; s++) {
            result[s].toNumber().should.be.equal(sk[s]);
        }
    });

    it('encoding skills with wrong forwardness and leftishness', async () =>  {
        sk = [16383, 13, 4, 56, 456];
        monthOfBirth = 4;
        playerId = 143;
        potential = 5;
        // leftishness = 0 only possible for goalkeepers:
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, potential, forwardness = 0, leftishness = 0).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, potential, forwardness = 1, leftishness = 0).should.be.rejected;
        // forwardness is 5 at max:
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, potential, forwardness = 5, leftishness = 1).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, potential, forwardness = 6, leftishness = 1).should.be.rejected;
        // leftishness is 7 at max:
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, potential, forwardness = 5, leftishness = 7).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, potential, forwardness = 5, leftishness = 8).should.be.rejected;
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