require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let engine = null;

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
    });
});