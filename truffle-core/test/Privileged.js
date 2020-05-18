const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Privileged = artifacts.require('Privileged');
const debug = require('../utils/debugUtils.js');

contract('Privileged', (accounts) => {
    let privileged = null;
    const epochInDays = 18387; // May 5th 2020
    const tz = 1;
    const countryIdxInTz = 1;

    const it2 = async(text, f) => {};

    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }
    
    function dayOfBirthToAgeYears(dayOfBirth){ 
        const now = Math.floor(new Date()/1000);
        ageYears = (secsToDays(now) - dayOfBirth)*7/365;
        return ageYears;
    }
    
    beforeEach(async () => {
        privileged = await Privileged.new().should.be.fulfilled;
    });

    it('create batch of world players', async () => {
        const playerValue = 3000;
        const seed = 4;
        const nPlayersPerForwardPos = [1, 2, 3, 4];
        const epochDays = Math.floor(1588668910 / (3600 * 24));
        const timezone = 1;
        const countryIdxInTz = 0;
        const maxPotential = 9;
        await privileged.createBuyNowPlayerIdBatch(
            playerValue,
            maxPotential,
            seed,
            nPlayersPerForwardPos,
            epochDays,
            timezone,
            countryIdxInTz,
        ).should.be.fulfilled;
    });
    
    it('creating buyNow players: ageModifier', async () =>  {
        mods = [];
        for (age = 16; age < 38; age += 3) {
            mod = await privileged.ageModifier(age).should.be.fulfilled;
            mods.push(mod);
        }
        expectedMods = [ 11500, 11110, 10720, 10330, 9940, 9550, 8050, 6550 ];
        debug.compareArrays(mods, expectedMods, toNum = true);
    });
    
    it('creating buyNow players: potentialModifier', async () =>  {
        mods = [];
        for (pot = 0; pot < 10; pot++) {
            mod = await privileged.potentialModifier(pot).should.be.fulfilled;
            mods.push(mod);
        }
        expectedMods = [ 8500, 8833, 9166, 9500, 9833, 10166, 10500, 10833, 11166, 11500 ];
        debug.compareArrays(mods, expectedMods, toNum = true);
    });
    
    it('creating one buyNow player', async () =>  {
        expectedSkills = [ 1526, 963, 1080, 875, 1547 ];
        expectedTraits = [0, 3, 6, 2];
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        var {0: skills, 1: ageYears, 2: traits, 3: internalId} = await privileged.createBuyNowPlayerIdPure(playerValue = 1000, maxPot = 9, seed, forwardPos = 3, tz, countryIdxInTz).should.be.fulfilled;
        // compare actual values
        debug.compareArrays(skills, expectedSkills, toNum = true);
        ageYears.toNumber().should.be.equal(29);
        debug.compareArrays(traits, expectedTraits, toNum = true);
        internalId.should.be.bignumber.equal("275260863937");
        // check that the average skill is as expected:
        expectedAvgSkill = await privileged.computeAvgSkills(playerValue, ageYears, traits[0]).should.be.fulfilled;
        sumSkills = expectedSkills.reduce((a, b) => a + b, 0);
        (Math.abs(expectedAvgSkill.toNumber() - sumSkills/5) < 20).should.be.equal(true);
        
        // test that you get the same via the non-pure function:
        var {0: finalId, 1: skills2, 2: dayOfBirth, 3: traits2, 4: internalId2} = await privileged.createBuyNowPlayerId(playerValue = 1000, maxPot = 9, seed, forwardPos = 3, epochInDays, tz, countryIdxInTz).should.be.fulfilled;
        debug.compareArrays(skills2, expectedSkills, toNum = true);
        debug.compareArrays(traits2, expectedTraits, toNum = true);
        internalId2.should.be.bignumber.equal(internalId);
        now = epochInDays*24*3600;
        expectedDayOfBirth = Math.floor(secsToDays(now) - ageYears*365/14);
        console.log(dayOfBirth.toNumber(), expectedDayOfBirth);
        (Math.abs(dayOfBirth.toNumber() - expectedDayOfBirth) < 14).should.be.equal(true);
        
    });

    it('creating buyNow players scales linearly with value, while other data remains the same', async () =>  {
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        var {0: skills, 1: ageYears, 2: traits, 3: internalId} = await privileged.createBuyNowPlayerIdPure(playerValue = 1000, maxPot = 9, seed, forwardPos = 3, tz, countryIdxInTz).should.be.fulfilled;
        var {0: skills2, 1: ageYears2, 2: traits2, 3: internalId2} = await privileged.createBuyNowPlayerIdPure(playerValue = 2000, maxPot = 9, seed, forwardPos = 3, tz, countryIdxInTz).should.be.fulfilled;
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
        expectedSkills = [ 797, 1093, 1257, 737, 1017 ];
        expectedTraits = [ 3, 3, 3, 1 ];
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        const nPlayersPerForwardPos = [0,0,0,2];
        var {0: playerIdArray, 1: skillsArray, 2: dayOfBirthArray, 3: traitsArray, 4: internalIdArray} = await privileged.createBuyNowPlayerIdBatch(
            playerValue = 1000, maxPot = 9, seed, nPlayersPerForwardPos, epochInDays, tz, countryIdxInTz
        ).should.be.fulfilled;

        // compare actual values
        debug.compareArrays(skillsArray[0], expectedSkills, toNum = true);
        debug.compareArrays(traitsArray[0], expectedTraits, toNum = true);
        internalIdArray[0].should.be.bignumber.equal("275195391431");
        internalIdArray[1].should.not.be.bignumber.equal("275195391431");
      
        // testing that they are created with the expected country and tz:
        var {0: tz2, 1: countryIdxInTz2} = await privileged.getTZandCountryIdxFromPlayerId(playerIdArray[0]).should.be.fulfilled;
        tz2.toNumber().should.be.equal(tz);
        countryIdxInTz2.toNumber().should.be.equal(countryIdxInTz);
    });
    
    it('creating a batch of buyNow players with less maxPotential', async () =>  {
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        const nPlayersPerForwardPos = [0,0,0,2];
        var {0: playerIdArray, 1: skillsArray, 2: dayOfBirthArray, 3: traitsArray, 4: internalIdArray} = await privileged.createBuyNowPlayerIdBatch(
            playerValue = 1000, maxPot = 4, seed, nPlayersPerForwardPos, epochInDays, tz, countryIdxInTz
        ).should.be.fulfilled;
        
        for (const trait of traitsArray) {
            (trait[0].toNumber() <= maxPot).should.be.equal(true);
        }
    });
    
    it('creating a batch of buyNow players and displaying', async () =>  {
        const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
        const nPlayersPerForwardPos = [10,10,10,10];
        var {0: playerIdArray, 1: skillsArray, 2: dayOfBirthArray, 3: traitsArray, 4: internalIdArray} = await privileged.createBuyNowPlayerIdBatch(
            playerValue = 1000, maxPot = 9, seed, nPlayersPerForwardPos, epochInDays, tz, countryIdxInTz
        ).should.be.fulfilled;
        h = web3.utils.keccak256(JSON.stringify(skillsArray) + JSON.stringify(traitsArray));
        assert.equal(h, '0x3e31ab49397c9131c1c8f4fff3ab9ecb25176a71f9252d42c6d72cfcb8732fcd', "createBuyNowPlayerIdBatch not as expected");

        if (false) {
            // traits: shoot, speed, pass, defence, endurance
            labels = ["GoalKeepers", "Defenders", "Midfielders", "Attackers"];
            st = "";
            st2 = "";
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
            console.log(st);
        }
    });
    
})