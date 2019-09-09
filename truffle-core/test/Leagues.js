const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;
    let encoding = null;
    let PLAYERS_PER_TEAM_MAX = null;
    let PLAYERS_PER_TEAM_INIT = null;
    let LEAGUES_PER_DIV = null;
    let TEAMS_PER_LEAGUE = null;
    let FREE_PLAYER_ID = null;
    let NULL_ADDR = null;
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const N_SKILLS = 5;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        encoding = leagues;
        await leagues.init().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = await leagues.PLAYERS_PER_TEAM_INIT().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await leagues.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        LEAGUES_PER_DIV = await leagues.LEAGUES_PER_DIV().should.be.fulfilled;
        TEAMS_PER_LEAGUE = await leagues.TEAMS_PER_LEAGUE().should.be.fulfilled;
        FREE_PLAYER_ID = await leagues.FREE_PLAYER_ID().should.be.fulfilled;
        NULL_ADDR = await leagues.NULL_ADDR().should.be.fulfilled;
        PLAYERS_PER_TEAM_INIT = PLAYERS_PER_TEAM_INIT.toNumber();
        PLAYERS_PER_TEAM_MAX = PLAYERS_PER_TEAM_MAX.toNumber();
        LEAGUES_PER_DIV = LEAGUES_PER_DIV.toNumber();
        TEAMS_PER_LEAGUE = TEAMS_PER_LEAGUE.toNumber();
        });

    it('check initial and max number of players per team', async () =>  {
        PLAYERS_PER_TEAM_INIT.should.be.equal(18);
        PLAYERS_PER_TEAM_MAX.should.be.equal(25);
        LEAGUES_PER_DIV.should.be.equal(16);
        TEAMS_PER_LEAGUE.should.be.equal(8);
    });

});