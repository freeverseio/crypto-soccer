const EncodingIDs = artifacts.require('EncodingIDs');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(EncodingIDs);
    })
    .catch(console.error);
};

