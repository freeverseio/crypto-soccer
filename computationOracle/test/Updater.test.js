const BigNumber = web3.BigNumber;
require('chai')
    .use(require('chai-as-promised'))
    .should();

const Updater = artifacts.require('Updater');

contract('Updater', accounts => {
    it('correct deployed', async () => {
        const updater = await Updater.new();
        updater.should.not.equal(null);
    });
});