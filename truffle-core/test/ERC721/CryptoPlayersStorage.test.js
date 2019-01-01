require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersStorageMock');

contract('CryptoPlayersStorage', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    });

    it('default name', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("");
    });

    it('default team', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const team = await contract.getTeam(id).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });
    
    it('default state', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const state = await contract.getState(id).should.be.fulfilled;
        state.toNumber().should.be.equal(0);
    });

    it('set name', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        await contract.setName(id, "player").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("player");
    });
    
    it('set team', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        await contract.setTeam(id, 1);
        const team = await contract.getTeam(id).should.be.fulfilled;
        team.toNumber().should.be.equal(1);
    });

    it('set state', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        await contract.setState(id, 1).should.be.fulfilled;
        const state = await contract.getState(id).should.be.fulfilled;
        state.toNumber().should.be.equal(1);
    });

    it('default skills', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const defence = await contract.getDefence(id).should.be.fulfilled;
        defence.toNumber().should.be.equal(0);
        const speed = await contract.getSpeed(id).should.be.fulfilled;
        speed.toNumber().should.be.equal(0);
        const pass = await contract.getPass(id).should.be.fulfilled;
        pass.toNumber().should.be.equal(0);
        const shoot = await contract.getShoot(id).should.be.fulfilled;
        shoot.toNumber().should.be.equal(0);
        const endurance = await contract.getEndurance(id).should.be.fulfilled;
        endurance.toNumber().should.be.equal(0);
    });
});
