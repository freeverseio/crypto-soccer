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
    const nullHash = '0x0';
    function hash32(x, y) {
        return web3.utils.keccak256(web3.eth.abi.encodeParameters(['bytes32', 'bytes32'], [x,y]));
    }
    function toBytes32(name) { return web3.utils.utf8ToHex(name); }

    beforeEach(async () => {
        merkle = await Merkle.new().should.be.fulfilled;
    });

    it('compatibility of hash function', async () => {
        leafs = Array.from(new Array(2), (x,i) => web3.utils.keccak256(i.toString()));
        resultBC = await merkle.hash_node(leafs[0], leafs[1]).should.be.fulfilled;
        resultJS = hash32(leafs[0], leafs[1]);
        resultBC.should.be.equal(resultJS)
        resultBC = await merkle.hash_node(nullHash, leafs[1]).should.be.fulfilled;
        resultJS = hash32(nullHash, leafs[1]);
        resultBC.should.be.equal(resultJS)
        resultBC = await merkle.hash_node(leafs[0], nullHash).should.be.fulfilled;
        resultJS = hash32(leafs[0], nullHash);
        resultBC.should.be.equal(resultJS)
        resultBC = await merkle.hash_node(nullHash, nullHash).should.be.fulfilled;
        resultJS = hash32(nullHash, nullHash);
        resultBC.should.be.equal(resultJS)
    });

    it('get merkle root', async () => {
        leafs = Array.from(new Array(4), (x,i) => web3.utils.keccak256(i.toString()));
        root = await merkle.merkleRoot(leafs, nLevels = 2).should.be.fulfilled;
        root1 = await merkle.hash_node(leafs[0], leafs[1]).should.be.fulfilled;
        root2 = await  merkle.hash_node(leafs[2], leafs[3]).should.be.fulfilled;
        myRoot = await merkle.hash_node(root1, root2).should.be.fulfilled;
        myRoot.should.be.equal(root)
    });

    it('verify', async () => {
        leafs = Array.from(new Array(4), (x,i) => web3.utils.keccak256(i.toString()));
        root = await merkle.merkleRoot(leafs, nLevels = 2).should.be.fulfilled;
        leafPos = 1;
        proof2 = await merkle.hash_node(leafs[2], leafs[3]).should.be.fulfilled;
        proof = [leafs[0], proof2];
        ok = await merkle.verify(root, proof, leafs[leafPos], leafPos).should.be.fulfilled;
        ok.should.be.equal(true);
    });

    it('build proof', async () => {
        leafs = Array.from(new Array(4), (x,i) => web3.utils.keccak256(i.toString()));
        root = await merkle.merkleRoot(leafs, nLevels = 2).should.be.fulfilled;
        leafPos = 1;
        proof = await merkle.buildProof(leafPos, leafs, nLevels).should.be.fulfilled; 
        ok = await merkle.verify(root, proof, leafs[leafPos], leafPos).should.be.fulfilled;
        ok.should.be.equal(true);
    });

});