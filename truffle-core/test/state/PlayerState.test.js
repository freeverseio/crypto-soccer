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

    // it('encoding and decoding skills', async () => {
    //     const sk = [16383, 13, 4, 56, 456]
    //     const skills = await playerStateLib.encodePlayerSkills(
    //         defence = sk[0],
    //         speed = sk[1],
    //         pass = sk[2],
    //         shoot = sk[3],
    //         endurance = sk[4],
    //         monthOfBirth = 4, 
    //         playerId = 143,
    //     ).should.be.fulfilled;
    //     result = await playerStateLib.getDefence(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(defence);
    //     result = await playerStateLib.getSpeed(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(speed);
    //     result = await playerStateLib.getPass(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(pass);
    //     result = await playerStateLib.getShoot(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(shoot);
    //     result = await playerStateLib.getEndurance(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(endurance);
    //     result = await playerStateLib.getMonthOfBirthInUnixTime(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(monthOfBirth);
    //     result = await playerStateLib.getPlayerIdFromSkills(skills).should.be.fulfilled;
    //     result.toNumber().should.be.equal(playerId);
    //     result = await playerStateLib.getSkillsVec(skills).should.be.fulfilled;
    //     for (s=0; s < sk.length; s++) {
    //         result[s].toNumber().should.be.equal(sk[s]);
    //     }
    // });

    
    it('encode decode player state', async () => {
        const playerId = 231;
        const currentTeamId = 432432;
        const currentShirtNum = 12;
        const prevPlayerTeamId = 32123;
        const lastSaleBlock = 3221;
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
    });
return;
// test individual changes when stuff is already full
    it('skills getters from state player', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        const playerState = await playerStateLib.playerStateCreate(
            defence,
            speed,
            pass,
            shoot,
            endurance,
            monthOfBirthInUnixTime = 40,
            playerId = 41,
            currentTeamId = 42,
            currentShirtNum = 3,
            prevLeagueId = 44,
            prevTeamPosInLeague = 45,
            prevShirtNumInLeague = 6,
            lastSaleBlock = 47).should.be.fulfilled;
        let result = await playerStateLib.getDefence(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(defence);
        result = await playerStateLib.getSpeed(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(speed);
        result = await playerStateLib.getPass(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(pass);
        result = await playerStateLib.getShoot(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(shoot);
        result = await playerStateLib.getEndurance(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(endurance);
        result = await playerStateLib.getCurrentTeamId(playerState).should.be.fulfilled;
        result.should.be.bignumber.equal('42');
        result = await playerStateLib.getCurrentShirtNum(playerState).should.be.fulfilled;
        result.should.be.bignumber.equal('3');
    });

    it('player state evolve', async () => {
        const playerState = await playerStateLib.playerStateCreate(
            defence = 3, 
            speed = 4, 
            pass = 6, 
            shoot = 11, 
            endurance = 9, 
            0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        const delta = 3;
        const updatedState = await playerStateLib.playerStateEvolve(playerState, delta).should.be.fulfilled;
        updatedState.should.be.bignumber.that.equals(playerState);
        let skill = await playerStateLib.getDefence(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence); // + delta);
        skill = await playerStateLib.getSpeed(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed); // + delta);
        skill = await playerStateLib.getPass(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass); // + delta);
        skill = await playerStateLib.getShoot(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot); // + delta);
        skill = await playerStateLib.getEndurance(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance); // + delta);
    });

    it('get skills', async () => {
        const playerState = await playerStateLib.playerStateCreate(
            defence = 0, 
            speed = 0, 
            pass = 0, 
            shoot = 0, 
            endurance = 1, 
            0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        const skills = await playerStateLib.getSkills(playerState).should.be.fulfilled;
        skills.toNumber().should.be.equal(1);
    });
});