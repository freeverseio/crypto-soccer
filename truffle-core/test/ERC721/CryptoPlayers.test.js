require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersMock');

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
        const name = "player";
        const state = 34324;
        const teamId = 1;
        await contract.addPlayer(name, state, teamId, accounts[0]).should.be.fulfilled;
        const count = await contract.getNCreatedPlayers().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
        const nameResult = await contract.getPlayerName(count);
        nameResult.should.be.equal(name);
        const stateResult = await contract.getPlayerState(count);
        stateResult.toNumber().should.be.equal(state);
    });
});
