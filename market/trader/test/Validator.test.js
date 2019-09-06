const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Accounts = require('web3-eth-accounts');

describe('Validator', () => {
    it('get the signer address by signatureObject', () => {
        var accounts = new Accounts();
        const account = accounts.create("iamaseller");
        const msg = "ciao";
        const hashedMsg = account.sign(msg);
        var address = accounts.recover(hashedMsg);
        address.should.be.equal(account.address);
    });

    it('get the signer address by message and signature', () => {
        var accounts = new Accounts();
        const account = accounts.create("iamaseller");
        const msg = "ciao";
        const hashedMsg = account.sign(msg);
        address = accounts.recover(msg, hashedMsg.signature);
        address.should.be.equal(account.address);
    });
}) 