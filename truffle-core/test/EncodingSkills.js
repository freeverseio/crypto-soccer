const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const Encoding = artifacts.require('EncodingSkills');

contract('Encoding', (accounts) => {

    beforeEach(async () => {
        encoding = await Encoding.new().should.be.fulfilled;
    });
    
    it('encodeTactics', async () =>  {
        PLAYERS_PER_TEAM_MAX = await encoding.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        lineup = Array.from(new Array(11), (x,i) => i);
        extraAttack = Array.from(new Array(10), (x,i) => i%2);
        encoded = await encoding.encodeTactics(lineup, extraAttack, tacticsId = 2).should.be.fulfilled;
        decoded = await encoding.decodeTactics(encoded).should.be.fulfilled;
        let {0: line, 1: attk, 2: tact} = decoded;
        tact.toNumber().should.be.equal(tacticsId);
        for (p = 0; p < 11; p++) {
            line[p].toNumber().should.be.equal(lineup[p]);
        }
        for (p = 0; p < 10; p++) {
            attk[p].should.be.equal(extraAttack[p] == 1 ? true : false);
        }
        // try to provide a lineup beyond range
        lineupWrong = lineup;
        lineupWrong[4] = PLAYERS_PER_TEAM_MAX;
        encoded = await encoding.encodeTactics(lineup, tacticsId = 2).should.be.rejected;
        // try to provide a tacticsId beyond range
        encoded = await encoding.encodeTactics(lineup, tacticsId = 64).should.be.rejected;
    });
    
    it('encoding and decoding skills', async () => {
        const sk = [16383, 13, 4, 56, 456]
        const skills = await encoding.encodePlayerSkills(
            sk,
            monthOfBirth = 4, 
            playerId = 143,
            [potential = 5,
            forwardness = 3,
            leftishness = 4,
            aggressiveness = 1],
            alignedLastHalf = true,
            redCardLastGame = true,
            gamesNonStopping = 2,
            injuryWeeksLeft = 6
        ).should.be.fulfilled;
        result = await encoding.getMonthOfBirth(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(monthOfBirth);
        result = await encoding.getPotential(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(potential);
        result = await encoding.getForwardness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(forwardness);
        result = await encoding.getLeftishness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(leftishness);
        result = await encoding.getAggressiveness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(aggressiveness);
        result = await encoding.getPlayerIdFromSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(playerId);
        result = await encoding.getAlignedLastHalf(skills).should.be.fulfilled;
        result.should.be.equal(alignedLastHalf);
        result = await encoding.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(redCardLastGame);
        result = await encoding.getGamesNonStopping(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(gamesNonStopping);
        result = await encoding.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(injuryWeeksLeft);
    });

    it('encoding skills with wrong forwardness and leftishness', async () =>  {
        sk = [16383, 13, 4, 56, 456];
        monthOfBirth = 4;
        playerId = 143;
        potential = 5;
        aggr = 2;
        alignedLastHalf = true;
        redCardLastGame = true;
        gamesNonStopping = 2;
        injuryWeeksLeft = 6;
        // leftishness = 0 only possible for goalkeepers:
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, [potential, forwardness = 0, leftishness = 0, aggr],
            alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, [potential, forwardness = 1, leftishness = 0, aggr],
            alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.rejected;
        // forwardness is 5 at max:
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, [potential, forwardness = 5, leftishness = 1, aggr],
            alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, [potential, forwardness = 6, leftishness = 1, aggr],
            alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.rejected;
        // leftishness is 7 at max:
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, [potential, forwardness = 5, leftishness = 7, aggr],
            alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, monthOfBirth, playerId, [potential, forwardness = 5, leftishness = 8, aggr],
            alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.rejected;
    });
    
});