require('chai')
    .use(require('chai-as-promised'))
    .should();

const Storage = artifacts.require('Storage');

contract('Storage', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await Storage.new().should.be.fulfilled;
    });
});