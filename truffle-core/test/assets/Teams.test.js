const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
    
const Assets = artifacts.require('Teams');

contract('Assets', (accounts) => {
    let assets = null;

    beforeEach(async () => {
        assets = await Assets.new().should.be.fulfilled;
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

    it('get playersId from teamId and pos in team', async () => {
        await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.rejected;
        await assets.createTeam(name = "Barca").should.be.fulfilled;
        await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=11).should.be.rejected;
        let playerId = await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=0).should.be.fulfilled;
        playerId.toNumber().should.be.equal(1);
        playerId = await assets.getPlayerIdFromTeamIdAndPos(teamId = 1, posInTeam=10).should.be.fulfilled;
        playerId.toNumber().should.be.equal(11);
    });

    it('sign team to league', async () => {
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 0).should.be.rejected;
        await assets.createTeam(name = "Barca").should.be.fulfilled;
        await assets.signToLeague(teamId = 1, leagueId = 1, posInLeague = 3).should.be.fulfilled;
        const currentHistory = await assets.getTeamCurrentHistory(1).should.be.fulfilled;
        currentHistory.currentLeagueId.should.be.bignumber.equal('1');
        currentHistory.posInCurrentLeague.should.be.bignumber.equal('3');
        currentHistory.prevLeagueId.should.be.bignumber.equal('0');
        currentHistory.posInPrevLeague.should.be.bignumber.equal('0');
    })
})