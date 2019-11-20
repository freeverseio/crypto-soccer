const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

    const Encoding = artifacts.require('EncodingSkills');
    const EncodingSet = artifacts.require('EncodingSkillsSetters');

contract('Encoding', (accounts) => {

    const it2 = async(text, f) => {};

    beforeEach(async () => {
        encoding = await Encoding.new().should.be.fulfilled;
        encodingSet = await EncodingSet.new().should.be.fulfilled;
    });
    
    it('encodeTactics', async () =>  {
        PLAYERS_PER_TEAM_MAX = await encoding.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        lineup = Array.from(new Array(14), (x,i) => i);
        substitutions = [4,10,2];
        subsRounds = [3,7,1];
        extraAttack = Array.from(new Array(10), (x,i) => (i%2 == 1 ? true: false));
        encoded = await encoding.encodeTactics(substitutions, subsRounds, lineup, extraAttack, tacticsId = 2).should.be.fulfilled;
        decoded = await encoding.decodeTactics(encoded).should.be.fulfilled;

        let {0: subs, 1: roun, 2: line, 3: attk, 4: tact} = decoded;
        tact.toNumber().should.be.equal(tacticsId);
        for (p = 0; p < 14; p++) {
            line[p].toNumber().should.be.equal(lineup[p]);
        }
        for (p = 0; p < 10; p++) {
            attk[p].should.be.equal(extraAttack[p]);
        }
        for (p = 0; p < 3; p++) {
            subs[p].toNumber().should.be.equal(substitutions[p]);
            roun[p].toNumber().should.be.equal(subsRounds[p]);
        }
        // // try to provide a tacticsId beyond range
        encoded = await encoding.encodeTactics(substitutions, subsRounds, lineup, extraAttack, tacticsId = 64).should.be.rejected;
        // try to provide a lineup beyond range
        lineupWrong = lineup;
        lineupWrong[4] = PLAYERS_PER_TEAM_MAX;
        encoded = await encoding.encodeTactics(substitutions, subsRounds, lineupWrong, extraAttack, tacticsId = 2).should.be.rejected;
    });

   
    it('encoding and decoding skills', async () => {
        sk = [2**16 - 16383, 2**16 - 13, 2**16 - 4, 2**16 - 56, 2**16 - 456]
        sumSkills = sk.reduce((a, b) => a + b, 0);

        skills = await encoding.encodePlayerSkills(
            sk,
            dayOfBirth = 4*365, 
            generation = 3,
            playerId = 143,
            [potential = 5,
            forwardness = 3,
            leftishness = 4,
            aggressiveness = 1],
            alignedEndOfLastHalf = true,
            redCardLastGame = true,
            gamesNonStopping = 2,
            injuryWeeksLeft = 6,
            substitutedLastHalf = true,
            sumSkills
        ).should.be.fulfilled;
        result = await encoding.getShoot(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[0]);
        result = await encoding.getSpeed(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[1]);
        result = await encoding.getPass(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[2]);
        result = await encoding.getDefence(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[3]);
        result = await encoding.getEndurance(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[4]);
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
        result = await encoding.getSubstitutedLastHalf(skills).should.be.fulfilled;
        result.should.be.equal(substitutedLastHalf);
        result = await encoding.getSumOfSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sumSkills);
        result = await encoding.getGeneration(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(generation);
        
        result =  await encoding.getIsSpecial(skills).should.be.fulfilled;
        result.should.be.equal(false);
        skills2 = await encoding.addIsSpecial(skills).should.be.fulfilled;
        result =  await encoding.getIsSpecial(skills2).should.be.fulfilled;
        result.should.be.equal(true);
        
        sk = [2**16 - 43, 2**16 - 567, 0, 2**16 - 356, 2**16 - 4556]
        sumSkills = sk.reduce((a, b) => a + b, 0);
        skills = await encodingSet.setShoot(skills, sk[0]).should.be.fulfilled;
        skills = await encodingSet.setSpeed(skills, sk[1]).should.be.fulfilled;
        skills = await encodingSet.setPass(skills, sk[2]).should.be.fulfilled;
        skills = await encodingSet.setDefence(skills, sk[3]).should.be.fulfilled;
        skills = await encodingSet.setEndurance(skills, sk[4]).should.be.fulfilled;
        result = await encoding.getShoot(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[0]);
        result = await encoding.getSpeed(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[1]);
        result = await encoding.getPass(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[2]);
        result = await encoding.getDefence(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[3]);
        result = await encoding.getEndurance(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sk[4]);

        skills = await encodingSet.setAlignedEndOfLastHalf(skills, true).should.be.fulfilled;
        result = await encoding.getAlignedEndOfLastHalf(skills).should.be.fulfilled;
        result.should.be.equal(true);

        skills = await encodingSet.setAlignedEndOfLastHalf(skills, false).should.be.fulfilled;
        result = await encoding.getAlignedEndOfLastHalf(skills).should.be.fulfilled;
        result.should.be.equal(false);
        
        skills = await encodingSet.setRedCardLastGame(skills, true).should.be.fulfilled;
        result = await encoding.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(true);

        skills = await encodingSet.setRedCardLastGame(skills, false).should.be.fulfilled;
        result = await encoding.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(false);
        
        skills = await encodingSet.setInjuryWeeksLeft(skills, 3).should.be.fulfilled;
        result = await encoding.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(3);

        skills = await encodingSet.setInjuryWeeksLeft(skills, 4).should.be.fulfilled;
        result = await encoding.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(4);

        skills = await encodingSet.setSumOfSkills(skills, sumSkills);
        result = await encoding.getSumOfSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sumSkills);
        
        result = await encoding.getTargetTeamId(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        
        skills = await encoding.setTargetTeamId(skills, targetTeamId = 2**40).should.be.fulfilled;
        result = await encoding.getTargetTeamId(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(targetTeamId);
        
    });

    it('encoding skills with wrong forwardness and leftishness', async () =>  {
        sk = [16383, 13, 4, 56, 456];
        dayOfBirth = 4;
        generation = 2;
        playerId = 143;
        potential = 5;
        aggr = 2;
        alignedEndOfLastHalf = true;
        redCardLastGame = true;
        gamesNonStopping = 2;
        injuryWeeksLeft = 6;
        substitutedLastHalf = true;
        sumSkills = sk.reduce((a, b) => a + b, 0);
        // leftishness = 0 only possible for goalkeepers:
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 0, leftishness = 0, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedLastHalf, sumSkills).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 1, leftishness = 0, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedLastHalf, sumSkills).should.be.rejected;
        // forwardness is 5 at max:
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 5, leftishness = 1, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedLastHalf, sumSkills).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 6, leftishness = 1, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedLastHalf, sumSkills).should.be.rejected;
        // leftishness is 7 at max:
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 5, leftishness = 7, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedLastHalf, sumSkills).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 5, leftishness = 8, aggr],
            alignedEndOfLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedLastHalf, sumSkills).should.be.rejected;
    });
    
});