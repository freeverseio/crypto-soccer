const Sort = artifacts.require('Sort');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Sort);
    })
    .catch(console.error);
};

