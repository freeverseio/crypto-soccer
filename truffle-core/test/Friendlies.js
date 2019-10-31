const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const Assets = artifacts.require('Assets');

contract('Assets', (accounts) => {
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const N_SKILLS = 5;
    let initTx = null;

    const it2 = async(text, f) => {};

    beforeEach(async () => {
        assets = await Assets.new().should.be.fulfilled;
        initTx = await assets.init().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = await assets.PLAYERS_PER_TEAM_INIT().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await assets.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        LEAGUES_PER_DIV = await assets.LEAGUES_PER_DIV().should.be.fulfilled;
        TEAMS_PER_LEAGUE = await assets.TEAMS_PER_LEAGUE().should.be.fulfilled;
        FREE_PLAYER_ID = await assets.FREE_PLAYER_ID().should.be.fulfilled;
        NULL_ADDR = await assets.NULL_ADDR().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = PLAYERS_PER_TEAM_INIT.toNumber();
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        LEAGUES_PER_DIV = LEAGUES_PER_DIV.toNumber();
        TEAMS_PER_LEAGUE = TEAMS_PER_LEAGUE.toNumber();
        });

        
    it('create special players', async () => {
    });
});