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
        await league.updateUsersAlongDataHash(id = 0, teamIdx = 0, tactic = [4, 4, 3]).should.be.rejected;
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
        await league.updateUsersAlongDataHash(id, teamIdx = 0, tactic = [4, 4, 2]).should.be.rejected;
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
        await league.updateUsersAlongDataHash(id, teamIdx = 0, tactic = [4, 4, 2]).should.be.fulfilled;
        let hash = await league.getUsersAlongDataHash(id).should.be.fulfilled;
        hash.should.be.equal('0xd0c6de7df5b6f49e55dab2d876310253a8da499e94d6de0e3a50fc437fd9d75d');
        await league.updateUsersAlongDataHash(id, teamIdx = 0, tactic = [4, 4, 2]).should.be.fulfilled;
        hash = await league.getUsersAlongDataHash(id).should.be.fulfilled;
        hash.should.be.equal('0x3d829cedc91a56c3fc097e80a4a70ee30d63ede10f0356d0565741208fcb8cb6');
    });

    // TODO: reactive
    // it('update with wrong teamIdx', async () => {
    //     const id = 0;
    //     await league.create(id, initBlock = 1, step = 100000, teamIds = [1, 2]).should.be.fulfilled;
    //     const finished = await league.hasFinished(id).should.be.fulfilled;
    //     finished.should.be.equal(false);
    //     await league.updateUsersAlongDataHash(id, teamIdx = 2, tactic = [4, 4, 2]).should.be.rejected;
    // });
}) 