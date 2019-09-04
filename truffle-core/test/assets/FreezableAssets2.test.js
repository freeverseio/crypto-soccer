const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Assets = artifacts.require('FreezableAssets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('FreezableAssets2', (accounts) => {
    let assets = null;
    let playerStateLib = null;
    const ALICE = accounts[1];
    const BOB = accounts[2];

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
    });

    it('change ownership', async () => {
        await assets.createTeam("Barca", ALICE).should.be.fulfilled;
        await assets.createTeam("Madrid", BOB).should.be.fulfilled;
        const player = 1;
        const originOwner = await assets.getPlayerOwner(player).should.be.fulfilled;
        originOwner.should.be.equal(ALICE);
        const toTeam = 2;
        await assets.transferPlayer(player, toTeam).should.be.fulfilled;
        const targetOwner = await assets.getPlayerOwner(player).should.be.fulfilled;
        targetOwner.should.be.equal(BOB);
    });
});