require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('PlayersMintableMock');

contract('PlayersMintable', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Players.new().should.be.fulfilled;
    });

    it('no initial players', async () => {
        const count = await contract.totalSupply().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('mint 2 player with same name', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        await contract.mint(accounts[0], "player").should.be.rejected;
    });

    it('name is correct', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("player");
    });

    it('compute id from name', async () => {
        const id = await contract.computeId("player").should.be.fulfilled;
        id.toNumber().should.be.equal(2.28092867984879e+76);
    });

    it('get player id of existing player', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        await contract.getPlayerId("player").should.be.fulfilled;
    });

    it('get player id of unexisting player', async () => {
        await contract.getPlayerId("player").should.be.rejected;
    });

    it('minted player skills sum is 250', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        const defence = await contract.getDefence(id).should.be.fulfilled;
        const speed = await contract.getSpeed(id).should.be.fulfilled;
        const pass = await contract.getPass(id).should.be.fulfilled;
        const shoot = await contract.getShoot(id).should.be.fulfilled;
        const endurance = await contract.getEndurance(id).should.be.fulfilled;
        const sum = defence.toNumber() + speed.toNumber() + pass.toNumber() + shoot.toNumber() + endurance.toNumber();
        sum.should.be.equal(250);
    });

    it('sum of computed skills is 250', async () => {
        for (let i = 0; i < 10; i++) {
            const skills = await contract.computeSkills(Math.random()).should.be.fulfilled;
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            sum.should.be.equal(250);
        }
    });
});
