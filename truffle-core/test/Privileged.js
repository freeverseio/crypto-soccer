const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Privileged = artifacts.require('Privileged');

contract('Privileged', (accounts) => {
    let privileged = null;

    beforeEach(async () => {
        privileged = await Privileged.new().should.be.fulfilled;
    });

    it('create batch of world players', async () => {
        const playerValue = 3000;
        const seed = 4;
        const nPlayersPerForwardPos = [1, 2, 3, 4];
        const epochDays = Math.floor(1588668910 / (3600 * 24));
        await privileged.createBuyNowPlayerIdBatch(
            playerValue,
            seed,
            nPlayersPerForwardPos,
            epochDays,
        ).should.be.fulfilled;
    });
})