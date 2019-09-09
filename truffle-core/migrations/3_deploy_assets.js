const Encoding = artifacts.require('Encoding');
const Assets = artifacts.require('Market');

module.exports = function (deployer) {
  deployer.then(async () => {
      deployReceipt = await deployer.deploy(Assets, Encoding.address);
    })
    .catch(console.error);
};

