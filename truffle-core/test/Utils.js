/*
 Tests for all functions in Privileged.sol
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Privileged = artifacts.require('Privileged');
const Utils = artifacts.require('Utils');
const debug = require('../utils/debugUtils.js');
const { isBigNumber } = require('web3-utils');


function getLevel(skills) {
    var level = 0;
    for (sk of skills) {
        level += Math.ceil(sk.toNumber()/1000);
    }
    return level;
}

contract('Privileged', (accounts) => {
    let privileged = null;
    const epochInDays = 18387; // May 5th 2020
    const tz = 1;
    const countryIdxInTz = 1;

    const it2 = async(text, f) => {};

    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }
    
    function dayOfBirthToAgeYears(dayOfBirth, nowInDays){ 
        ageYears = (nowInDays - dayOfBirth)*14/365;
        return ageYears;
    }
    
    beforeEach(async () => {
        utils = await Utils.new().should.be.fulfilled;
    });

    it('creating a team works', async () =>  {
        const timeZone = 10;
        const countryIdxInTZ = 0;
        const teamIdxInTZ = 0;
        const deployTimeInUnixEpochSecs = 1728813174000; // Date and time (GMT): Sunday, 13 October 2024 09:52:54
        const divisionCreationRound = 0;
        const {teamId, playerIds, playerSkillsAtBirth} = await utils.createTeam(timeZone, countryIdxInTZ, teamIdxInTZ, deployTimeInUnixEpochSecs, divisionCreationRound);
        assert.equal(teamId.toString(), "2748779069440");
        assert.equal(playerIds.length, 18);
        assert.equal(playerSkillsAtBirth.length, 18);
        for (let i = 0; i < 18; i++) {
            const decodedSkills = await utils.fullDecodeSkills(playerSkillsAtBirth[i]);
            assert.equal(decodedSkills.playerId.toString(), playerIds[i].toString())
            console.log(decodedSkills.dayOfBirth.toString())
        }
    });

    
})