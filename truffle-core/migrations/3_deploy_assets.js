const Assets = artifacts.require('Market');

module.exports = function (deployer) {
  deployer.then(async () => {
      deployReceipt = await deployer.deploy(Assets);
    })
    .catch(console.error);
};

