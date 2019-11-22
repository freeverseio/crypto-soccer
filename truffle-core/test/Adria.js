const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');

contract('Engine', (accounts) => {
    
    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        await engine.setPreCompAddr(precomp.address).should.be.fulfilled;
    });
    
    it('play one half of a match', async () => {
        // inputs:
        nowInSecs  = 1570147200
        seed = '0xb0ae22e2f60d41a9c23f77adac5bfdccb8228e51dbd13aa2a3654c5276b2c04a'  // = web3.utils.toBN(web3.utils.keccak256("32123"));
        teamState = Array.from(new Array(25), (x,i) => '0xa000998000020896142b600010001000100010001'); 
        tactics = '0x5cc299ac5a928398a4188200000'
        firstHalfLog = [0, 0]
        matchBooleans = [ false, false, false]
        
        // calling the view function in Engine.sol:
        result = await engine.playHalfMatch(seed, nowInSecs, [teamState, teamState], [tactics, tactics], firstHalfLog, matchBooleans).should.be.fulfilled;
    });

});