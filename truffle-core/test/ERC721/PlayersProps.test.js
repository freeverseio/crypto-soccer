require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('PlayersPropsMock');

contract('PlayersProps', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Players.new().should.be.fulfilled;
    });

    it('default name', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("");
    });

    it('set name', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        await contract.setName(id, "player").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("player");
    });

    it('default genome', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const genome = await contract.getPlayerState(id).should.be.fulfilled;
        genome.toString(16).should.be.equal('0');
    });

    it('set genome', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const birth = 12;
        const defence = 0x01;
        const speed = 0x02;
        const pass = 0x03;
        const shoot = 0x04;
        const endurance = 0x05;
        const currentTeamId = 0x01;
        const currentShirtNum = 0x01;
        const prevLeagueId = 0x01;
        const prevTeamPosInLeague = 0x01;
        const prevShirtNumInLeague = 0x01;
        const lastSaleBlock = 0x01;
        await contract.setGenome(
            id, birth, defence, speed, pass, shoot, endurance,
            currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock 
        ).should.be.fulfilled;
        const genome = await contract.getPlayerState(id).should.be.fulfilled;
//        genome.toString(16).should.be.equal('14004000c002000400c');
        genome.toString(16).should.be.equal('40000000088080000440000014004000c002000400c');
    });

    it('get infos of unexistent player', async () => {
        const id = 1;
        await contract.getPlayerState(id).should.be.rejected;
        await contract.getBirth(id).should.be.rejected;
        await contract.getDefence(id).should.be.rejected;
        await contract.getSpeed(id).should.be.rejected;
        await contract.getPass(id).should.be.rejected;
        await contract.getShoot(id).should.be.rejected;
        await contract.getEndurance(id).should.be.rejected;
    });

    it('get infos coded into genome', async () => {
        const id = 1;
        await contract.mint(accounts[0], id).should.be.fulfilled;
        const birth = 12;
        const defence = 0x01;
        const speed = 0x02;
        const pass = 0x03;
        const shoot = 0x04;
        const endurance = 0x05;
        const currentTeamId = 0x01;
        const currentShirtNum = 0x01;
        const prevLeagueId = 0x01;
        const prevTeamPosInLeague = 0x01;
        const prevShirtNumInLeague = 0x01;
        const lastSaleBlock = 0x01;
        await contract.setGenome(
            id, birth, defence, speed, pass, shoot, endurance,
            currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock 
        ).should.be.fulfilled;
        let result = await contract.getBirth(id).should.be.fulfilled;
        result.toNumber().should.be.equal(birth);
        result = await contract.getDefence(id).should.be.fulfilled;
        result.toNumber().should.be.equal(defence);
        result = await contract.getSpeed(id).should.be.fulfilled;
        result.toNumber().should.be.equal(speed);
        result = await contract.getPass(id).should.be.fulfilled;
        result.toNumber().should.be.equal(pass);
        result = await contract.getShoot(id).should.be.fulfilled;
        result.toNumber().should.be.equal(shoot);
        result = await contract.getEndurance(id).should.be.fulfilled;
        result.toNumber().should.be.equal(endurance);
    });
});
