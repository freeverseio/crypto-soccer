require('chai')
    .use(require('chai-as-promised'))
    .should();

const Assets = artifacts.require('Teams');

contract('Assets', (accounts) => {
    let assets = null;

    beforeEach(async () => {
        assets = await Assets.new().should.be.fulfilled;
    });

    it('initial number of team', async () => {
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('create team', async () => {
        const receipt = await assets.createTeam(name = "Barca").should.be.fulfilled;
        const count = await assets.countTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        let teamName = receipt.logs[0].args.teamName;
        teamName.should.be.equal("Barca");
        const teamId = receipt.logs[0].args.teamId.toNumber();
        teamId.should.be.equal(1);
        teamName = await assets.getTeamName(teamId).should.be.fulfilled;
        teamName.should.be.equal("Barca");
    });

    it('get name of invalid team', async () => {
        await assets.getTeamName(0).should.be.rejected;
    });

    it('get name of unexistent team', async () => {
        await assets.getTeamName(1).should.be.rejected;
    });

    it('get playersId from teamId and pos in team', async () => {
        await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.rejected;
        await assets.createTeam(name = "Barca").should.be.fulfilled;
        await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=11).should.be.rejected;
        let playerId = await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.fulfilled;
        playerId.toNumber().should.be.equal(1);
        playerId = await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=10).should.be.fulfilled;
        playerId.toNumber().should.be.equal(11);
    });
})