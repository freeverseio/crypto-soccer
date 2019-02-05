require('chai')
    .use(require('chai-as-promised'))
    .should();

const Teams = artifacts.require('Teams');

contract('Teams', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await Teams.new().should.be.fulfilled;
    });
});
