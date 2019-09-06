const AssetsLib = artifacts.require('AssetsLib');
const Assets = artifacts.require('Market');

module.exports = function (deployer) {
  deployer.then(async () => {
      deployReceipt = await deployer.deploy(Assets, AssetsLib.address);
    })
    .catch(console.error);
};

