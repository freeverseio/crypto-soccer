const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const debug = require('../utils/debugUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Encoding = artifacts.require('EncodingSkills');
const EncodingTact = artifacts.require('EncodingTacticsPart1');
const EncodingSet = artifacts.require('EncodingSkillsSetters');
const EncodingGet = artifacts.require('EncodingSkillsGetters');
const Utils = artifacts.require('Utils');
const Privileged = artifacts.require('Privileged');

function secsToDays(secs) {
    return secs/ (24 * 3600);
}


function dayOfBirthToAgeYears(dayOfBirth){ 
    const now = Math.floor(new Date()/1000);
    ageYears = (secsToDays(now) - dayOfBirth)*7/365;
    return ageYears;
}

contract('Encoding', (accounts) => {

    const it2 = async(text, f) => {};

    beforeEach(async () => {
        constants = await ConstantsGetters.new().should.be.fulfilled;
        encoding = await Encoding.new().should.be.fulfilled;
        utils = await Utils.new().should.be.fulfilled;
        privileged = await Privileged.new().should.be.fulfilled;
        encodingSet = await EncodingSet.new().should.be.fulfilled;
        encodingGet = await EncodingGet.new().should.be.fulfilled;
        encodingTact = await EncodingTact.new().should.be.fulfilled;
    });

    it('creating buyNow players: ageModifier', async () =>  {
        mods = [];
        for (age = 16; age < 38; age += 3) {
            mod = await privileged.ageModifier(age).should.be.fulfilled;
            mods.push(mod);
        }
        expectedMods = [ 10000, 9610, 9220, 8830, 8440, 8050, 6550, 5050 ];
        debug.compareArrays(mods, expectedMods, toNum = true, verbose = false);
    });
    
    it('creating buyNow players: potentialModifier', async () =>  {
        mods = [];
        for (pot = 0; pot < 10; pot++) {
            mod = await privileged.potentialModifier(pot).should.be.fulfilled;
            mods.push(mod);
        }
        expectedMods = [ 8500, 8833, 9166, 9500, 9833, 10166, 10500, 10833, 11166, 11500 ];
        debug.compareArrays(mods, expectedMods, toNum = true, verbose = false);
    });
    
    it('creating one buyNow player', async () =>  {
        expectedSkills = [ 1740, 1219, 979, 1226, 1903 ];
        expectedTraits = [0, 3, 6, 1];
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        var {0: skills, 1: ageYears, 2: traits, 3: internalId} = await privileged.createBuyNowPlayerIdPure(playerValue = 1000, seed, forwardPos = 3).should.be.fulfilled;
        // check that the average skill is as expected:
        expectedAvgSkill = await privileged.computeAvgSkills(playerValue, ageYears, traits[0]).should.be.fulfilled;
        sumSkills = expectedSkills.reduce((a, b) => a + b, 0);
        (Math.abs(expectedAvgSkill.toNumber() - sumSkills/5) < 20).should.be.equal(true);
        // compare actual values
        debug.compareArrays(skills, expectedSkills, toNum = true, verbose = false);
        ageYears.toNumber().should.be.equal(29);
        debug.compareArrays(traits, expectedTraits, toNum = true, verbose = false);
        internalId.should.be.bignumber.equal("1247534008908");
        
        // test that you get the same via the non-pure function:
        var {0: finalId, 1: skills2, 2: dayOfBirth, 3: traits2, 4: internalId2} = await privileged.createBuyNowPlayerId(playerValue = 1000, seed, forwardPos = 3).should.be.fulfilled;
        debug.compareArrays(skills2, expectedSkills, toNum = true, verbose = false);
        debug.compareArrays(traits2, expectedTraits, toNum = true, verbose = false);
        internalId2.should.be.bignumber.equal(internalId);

        const now = Math.floor(new Date()/1000);
        expectedDayOfBirth = Math.floor(secsToDays(now) - ageYears*365/7);
        (Math.abs(dayOfBirth.toNumber() - expectedDayOfBirth) < 10).should.be.equal(true);
        
    });
    
    it('creating buyNow players scales linearly with value, while other data remains the same', async () =>  {
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        var {0: skills, 1: ageYears, 2: traits, 3: internalId} = await privileged.createBuyNowPlayerIdPure(playerValue = 1000, seed, forwardPos = 3).should.be.fulfilled;
        var {0: skills2, 1: ageYears2, 2: traits2, 3: internalId2} = await privileged.createBuyNowPlayerIdPure(playerValue = 2000, seed, forwardPos = 3).should.be.fulfilled;
        for (s = 0; s < skills.length; s++) {
            (Math.abs(skills2[s].toNumber() - 2*skills[s].toNumber()) < 20).should.be.equal(true);
        }
        for (t = 0; t < traits.length; t++) {
            traits2[t].toNumber().should.be.equal(traits[t].toNumber());
        }
        internalId2.should.be.bignumber.equal(internalId);
        ageYears.should.be.bignumber.equal(ageYears2);
    });

    it('creating a batch of buyNow players', async () =>  {
        expectedSkills = [ 1740, 1219, 979, 1226, 1903 ];
        expectedTraits = [0, 3, 6, 1];
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        const nPlayersPerForwardPos = [0,0,0,2];
        var {0: playerIdArray, 1: skillsArray, 2: dayOfBirthArray, 3: traitsArray, 4: internalIdArray} = await privileged.createBuyNowPlayerIdBatch(
            playerValue = 1000, seed, nPlayersPerForwardPos
        ).should.be.fulfilled;

        // compare actual values
        debug.compareArrays(skillsArray[0], expectedSkills, toNum = true, verbose = false);
        debug.compareArrays(traitsArray[0], expectedTraits, toNum = true, verbose = false);
        internalIdArray[0].should.be.bignumber.equal("1247534008908");
        internalIdArray[1].should.not.be.bignumber.equal("1247534008908");
    });
    
    it('creating a batch of buyNow players and displaying', async () =>  {
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        const nPlayersPerForwardPos = [10,10,10,10];
        var {0: playerIdArray, 1: skillsArray, 2: dayOfBirthArray, 3: traitsArray, 4: internalIdArray} = await privileged.createBuyNowPlayerIdBatch(
            playerValue = 1000, seed, nPlayersPerForwardPos
        ).should.be.fulfilled;
        
        // traits: shoot, speed, pass, defence, endurance
        labels = ["GoalKeepers", "Defenders", "Midfielders", "Attackers"];
        st = "";
        counter = 0;
        for (pos = 0; pos < nPlayersPerForwardPos.length; pos++) {
            st += labels[pos];
            for (p = 0; p < nPlayersPerForwardPos[pos]; p++) {
                st += "\nPot: " + traitsArray[counter][0];
                st += " | Age: " + Math.floor(dayOfBirthToAgeYears(dayOfBirthArray[counter]));
                st += " | Shoot: " + skillsArray[counter][0];
                st += " | Speed: " + skillsArray[counter][1];
                st += " | Pass: " + skillsArray[counter][2];
                st += " | Defence: " + skillsArray[counter][3];
                st += " | Endurance: " + skillsArray[counter][4];
                counter++;
            }
            st += "\n"
        }
        console.log(st)
    });
    
    it('encodeTactics incorrect lineup', async () =>  {
        PLAYERS_PER_TEAM_MAX = await constants.get_PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        lineup = Array.from(new Array(14), (x,i) => i);
        substitutions = [4,10,2];
        subsRounds = [3,7,1];
        extraAttack = Array.from(new Array(10), (x,i) => (i%2 == 1 ? true: false));
        encoded = await encodingTact.encodeTactics(substitutions, subsRounds, lineup, extraAttack, tacticsId = 2).should.be.fulfilled;
        lineup[0] = PLAYERS_PER_TEAM_MAX;
        encoded = await encodingTact.encodeTactics(substitutions, subsRounds, lineup, extraAttack, tacticsId = 2).should.be.fulfilled;
    })
    
    it('encodeTactics', async () =>  {
        PLAYERS_PER_TEAM_MAX = await constants.get_PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        lineup = Array.from(new Array(14), (x,i) => i);
        substitutions = [4,10,2];
        subsRounds = [3,7,1];
        extraAttack = Array.from(new Array(10), (x,i) => (i%2 == 1 ? true: false));
        encoded = await encodingTact.encodeTactics(substitutions, subsRounds, lineup, extraAttack, tacticsId = 2).should.be.fulfilled;
        decoded = await encodingTact.decodeTactics(encoded).should.be.fulfilled;

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
        encoded = await encodingTact.encodeTactics(substitutions, subsRounds, lineup, extraAttack, tacticsId = 64).should.be.rejected;
        // try to provide a lineup beyond range
        lineupWrong = lineup;
        lineupWrong[4] = PLAYERS_PER_TEAM_MAX + 1;
        encoded = await encodingTact.encodeTactics(substitutions, subsRounds, lineupWrong, extraAttack, tacticsId = 2).should.be.rejected;
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
            alignedEndOfFirstHalf = true,
            redCardLastGame = true,
            gamesNonStopping = 2,
            injuryWeeksLeft = 6,
            substitutedFirstHalf = true,
            sumSkills
        ).should.be.fulfilled;

        skills.should.be.bignumber.equal('40439920000726868070503716865792521545121682176182486071370780491777');

        N_SKILLS = 5;
        resultSkills = [];
        for (s = 0; s < N_SKILLS; s++) {
            result = await encodingGet.getSkill(skills, s).should.be.fulfilled;
            resultSkills.push(result);
        }
        debug.compareArrays(resultSkills, sk, toNum = true, verbose = false);

        result = await encodingGet.getBirthDay(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(dayOfBirth);
        result = await encodingGet.getPotential(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(potential);
        result = await encodingGet.getForwardness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(forwardness);
        result = await encodingGet.getLeftishness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(leftishness);
        result = await encodingGet.getAggressiveness(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(aggressiveness);
        result = await encodingGet.getPlayerIdFromSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(playerId);
        result = await encodingGet.getAlignedEndOfFirstHalf(skills).should.be.fulfilled;
        result.should.be.equal(alignedEndOfFirstHalf);
        result = await encodingGet.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(redCardLastGame);

        result = await encodingGet.getGamesNonStopping(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(gamesNonStopping);
        gamesNonStopping = 7;        
        skills = await encodingSet.setGamesNonStopping(skills, gamesNonStopping).should.be.fulfilled;
        result = await encodingGet.getGamesNonStopping(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(gamesNonStopping);
        skillsdummy = await encodingSet.setGamesNonStopping(skills, 8).should.be.rejected;
        

        skills = await encodingSet.setPotential(skills, potential+1).should.be.fulfilled;
        result = await encodingGet.getPotential(skills).should.be.fulfilled;
        potential = potential + 1;
        result.toNumber().should.be.equal(potential);

        result = await encodingGet.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(injuryWeeksLeft);

        result = await encodingGet.getSubstitutedFirstHalf(skills).should.be.fulfilled;
        result.should.be.equal(substitutedFirstHalf);
        substitutedFirstHalf = !substitutedFirstHalf;
        skills = await encodingSet.setSubstitutedFirstHalf(skills, substitutedFirstHalf).should.be.fulfilled;
        result = await encodingGet.getSubstitutedFirstHalf(skills).should.be.fulfilled;
        result.should.be.equal(substitutedFirstHalf);

        result = await encodingGet.getSumOfSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sumSkills);
        result = await encodingGet.getGeneration(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(generation);
        
        result =  await encodingGet.getIsSpecial(skills).should.be.fulfilled;
        result.should.be.equal(false);
        skills2 = await encodingGet.addIsSpecial(skills).should.be.fulfilled;
        result =  await encodingGet.getIsSpecial(skills2).should.be.fulfilled;
        result.should.be.equal(true);
        
        sk = [2**16 - 43, 2**16 - 567, 0, 2**16 - 356, 2**16 - 4556]
        sumSkills = sk.reduce((a, b) => a + b, 0);
        for (s = 0; s < N_SKILLS; s++) {
            skills = await encodingSet.setSkill(skills, sk[s], s).should.be.fulfilled;
        }
        resultSkills = [];
        for (s = 0; s < N_SKILLS; s++) {
            result = await encodingGet.getSkill(skills, s).should.be.fulfilled;
            resultSkills.push(result);
        }
        debug.compareArrays(resultSkills, sk, toNum = true, verbose = false);

        alignedEndOfFirstHalf = !alignedEndOfFirstHalf;
        skills = await encodingSet.setAlignedEndOfFirstHalf(skills, alignedEndOfFirstHalf).should.be.fulfilled;
        result = await encodingGet.getAlignedEndOfFirstHalf(skills).should.be.fulfilled;
        result.should.be.equal(alignedEndOfFirstHalf);

        alignedEndOfFirstHalf = !alignedEndOfFirstHalf;
        skills = await encodingSet.setAlignedEndOfFirstHalf(skills, alignedEndOfFirstHalf).should.be.fulfilled;
        result = await encodingGet.getAlignedEndOfFirstHalf(skills).should.be.fulfilled;
        result.should.be.equal(alignedEndOfFirstHalf);
        
        redCardLastGame = !redCardLastGame;
        skills = await encodingSet.setRedCardLastGame(skills, redCardLastGame).should.be.fulfilled;
        result = await encodingGet.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(redCardLastGame);

        redCardLastGame = !redCardLastGame;
        skills = await encodingSet.setRedCardLastGame(skills, redCardLastGame).should.be.fulfilled;
        result = await encodingGet.getRedCardLastGame(skills).should.be.fulfilled;
        result.should.be.equal(redCardLastGame);
        
        injuryWeeksLeft -= 2;
        skills = await encodingSet.setInjuryWeeksLeft(skills, injuryWeeksLeft).should.be.fulfilled;
        result = await encodingGet.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(injuryWeeksLeft);

        injuryWeeksLeft += 1;
        skills = await encodingSet.setInjuryWeeksLeft(skills, injuryWeeksLeft).should.be.fulfilled;
        result = await encodingGet.getInjuryWeeksLeft(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(injuryWeeksLeft);

        skills = await encodingSet.setSumOfSkills(skills, sumSkills);
        result = await encodingGet.getSumOfSkills(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(sumSkills);
        
        generation += 2;
        skills = await encodingSet.setGeneration(skills, generation).should.be.fulfilled;
        result = await encodingGet.getGeneration(skills).should.be.fulfilled;
        result.toNumber().should.be.equal(generation);
        
        // testing full decode
        const {0: _skills, 1: _day, 2: _traits, 3: _playerId, 4: _alignedSubstRed, 5: _genNonstopInj} = await utils.fullDecodeSkills(skills).should.be.fulfilled     
        _day.toNumber().should.be.equal(dayOfBirth);
        _playerId.toNumber().should.be.equal(playerId);
        debug.compareArrays(_skills, sk, toNum = true, verbose = false);
        expectedTraits = [potential, forwardness, leftishness, aggressiveness];
        debug.compareArrays(_traits, expectedTraits, toNum = true, verbose = false);
        expectedBools = [alignedEndOfFirstHalf, substitutedFirstHalf, redCardLastGame];
        debug.compareArrays(_alignedSubstRed, expectedBools, toNum = false, verbose = false);
        expectedGenGameInj = [generation, gamesNonStopping, injuryWeeksLeft];
        debug.compareArrays(_genNonstopInj, expectedGenGameInj, toNum = true, verbose = false);
    });

    it('encoding skills with wrong forwardness and leftishness', async () =>  {
        sk = [16383, 13, 4, 56, 456];
        dayOfBirth = 4;
        generation = 2;
        playerId = 143;
        potential = 5;
        aggr = 2;
        alignedEndOfFirstHalf = true;
        redCardLastGame = true;
        gamesNonStopping = 2;
        injuryWeeksLeft = 6;
        substitutedFirstHalf = true;
        sumSkills = sk.reduce((a, b) => a + b, 0);
        // leftishness = 0 only possible for goalkeepers:
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 0, leftishness = 0, aggr],
            alignedEndOfFirstHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedFirstHalf, sumSkills).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 1, leftishness = 0, aggr],
            alignedEndOfFirstHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedFirstHalf, sumSkills).should.be.rejected;
        // forwardness is 5 at max:
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 5, leftishness = 1, aggr],
            alignedEndOfFirstHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedFirstHalf, sumSkills).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 6, leftishness = 1, aggr],
            alignedEndOfFirstHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedFirstHalf, sumSkills).should.be.rejected;
        // leftishness is 7 at max:
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 5, leftishness = 7, aggr],
            alignedEndOfFirstHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedFirstHalf, sumSkills).should.be.fulfilled;
        await encoding.encodePlayerSkills(sk, dayOfBirth, generation, playerId, [potential, forwardness = 5, leftishness = 8, aggr],
            alignedEndOfFirstHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft, substitutedFirstHalf, sumSkills).should.be.rejected;
    });
    
});