/*
 Tests for all functions in EncodingIDs.sol and contracts inherited by it
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingGetters = artifacts.require('EncodingMatchEventsGetters');
const EncodingSetters = artifacts.require('EncodingMatchEventsSetters');

const fs = require('fs');
const { assert } = require('chai');

// async function idWrapper(id) {
//     const {0: timeZone, 1: country, 2: value} = await encoding.decodeTZCountryAndVal(id).should.be.fulfilled;
//     const result = {
//         encodedId: id.toString(),
//         timezone: Number(timeZone), 
//         country: Number(country),
//         val: Number(value),
//     };
//     return result;
// }

contract('EncodingIDs', (accounts) => {

    beforeEach(async () => {
        encodingGet = await EncodingGetters.new().should.be.fulfilled;
        encodingSet = await EncodingSetters.new().should.be.fulfilled;
    });
   
    it('encoding of TZ and country in teamId and playerId', async () =>  {
        N_ROUNDS = 12;
        teamThatAttacks = Array.from(new Array(N_ROUNDS), (x,i) => i%2);
        shooter = Array.from(new Array(N_ROUNDS), (x,i) => 15-i%5);
        assister = Array.from(new Array(N_ROUNDS), (x,i) => 15-i%4);
        isGoal = Array.from(new Array(N_ROUNDS), (x,i) => i%2 == 0);
        managesToShoot = Array.from(new Array(N_ROUNDS), (x,i) => i%2 == 1);

        encoded = 0;
        for (r = 0; r < N_ROUNDS; r++) {
            encoded = await encodingSet.setTeamThatAttacks(encoded, r, teamThatAttacks[r]).should.be.fulfilled;
            encoded = await encodingSet.setManagesToShoot(encoded, r, managesToShoot[r]).should.be.fulfilled;
            encoded = await encodingSet.setIsGoal(encoded, r, isGoal[r]).should.be.fulfilled;
            encoded = await encodingSet.setShooter(encoded, r, shooter[r]).should.be.fulfilled;
            encoded = await encodingSet.setAssister(encoded, r, assister[r]).should.be.fulfilled;
        }
        for (r = 0; r < N_ROUNDS; r++) {
            assert.equal(await encodingGet.getTeamThatAttacksFromEvents(encoded, r), teamThatAttacks[r]);
            assert.equal(await encodingGet.getManagesToShootFromEvents(encoded, r), managesToShoot[r]);
            assert.equal(await encodingGet.getIsGoalFromEvents(encoded, r), isGoal[r]);
            assert.equal(await encodingGet.getShooterFromEvents(encoded, r), shooter[r]);
            assert.equal(await encodingGet.getAssisterFromEvents(encoded, r), assister[r]);
        }
        
    });

});