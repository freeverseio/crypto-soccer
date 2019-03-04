require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeagueState = artifacts.require('LeagueState');

contract('LeagueState', (accounts) => {
    let instance = null;

    // beforeEach(async () => {
    //     instance = await LeagueState.new().should.be.fulfilled;
    //     TEAMSTATEEND = await instance.TEAMSTATEEND().should.be.fulfilled;
    //     LEAGUESTATEDIVIDER = await instance.LEAGUESTATEDIVIDER().should.be.fulfilled;
    // });

    // it('count of empty state', async () => {
    //     const state = await instance.leagueStatePerDayCreate().should.be.fulfilled;
    //     const count = await instance.leagueStatePerDayCount(state).should.be.fulfilled;
    //     count.toNumber().should.be.equal(0);
    // });

    // it('append 1 empty leagueState', async () => {
    //     const dayState = await instance.dayStateCreate().should.be.fulfilled;
    //     let leagueStatePerDay = await instance.leagueStatePerDayCreate().should.be.fulfilled;
    //     leagueStatePerDay = await instance.leagueStatePerDayAppend(leagueStatePerDay, leagueState).should.be.fulfilled;
    //     const count = await instance.leagueStatePerDayCount(leagueStatePerDay).should.be.fulfilled;
    //     count.toNumber().should.be.equal(1);
    // });

    // it('append 2 empty leagueState', async () => {
    //     const leagueState = await instance.leagueStateCreate().should.be.fulfilled;
    //     let leagueStatePerDay = await instance.leagueStatePerDayCreate().should.be.fulfilled;
    //     leagueStatePerDay = await instance.leagueStatePerDayAppend(leagueStatePerDay, leagueState).should.be.fulfilled;
    //     leagueStatePerDay = await instance.leagueStatePerDayAppend(leagueStatePerDay, leagueState).should.be.fulfilled;
    //     const count = await instance.leagueStatePerDayCount(leagueStatePerDay).should.be.fulfilled;
    //     count.toNumber().should.be.equal(2);
    // });
});