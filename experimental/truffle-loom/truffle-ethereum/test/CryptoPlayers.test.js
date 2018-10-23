require('chai')
    .use(require('chai-as-promised'))
    .should();

const Gateway = artifacts.require('Gateway');
const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    it('correct deployed', async () => {
        const gateway = await Gateway.deployed();
        gateway.should.not.equal(null);
        const player = await CryptoPlayers.new(gateway.address);
        player.should.not.equal(null);
    });
});
