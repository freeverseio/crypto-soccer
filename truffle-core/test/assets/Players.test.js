const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Players = artifacts.require('PlayersMock');
const PlayerStateLib = artifacts.require('PlayerState');

contract('Players', (accounts) => {
    let players = null;
    let playerStateLib = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        players = await Players.new(playerStateLib.address).should.be.fulfilled;
    });

    
});
 