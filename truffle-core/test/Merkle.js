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
        leafs = Array.from(new Array(4), (x,i) => web3.utils.keccak256(i.toString()));
        console.log(leafs);
        root = await merkle.merkleRoot(leafs, nLevels = 2).should.be.fulfilled;
        console.log(root);
    });

});