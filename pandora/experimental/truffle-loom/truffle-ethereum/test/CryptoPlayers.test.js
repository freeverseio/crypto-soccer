require('chai')
    .use(require('chai-as-promised'))
    .should();

const Gateway = artifacts.require('Gateway');
const Players = artifacts.require('Players');

contract('Players', (accounts) => {
    it('correct deployed', async () => {
        const gateway = await Gateway.deployed();
        gateway.should.not.equal(null);
        const player = await Players.new(gateway.address);
        player.should.not.equal(null);
    });
});
