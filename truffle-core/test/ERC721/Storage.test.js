require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('Storage');

contract('CryptoTeams', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Storage.new().should.be.fulfilled;
    });

    it('no initial players', async () =>{
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('no initial teams', async () =>{
        const count = await contract.getNCreatedTeams().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    })
});
