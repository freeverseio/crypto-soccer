const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;
    let TEAMS_PER_LEAGUE = null;
    let MATCHDAYS = null;
    let MATCHES_PER_DAY = null;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        encoding = leagues;
        await leagues.init().should.be.fulfilled;
        TEAMS_PER_LEAGUE = await leagues.TEAMS_PER_LEAGUE().should.be.fulfilled;
        MATCHDAYS = await leagues.MATCHDAYS().should.be.fulfilled;
        MATCHES_PER_DAY = await leagues.MATCHES_PER_DAY().should.be.fulfilled;
        });

    it('check initial constants', async () =>  {
        MATCHDAYS.toNumber().should.be.equal(14);
        MATCHES_PER_DAY.toNumber().should.be.equal(4);
        TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    });

    it('get teams for match in wrong day', async () => {
        matchIdxInDay = 0; 
        day = MATCHDAYS-1; 
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        day = MATCHDAYS; 
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    });

    it('get teams for match in wrong match in day', async () => {
        day = 0;
        matchIdxInDay = MATCHES_PER_DAY-1;
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        matchIdxInDay = MATCHES_PER_DAY;
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        day = 0;
        matchIdxInDay = 0;
        teams = await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        day = Math.floor(MATCHDAYS/2);
        teams = await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });

});