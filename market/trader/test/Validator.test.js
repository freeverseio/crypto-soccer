const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Accounts = require('web3-eth-accounts');

describe('Validator', () => {
    it('create signature of "ciao"', async () => {
        var accounts = new Accounts();
        const account = await accounts.create("iamaseller");
        const hashedMsg = account.sign("ciao");
        const address = accounts.recover(hashedMsg);
        address.should.be.equal(account.address);
    });
}) 