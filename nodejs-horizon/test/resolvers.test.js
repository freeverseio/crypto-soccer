const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const ganache = require('ganache-cli');

const identity = {
    address: '0x3Abf1775944E2B2C15c05D044632831f0Dfc9130',
    privateKey: '0x0a69684608770d018143dd70dc5dc5b6beadc366b87e45fcb567fc09407e7fe5'
};

// we preset the balance of our identities to 100 ether
const provider = ganache.provider({
    accounts: [{ secretKey: identity.privateKey, balance: Web3Utils.toWei('100', 'ether') }]
});

const web3 = new Web3(provider, null, {});

