require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('Storage');

contract('CryptoTeams', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Storage.new().should.be.fulfilled;
    });

    it('no initial players', async () =>{
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('no initial teams', async () =>{
        const count = await contract.getNCreatedTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    })

    it('add player', async () => {
        const name = "player";
        const state = 34324;
        await contract.addPlayer(name, state).should.be.fulfilled;
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(2);
        const nameResult = await contract.getPlayerName(count-1);
        nameResult.should.be.equal(name);
        const stateResult = await contract.getPlayerState(count-1);
        stateResult.toNumber().should.be.equal(state);
    });

    it('team name by player', async () => {
        const team = "team";
        const player = "player";
        const playerState = 44535;
        let name = await contract.teamNameByPlayer(player).should.be.fulfilled;
        name.should.be.equal("");
        await contract.addTeam(team);
        await contract.addPlayer(player, playerState);
        await contract.addPlayerToTeam(player, 0);
        name = await contract.teamNameByPlayer(player).should.be.fulfilled;
        name.should.be.equal(team);
    });
});
