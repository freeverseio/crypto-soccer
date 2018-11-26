require('chai')
    .use(require('chai-as-promised'))
    .should();


const TeamFactory = artifacts.require('TeamFactoryMock');
const GameEngine = artifacts.require('GameEngine');
const f = require('../../jsCommons/functions.js');
const k = require('../../jsCommons/constants.js');

contract('GameEngine', (accounts) => {
    let teamFactory;
    let instance;

    beforeEach(async () => {
        teamFactory = await TeamFactory.new().should.be.fulfilled;
        instance = await GameEngine.new(teamFactory.address).should.be.fulfilled;
    });

    it('play game with unexistent teams', async () => {
        await instance.playGame(0,0, 444).should.be.rejected;
    });
    
    it('play game', async () => {
        const team0 = await f.createTeam(
            teamFactory, 
            "teamName0", 
            "playerBaseName0", 
            k.MaxPlayersInTeam, 
            f.createAlineacion(4,3,3)
            ).should.be.fulfilled;
        team0.toNumber().should.be.equal(0);
        const team1 = await f.createTeam(
            teamFactory, 
            "teamName1", 
            "playerBaseName1", 
            k.MaxPlayersInTeam, 
            f.createAlineacion(4,3,3)
            ).should.be.fulfilled;
        team1.toNumber().should.be.equal(1);
        await instance.playGame(team0, team1, 444).should.be.fulfilled;
    });

});