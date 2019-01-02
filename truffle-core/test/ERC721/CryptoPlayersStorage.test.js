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

    it('set genome', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const defence = 0x01;
        const speed = 0x02;
        const pass = 0x03;
        const shoot = 0x04;
        const endurance = 0x05;
        await contract.setGenome(id, defence, speed, pass, shoot, endurance).should.be.fulfilled;
        const genome = await contract.getGenome(id).should.be.fulfilled;
        genome.toString(16).should.be.equal('50004000300020001');
    });

    it('get skills', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const defence = 0x01;
        const speed = 0x02;
        const pass = 0x03;
        const shoot = 0x04;
        const endurance = 0x05;
        await contract.setGenome(id, defence, speed, pass, shoot, endurance).should.be.fulfilled;
        let skill = await contract.getDefence(id).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence);
        skill = await contract.getSpeed(id).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed);
        skill = await contract.getPass(id).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass);
        skill = await contract.getShoot(id).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot);
        skill = await contract.getEndurance(id).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance);
    });
});
