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
    });

    it('append player state 2d', async () => {
        const state2D = await instance.leagueStateCreate().should.be.fulfilled;
        let state3D = await instance.playerState3DCreate().should.be.fulfilled;
        let count = await instance.playerState3DSize(state3D).should.be.fulfilled;
        count.toString().should.be.equal('0');
        state3D = await instance.playerState3DAppend(state3D, state2D).should.be.fulfilled;
        count = await instance.playerState3DSize(state3D).should.be.fulfilled;
        count.toString().should.be.equal('1');
        state3D = await instance.playerState3DAppend(state3D, state2D).should.be.fulfilled;
        count = await instance.playerState3DSize(state3D).should.be.fulfilled;
        count.toString().should.be.equal('2');
    });
}) 