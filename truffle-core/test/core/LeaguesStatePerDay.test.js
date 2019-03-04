require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeaguesStatePerDay = artifacts.require('LeaguesStatePerDay');

contract('LeaguesStatePerDay', (accounts) => {
    let instance = null;
    let TEAMSTATEDIVIDER = null;
    let LEAGUESTATEDIVIDER = null;
    const initBlock = 1;
    const step = 1;
    const id = 0;
    const teamIds = [1, 2];

    beforeEach(async () => {
        instance = await LeaguesStatePerDay.new().should.be.fulfilled;
        TEAMSTATEDIVIDER = await instance.TEAMSTATEDIVIDER().should.be.fulfilled;
        LEAGUESTATEDIVIDER = await instance.LEAGUESTATEDIVIDER().should.be.fulfilled;
    });

    it('count of empty state', async () => {
        const state = await instance.leagueStatePerDayCreate().should.be.fulfilled;
        const count = await instance.leagueStatePerDayCount(state).should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('append 1 empty leagueState', async () => {
        const leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        let leagueStatePerDay = await instance.leagueStatePerDayCreate().should.be.fulfilled;
        leagueStatePerDay = await instance.leagueStatePerDayAppend(leagueStatePerDay, leagueState).should.be.fulfilled;
        const count = await instance.leagueStatePerDayCount(leagueStatePerDay).should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('append 2 empty leagueState', async () => {
        const leagueState = await instance.leagueStateCreate().should.be.fulfilled;
        let leagueStatePerDay = await instance.leagueStatePerDayCreate().should.be.fulfilled;
        leagueStatePerDay = await instance.leagueStatePerDayAppend(leagueStatePerDay, leagueState).should.be.fulfilled;
        leagueStatePerDay = await instance.leagueStatePerDayAppend(leagueStatePerDay, leagueState).should.be.fulfilled;
        const count = await instance.leagueStatePerDayCount(leagueStatePerDay).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });
});