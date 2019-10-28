const SortIdxs = artifacts.require('SortIdxs');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(SortIdxs);
    })
    .catch(console.error);
};

