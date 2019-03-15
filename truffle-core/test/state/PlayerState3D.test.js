const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('PlayerState3D');

contract('PlayerState3D', (accounts) => {
    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
    });
}) 