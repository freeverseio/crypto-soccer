require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('Storage');

contract('Storage', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Storage.new().should.be.fulfilled;
    });

    it('no initial players', async () =>{
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('no initial teams', async () =>{
        const count = await contract.getNCreatedTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    })

    it('add player', async () => {
        const name = "player";
        const state = 34324;
        const teamId = 1;
        await contract.addPlayer(name, state, teamId).should.be.fulfilled;
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        const nameResult = await contract.getPlayerName(count-1);
        nameResult.should.be.equal(name);
        const stateResult = await contract.getPlayerState(count-1);
        stateResult.toNumber().should.be.equal(state);
    });

    it ('create team', async () => {
        await contract.addTeam("team", accounts[0]).should.be.fulfilled;
        const name  = await contract.getTeamName(1).should.be.fulfilled;
        name.should.be.equal("team");
    })

    it('team name', async() => {
        const team = "team";
        await contract.getTeamName(0).should.be.rejected;
        await contract.getTeamName(1).should.be.rejected;
        await contract.addTeam(team, accounts[0]).should.be.fulfilled;
        const name = await contract.getTeamName(1).should.be.fulfilled;
        name.should.be.equal(team);
    });

    it('team owner', async () => {
        let owner = await contract.teamOwnerOf(0).should.be.rejected;
        await contract.addTeam("team", accounts[0]).should.be.fulfilled;
        owner = await contract.teamOwnerOf(1).should.be.fulfilled;
        owner.should.be.equal(accounts[0]);
    })

    it('team name of unexistent player', async () => {
        await contract.teamNameByPlayer("unexistent").should.be.rejected;
    });

    it('team name by player', async () => {
        const team = "team";
        const player = "player";
        const playerState = 44535;
        await contract.addTeam(team, accounts[0]);
        await contract.addPlayer(player, playerState, 1);
        const index = await contract.getTeamIndexByPlayer(player).should.be.fulfilled;
        index.toNumber().should.be.equal(1);
        const name = await contract.teamNameByPlayer(player).should.be.fulfilled;
        name.should.be.equal(team);
    });
});
