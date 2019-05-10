const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Players = artifacts.require('PlayersMock');

contract('Players', (accounts) => {
    let players = null;

    beforeEach(async () => {
        players = await Players.new().should.be.fulfilled;
    });

    it('query null player id', async () => {
        await players.getPlayerTeam(0).should.be.rejected;
    });

    it('query non created player id', async () => {
        await players.getPlayerTeam(1).should.be.rejected;
    });

    it('get player team of existing player', async () => {
        const nPLayersPerTeam = await players.getPlayersPerTeam().should.be.fulfilled;
        await players.addTeam("Barca").should.be.fulfilled;
        for (let playerId=1 ; playerId <= nPLayersPerTeam ; playerId++){
            const teamId = await players.getPlayerTeam(playerId).should.be.fulfilled;
            teamId.toNumber().should.be.equal(1);
        }
        await players.getPlayerTeam(nPLayersPerTeam+1).should.be.rejected;
    });

    it('computed skills with rnd = 0 is 50 each', async () => {
        let skills = await players.computeSkills(0).should.be.fulfilled;
        skills.forEach(skill => (skill.toNumber().should.be.equal(50)));
    });

    it('int hash is deterministic', async () => {
        const rand0 = await players.intHash("Barca0").should.be.fulfilled;
        const rand1 = await players.intHash("Barca0").should.be.fulfilled;
        rand0.should.be.bignumber.equal(rand1);
        const rand2 = await players.intHash("Barca1").should.be.fulfilled;
        rand0.should.be.bignumber.not.equal(rand2);
        rand0.should.be.bignumber.equal('64856073772839990506814373782217928521534618466099710722049665631602958392435');
    });

    it('sum of computed skills is 250', async () => {
        for (let i = 0; i < 10; i++) {
            const seed = await players.intHash("Barca" + i).should.be.fulfilled;
            const skills = await players.computeSkills(seed).should.be.fulfilled;
            const sum = skills.reduce((a, b) => a + b.toNumber(), 0);
            sum.should.be.equal(250);
        }
    });


    // it('minted player skills sum is 250', async () => {
    //     await contract.mint(accounts[0], "player").should.be.fulfilled;
    //     await 
    //     const id = await contract.getPlayerId("player").should.be.fulfilled;
    //     const defence = await contract.getDefence(id).should.be.fulfilled;
    //     const speed = await contract.getSpeed(id).should.be.fulfilled;
    //     const pass = await contract.getPass(id).should.be.fulfilled;
    //     const shoot = await contract.getShoot(id).should.be.fulfilled;
    //     const endurance = await contract.getEndurance(id).should.be.fulfilled;
    //     const sum = defence.toNumber() + speed.toNumber() + pass.toNumber() + shoot.toNumber() + endurance.toNumber();
    //     sum.should.be.equal(250);
    // });

    

    // it('get skills of player', async () => {
    //     await players
    // })
});
 