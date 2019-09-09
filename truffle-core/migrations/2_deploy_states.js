const Encoding = artifacts.require('Encoding');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Encoding);
    })
    .catch(console.error);
};

