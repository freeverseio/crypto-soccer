require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('Storage');

contract('Storage', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await Storage.new().should.be.fulfilled;
    });

    it('number of players per team', async () => {
        const result = await instance.getPlayersPerTeam().should.be.fulfilled;
        result.toNumber().should.be.equal(11);
    })
});