const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Contract = artifacts.require('PlayerState3D');

contract('PlayerState3D', (accounts) => {
    beforeEach(async () => {
        instance = await Contract.new().should.be.fulfilled;
    });

    it('create player state 3D', async () => {
        const state = await instance.playerState3DCreate().should.be.fulfilled;
        const valid = await instance.isValidPlayerState3D(state).should.be.fulfilled;
        valid.should.be.equal(true);
    })
}) 