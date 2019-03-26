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
        await league.updateUsersAlongDataHash(id = 0, '0x43').should.be.rejected;
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
        await league.updateUsersAlongDataHash(id, '0x34').should.be.rejected;
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
        await league.updateUsersAlongDataHash(id, '0x5464').should.be.fulfilled;
    });

    it('compute user along data hash', async () => {
        let hash = await league.computeUsersAlongDataHash(teamIds = [0], tactic = [[4, 4, 2]], block = [3]).should.be.fulfilled;
        hash.should.be.equal('0x23f31280f69accf85f4ed1f35b9b7c8120241435f7f1c7005d1a397e09035c4b');
        hash = await league.computeUsersAlongDataHash(teamIds = [0], tactic = [[4, 4, 2]], block = [2]).should.be.fulfilled;
        hash.should.be.equal('0x94cc21c8dfb0a81fb883059124ef97d417f42f86c1caa0c248ae05eda99ff245');
        hash = await league.computeUsersAlongDataHash(teamIds = [0, 1], tactic = [[4, 4, 2], [4, 4, 2]], block = [2, 4]).should.be.fulfilled;
        hash.should.be.equal('0xb854db7f540d4de46dd8e42fdaf48fed19057ae2e0e60e8be3c460647ceae2d6');
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