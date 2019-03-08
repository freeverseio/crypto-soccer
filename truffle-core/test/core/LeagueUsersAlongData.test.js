const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueUsersAlongData');

contract('LeagueUsersAlongData', (accounts) => {
    const id = 0;

    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
        const initBlock = 1;
        const step = 1;
        const teamIds = [1, 2];
        await league.create(id, initBlock, step, teamIds).should.be.fulfilled;
    });

    it('initial hash of unexistent league', async () => {
        await league.getUsersAlongDataHash(3).should.be.rejected;
    });

    it('initial hash of existing league', async () => {
        const hash = await league.getUsersAlongDataHash(id).should.be.fulfilled;
        hash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
    });
}) 