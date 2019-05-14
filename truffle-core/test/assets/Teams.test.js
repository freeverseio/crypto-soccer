const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
    
const Assets = artifacts.require('Teams');
const PlayerStateLib = artifacts.require('PlayerState');

contract('Assets', (accounts) => {
    let assets = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
    });

    
})