const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const timeTravel = require('../utils/TimeTravel.js');
const merkleUtils = require('../utils/merkleUtils.js');
const chllUtils = require('../utils/challengeUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');



contract('Updates', (accounts) => {
    const nullHash = web3.eth.abi.encodeParameter('bytes32','0x0');
    
    const it2 = async(text, f) => {};
    
    const moveToNextVerse = async (updates, extraSecs = 0) => {
        now = await updates.getNow().should.be.fulfilled;
        nextTime = await updates.getNextVerseTimestamp().should.be.fulfilled;
        await timeTravel.advanceTime(nextTime - now + extraSecs);
        await timeTravel.advanceBlock().should.be.fulfilled;
    };

    beforeEach(async () => {
        constants = await ConstantsGetters.new().should.be.fulfilled;
    });

    it('test that cannot initialize updates twice', async () =>  {
        var {0: activeTeams, 1: orgMap} = chllUtils.createUniverse(10);
        console.log(activeTeams)
        console.log(orgMap)
        var team = chllUtils.createLeague();
        console.log(team)
    });

});