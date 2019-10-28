const SortValues = artifacts.require('SortValues');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(SortValues);
    })
    .catch(console.error);
};

