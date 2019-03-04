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
});