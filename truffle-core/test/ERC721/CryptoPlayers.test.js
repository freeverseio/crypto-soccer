require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');

contract('CryptoPlayers', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    })

    it('check name and symbol', async () => {
        const contract = await CryptoPlayers.new();
        await contract.name().should.eventually.equal("CryptoSoccerPlayers");
        await contract.symbol().should.eventually.equal("CSP");
    });

    it('no initial players', async () => {
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('add player', async () => {
        const state = 34324;
        await contract.addPlayer("player", state, accounts[0]).should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        const name = await contract.getName(id);
        name.should.be.equal("player");
        const stateResult = await contract.getState(id);
        stateResult.toNumber().should.be.equal(state);
    });
});
