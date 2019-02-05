require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');

contract('Players', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Players.new().should.be.fulfilled;
    })
});
