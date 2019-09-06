const AssetsLib = artifacts.require('AssetsLib');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(AssetsLib);
    })
    .catch(console.error);
};

