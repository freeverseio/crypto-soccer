const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueUsersAlongData');

contract('LeagueUsersAlongData', (accounts) => {
    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
    });

    it('initial hash of unexistent league', async () => {
        await league.getUsersAlongDataHash(id = 3).should.be.rejected;
    });

    it('initial hash of existing league', async () => {
        await league.create(
            id = 0, 
            initBlock = 1, 
            step = 1, 
            teamIds = [1, 2], 
            tactics = [[4, 4, 3], [4, 4, 3]]
        ).should.be.fulfilled;
        const hash = await league.getUsersAlongDataHash(id).should.be.fulfilled;
        hash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
    });

    it('update unexistent league', async () => {
        await league.updateUsersAlongDataHash(id = 0, teamIdx = [0], tactic = [[4, 4, 3]]).should.be.rejected;
    })

    it('update finished league', async () => {
        await league.create(
            id = 0, 
            initBlock = 1, 
            step = 1, 
            teamIds = [1, 2], 
            tactics = [[4, 4, 3], [4, 4, 3]]
        ).should.be.fulfilled;
        const finished = await league.hasFinished(id).should.be.fulfilled;
        finished.should.be.equal(true);
        await league.updateUsersAlongDataHash(id, teamIdx = [0], tactic = [[4, 4, 2]]).should.be.rejected;
    });

    it('update', async () => {
        await league.create(
            id = 0, 
            initBlock = 1, 
            step = 100000, 
            teamIds = [1, 2], 
            tactics = [[4, 4, 3], [4, 4, 3]]
        ).should.be.fulfilled;
        const finished = await league.hasFinished(id).should.be.fulfilled;
        finished.should.be.equal(false);
        await league.updateUsersAlongDataHash(id, teamIdx = [0], tactic = [[4, 4, 2]]).should.be.fulfilled;
        await league.updateUsersAlongDataHash(id, teamIdx = [0], tactic = [[4, 4, 2]]).should.be.fulfilled;
    });

    // it('update with wrong teamIdx', async () => {
    //     await league.create(
    //         id = 0, 
    //         initBlock = 1, 
    //         step = 100000, 
    //         teamIds = [1, 2], 
    //         tactics = [[4, 4, 3], [4, 4, 3]]
    //     ).should.be.fulfilled;
    //     const finished = await league.hasFinished(id).should.be.fulfilled;
    //     finished.should.be.equal(false);
    //     await league.updateUsersAlongDataHash(id, teamIdx = [1, 3], tactic = [[4, 4, 2], [4, 4, 2]]).should.be.rejected;
    // });
}) 