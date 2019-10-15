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
        lineup = Array.from(new Array(14), (x,i) => i);
        substitutions = [4,10,2];
        rounds = [3,7,1];
        extraAttack = Array.from(new Array(10), (x,i) => i%2);
        encoded = await encoding.encodeTactics(substitutions, rounds, lineup, extraAttack, tacticsId = 2).should.be.fulfilled;
        decoded = await encoding.decodeTactics(encoded).should.be.fulfilled;
        let {0: subs, 1: roun, 2: line, 3: attk, 4: tact} = decoded;
        tact.toNumber().should.be.equal(tacticsId);
        for (p = 0; p < 14; p++) {
            line[p].toNumber().should.be.equal(lineup[p]);
        }
        for (p = 0; p < 10; p++) {
            attk[p].should.be.equal(extraAttack[p] == 1 ? true : false);
        }
        for (p = 0; p < 3; p++) {
            subs[p].toNumber().should.be.equal(substitutions[p]);
            roun[p].toNumber().should.be.equal(rounds[p]);
        }
        // // try to provide a tacticsId beyond range
        encoded = await encoding.encodeTactics(substitutions, rounds, lineup, extraAttack, tacticsId = 64).should.be.rejected;
        // try to provide a lineup beyond range
        lineupWrong = lineup;
        lineupWrong[4] = PLAYERS_PER_TEAM_MAX;
        encoded = await encoding.encodeTactics(substitutions, rounds, lineupWrong, extraAttack, tacticsId = 2).should.be.rejected;
    });
    
    it('encoding and decoding skills', async () => {
        const sk = [16383, 13, 4, 56, 456]
        const skills = await encoding.encodePlayerSkills(
            sk,
            dayOfBirth = 4*365, 
            playerId = 143,
            [potential = 5,
            forwardness = 3,
            leftishness = 4,
            aggressiveness = 1],
            alignedEndOfLastHalf = true,
            redCardLastGame = true,
            gamesNonStopping = 2,
            injuryWeeksLeft = 6
        ).should.be.fulfilled;
        result = await encoding.getBirthDay(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(dayOfBirth);
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
        result = await encoding.getAlignedEndOfLastHalf(skills).should.be.fulfilled;
        result.should.be.equal(alignedEndOfLastHalf);
        result = await encoding.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(redCardLastGame);
        result = await encoding.getGamesNonStopping(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(gamesNonStopping);
        result = await encoding.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(injuryWeeksLeft);
    });

    it('encoding skills with wrong forwardness and leftishness', async () =>  {
        sk = [16383, 13, 4, 56, 456];
        dayOfBirth = 4;
        playerId = 143;
        potential = 5;
        aggr = 2;
        alignedEndOfLastHalf = true;
        redCardLastGame = true;
        gamesNonStopping = 2;
        injuryWeeksLeft = 6;
        // leftishness = 0 only possible for goalkeepers:
        await encoding.encodePlayerSkills(sk, dayOfBirth, playerId, [potential, forwardness = 0, leftishness = 0, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, playerId, [potential, forwardness = 1, leftishness = 0, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.rejected;
        // forwardness is 5 at max:
        await encoding.encodePlayerSkills(sk, dayOfBirth, playerId, [potential, forwardness = 5, leftishness = 1, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, playerId, [potential, forwardness = 6, leftishness = 1, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.rejected;
        // leftishness is 7 at max:
        await encoding.encodePlayerSkills(sk, dayOfBirth, playerId, [potential, forwardness = 5, leftishness = 7, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, playerId, [potential, forwardness = 5, leftishness = 8, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.rejected;
    });
    
});