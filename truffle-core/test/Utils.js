/*
 Tests for all functions in Privileged.sol
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Utils = artifacts.require('Utils');
const debug = require('../utils/debugUtils.js');
const { isBigNumber } = require('web3-utils');
const { assert } = require('chai');


function getLevel(skills) {
    var level = 0;
    for (sk of skills) {
        level += Math.ceil(sk.toNumber()/1000);
    }
    return level;
}

contract('Utils', (accounts) => {
    let utils = null;
    const N_PLAYERS_IN_TEAM = 18;

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

    it('creating a team in (10,0,0) works', async () =>  {
        const timeZone = 10;
        const countryIdxInTZ = 0;
        const teamIdxInTZ = 0;
        const deployTimeInUnixEpochSecs = 1728813174000; // Date and time (GMT): Sunday, 13 October 2024 09:52:54
        const divisionCreationRound = 0;
        const {teamId, playerIds, playerSkillsAtBirth} = await utils.createTeam(timeZone, countryIdxInTZ, teamIdxInTZ, deployTimeInUnixEpochSecs, divisionCreationRound);
        assert.equal(teamId.toString(), "2748779069440");
        assert.equal(playerIds.length, N_PLAYERS_IN_TEAM);
        assert.equal(playerSkillsAtBirth.length, N_PLAYERS_IN_TEAM);
        const teamDecodedSkills = await utils.fullDecodeSkillsForEntireTeam(playerSkillsAtBirth);

        for (let i = 0; i < N_PLAYERS_IN_TEAM; i++) {
            assert.equal(teamDecodedSkills.skills[i].length, 5);
            const decodedSkills = await utils.fullDecodeSkills(playerSkillsAtBirth[i]);
            debug.compareArrays(teamDecodedSkills, decodedSkills, toNum = true);
            assert.equal(decodedSkills.playerId.toString(), playerIds[i].toString())
            const decodedPlayerId = await utils.decodeTZCountryAndVal(playerIds[i]);
            assert.equal(decodedPlayerId[0].toString(), timeZone.toString())
            assert.equal(decodedPlayerId[1].toString(), countryIdxInTZ.toString())
            assert.equal(decodedPlayerId[2].toString(), i.toString())
        }
    });

    it('creating a team in (23,3,5) works', async () =>  {
        const timeZone = 23;
        const countryIdxInTZ = 3;
        const teamIdxInTZ = 5;
        const deployTimeInUnixEpochSecs = 1728813174000; // Date and time (GMT): Sunday, 13 October 2024 09:52:54
        const divisionCreationRound = 0;
        const {teamId, playerIds, playerSkillsAtBirth} = await utils.createTeam(timeZone, countryIdxInTZ, teamIdxInTZ, deployTimeInUnixEpochSecs, divisionCreationRound);
        assert.equal(teamId.toString(), "6322997166085");
        assert.equal(playerIds.length, N_PLAYERS_IN_TEAM);
        assert.equal(playerSkillsAtBirth.length, N_PLAYERS_IN_TEAM);
        const teamDecodedSkills = await utils.fullDecodeSkillsForEntireTeam(playerSkillsAtBirth);
        for (let i = 0; i < N_PLAYERS_IN_TEAM; i++) {
            const decodedSkills = await utils.fullDecodeSkills(playerSkillsAtBirth[i]);
            debug.compareArrays(teamDecodedSkills, decodedSkills, toNum = true);
            assert.equal(decodedSkills.playerId.toString(), playerIds[i].toString())
            const decodedPlayerId = await utils.decodeTZCountryAndVal(playerIds[i]);
            assert.equal(decodedPlayerId[0].toString(), timeZone.toString())
            assert.equal(decodedPlayerId[1].toString(), countryIdxInTZ.toString())
            assert.equal(decodedPlayerId[2].toString(), (i + 5 * N_PLAYERS_IN_TEAM).toString())
        }
    });

    
})