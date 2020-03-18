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

    it('get merkle root', async () => {
        leafs = Array.from(new Array(4), (x,i) => web3.utils.keccak256(i.toString()));
        root = await merkle.merkleRoot(leafs, nLevels = 2).should.be.fulfilled;
        root1 = await merkle.hash_node(leafs[0], leafs[1]).should.be.fulfilled;
        root2 = await merkle.hash_node(leafs[2], leafs[3]).should.be.fulfilled;
        myRoot = await merkle.hash_node(root1, root2).should.be.fulfilled;
        myRoot.should.be.equal(root)
    });


    it('create and verify', async () => {
        leafs = Array.from(new Array(4), (x,i) => web3.utils.keccak256(i.toString()));
        root = await merkle.merkleRoot(leafs, nLevels = 2).should.be.fulfilled;
        leafHash = leafs[1];
        proof2 = await merkle.hash_node(leafs[2], leafs[3]).should.be.fulfilled;
        ok = await merkle.verify(root, [leafs[0], proof2], leafHash, 1).should.be.fulfilled;
        ok.should.be.equal(true);
    });

});