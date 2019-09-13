const Updates = artifacts.require('Updates');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(Updates);
    })
    .catch(console.error);
};

