const BigNumber = web3.BigNumber;
require('chai')
    .use(require('chai-bignumber')(BigNumber))
    .use(require('chai-as-promised'))
    .should();

const Gateway = artifacts.require('Gateway');

contract('Gateway', (accounts) => {

    console.log(accounts[0])
    it('correct deployed', async () => {
        const validator = accounts[9]
        console.log("TODO: understand what's second, third constuctor param");
        const gateway = await Gateway.new([validator], 3, 4);
        gateway.should.not.equal(null);
    });

    it('deposit ETH', async () => {
        const gateway = await Gateway.new([accounts[9]], 3, 4);
        const amount = 10;
        await gateway.send(amount, {from: accounts[0]}).should.be.fulfilled;
        let result = await gateway.getETH(accounts[0]).should.be.fulfilled;
        result.should.be.bignumber.equal(amount);
    });

    // it('withdraw ETH', async () => {
    //     const gateway = await Gateway.new([accounts[9]], 3, 4);
    //     const amount = 10;
    //     await gateway.send(amount, { from: accounts[0] }).should.be.fulfilled;
    //     const signature= web3.eth.sign(accounts[0], '0x' + amount);
    //     await gateway.withdrawETH(amount, signature); //.should.be.fulfilled;
    // });
});
