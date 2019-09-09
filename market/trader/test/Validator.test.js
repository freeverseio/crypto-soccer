const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Accounts = require('web3-eth-accounts');
const Validator = require('../src/Validator');

describe('Validator', () => {
    const accounts = new Accounts();
    const privateKey = '0x0a69684608770d018143dd70dc5dc5b6beadc366b87e45fcb567fc09407e7fe5';
    const account = accounts.privateKeyToAccount(privateKey);

    it('get the signer address by signatureObject', () => {
        const msg = "ciao";
        const signatureObject = account.sign(msg);
        const address = accounts.recover(signatureObject);
        address.should.be.equal(account.address);
        signatureObject.signature.should.be.equal('0xcf0a59da3b50f2827d9b15fc83391cd5feaf9b25131c2f4f20e7ae2d4fba811b41f35b6b17ba566c38a5c3737a759018be1f9064b7c8f56daaf4c00e51c7df281b');
    });

    it('get the signer address by message and signature', () => {
        const msg = "ciao";
        const signatureObject = account.sign(msg);
        address = accounts.recover(signatureObject.messageHash, signatureObject.signature, true);
        address.should.be.equal(account.address);
        signatureObject.signature.should.be.equal('0xcf0a59da3b50f2827d9b15fc83391cd5feaf9b25131c2f4f20e7ae2d4fba811b41f35b6b17ba566c38a5c3737a759018be1f9064b7c8f56daaf4c00e51c7df281b');
    });

    it('get the signer address by message and v,r,s', () => {
        const accounts = new Accounts();
        const account = accounts.privateKeyToAccount(privateKey);
        const msg = "ciao";
        const hashedMsg = account.sign(msg);
        address = accounts.recover(hashedMsg.messageHash, hashedMsg.v, hashedMsg.r, hashedMsg.s, true);
        address.should.be.equal(account.address);
        hashedMsg.signature.should.be.equal('0xcf0a59da3b50f2827d9b15fc83391cd5feaf9b25131c2f4f20e7ae2d4fba811b41f35b6b17ba566c38a5c3737a759018be1f9064b7c8f56daaf4c00e51c7df281b');
    });

    it('get signer address', () => {
        const accounts = new Accounts();
        const account = accounts.privateKeyToAccount(privateKey);
        const msg = "ciao";
        const hashedMsg = account.sign(msg);
        const validator = new Validator();
        const address = validator.recoverSignerAddress(msg, hashedMsg.signature);
        address.should.be.equal(account.address);
    });
}) 