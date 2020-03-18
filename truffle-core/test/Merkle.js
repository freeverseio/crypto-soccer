const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const delegateUtils = require('../utils/delegateCallUtils.js');

const Merkle = artifacts.require('Merkle');


contract('Assets', (accounts) => {
    
    const it2 = async(text, f) => {};
    function toBytes32(name) { return web3.utils.utf8ToHex(name); }

    beforeEach(async () => {
        merkle = await Merkle.new().should.be.fulfilled;
    });
        
    it('create and verify', async () => {
        2+2;
    });

});